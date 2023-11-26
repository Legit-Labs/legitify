const core = require("@actions/core");
const fs = require("fs");
const zlib = require("zlib");
const request = require("request");
const tar = require("tar-fs");
const path = require("path");
const fetch = require("node-fetch");
const exec = require("@actions/exec");
const { context, getOctokit } = require("@actions/github");
const artifact = require("@actions/artifact");
const { exit } = require("process");

const logoMarkdown = `<div align="left">
<a href="https://www.legitsecurity.com">
 <img width="120" alt="Legitify Logo" src="https://github.com/Legit-Labs/legitify/assets/74864790/c76dc765-e8fd-498e-ab92-1228eb5a1f2d">
 </a>
</div>`

async function uploadErrorLog() {
  const fileName = "error.log";
  try {
    if (fs.existsSync(fileName)) {
      const client = artifact.create();
      await client.uploadArtifact(fileName, [fileName], ".", context.runId);
      console.log(`Uploaded ${fileName} to the workflow artifact`);
    } else {
      console.log(`File ${fileName} does not exist so skipping upload`);
    }
  } catch (error) {
    console.error(error);
  }
}

async function isPrivateRepo(token) {
  const repoName = process.env.GITHUB_REPOSITORY;
  if (!repoName) {
    core.setFailed("No repository detected.");
    return;
  }
  const parts = repoName.split("/");
  const owner = parts[0];
  const repo = parts[1];

  const octokit = getOctokit(token);
  const repoResp = await octokit.rest.repos.get({
    owner,
    repo,
  });

  return repoResp.data.private;
}

async function executeLegitify(token, args, uploadCodeScanning) {
  let myOutput = "";
  let myError = "";

  const options = {};
  options.listeners = {
    stderr: (data) => {
      myError += data.toString();
    },
  };
  options.env = { SCM_TOKEN: token };
  options.silent = true
  const isPrivate = await isPrivateRepo(token)
  core.setOutput("is_private", isPrivate)
  if (!isPrivate) {
    console.log("running in a public repository - only uploading results to Code Scanning (Security Center)")
  }

  try {
    const sarifFile = "legitify-output.sarif"

    // generate the output as json
    const jsonFile = "legitify-output.json"
    const analyzeArgs = ["analyze", ...args, "--output-format", "json", "--output-file", jsonFile]
    console.log("execute legitify analyze:", analyzeArgs)
    await exec.exec("./legitify", analyzeArgs, options)

    // generate a sarif version for the code scanning
    if (uploadCodeScanning) {
      myError = ""
      const convertSarifArgs = ["convert", "--input-file", jsonFile, "--output-format", "sarif", "--output-file", sarifFile]
      console.log("execute legitify convert sarif:", convertSarifArgs)
      await exec.exec("./legitify", convertSarifArgs, options);
    }

    // generate a markdown version for the action output
    myError = ""
    options.listeners.stdout = (data) => {
      myOutput += data.toString();
    }
    const convertMarkdownArgs = ["convert", "--input-file", jsonFile, "--output-format", "markdown"]
    console.log("execute legitify convert markdown:", convertMarkdownArgs)
    await exec.exec('"./legitify"', convertMarkdownArgs, options)
    if (isPrivate) {
      console.log(logoMarkdown)
      fs.writeFileSync(process.env.GITHUB_STEP_SUMMARY, logoMarkdown + '\n\n' + myOutput)
    } else {
      fs.unlinkSync(jsonFile);
      fs.unlinkSync("error.log");
    }
  } catch (error) {
    console.log(error.toString() + " | stderr: " + myError.toString())
    fs.writeFileSync(process.env.GITHUB_STEP_SUMMARY, "legitify failed with:\n" + myError)
    core.setFailed(error);
    exit(1);
  }
}

async function fetchLegitifyReleaseUrl(baseVersion) {
  try {
    const response = await fetch(
      "https://api.github.com/repos/Legit-Labs/legitify/releases"
    );
    if (!response.ok) {
      core.setFailed(`Failed to fetch releases: ${response.statusText}`);
      exit(1);
    }
    const releases = await response.json();

    for (const release of releases) {
      const version = release.tag_name.slice(1);
      if (version.startsWith(baseVersion)) {
        const linuxAsset = release.assets.find(
          (asset) =>
            asset.name.endsWith(".tar.gz") && asset.name.includes("linux_amd64")
        );
        return linuxAsset.browser_download_url;
      }
    }

    throw new Error(
      `No releases found with version starting with ${baseVersion}`
    );
  } catch (error) {
    core.setFailed(error);
    exit(1);
  }
}

function breakStringToParams(str) {
  const pattern = /(-{1,2}\S+)(?:\s+((?!\s*-).+?)(?=\s+-{1,2}|\s*$))?/g;
  return str.match(pattern)
}

function generateAnalyzeArgs(repo, owner) {
  let args = [];

  const scorecard = process.env["scorecard"];
  if (scorecard === "yes" || scorecard === "verbose") {
    args.push("--scorecard");
    args.push(scorecard);
  }

  if (process.env["analyze_self_only"] === "true") {
    args.push("--repo");
    args.push(repo);
  } else if (process.env["repositories"] !== "") {
    args.push("--repo");
    args.push(process.env["repositories"]);
  } else {
    args.push("--org");
    args.push(owner);
  }

  if (process.env["ignore-policies-file"] !== "") {
    args.push("--ignore-policies-file");
    args.push(process.env["ignore-policies-file"]);
  }

  const extra = process.env["extra"]
  if (extra !== "") {
    const splitArgs = breakStringToParams(extra)
    args.push(...splitArgs)
  }

  return args;
}

function downloadAndExtract(fileUrl, filePath) {
  console.log(
    `downloading legitify binary from the following release URL: ${fileUrl}`
  );
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(filePath);

    request(fileUrl)
      .on("error", (error) => {
        reject(error);
      })
      .pipe(file)
      .on("close", () => {
        const readStream = fs.createReadStream(filePath);
        const extractor = zlib.createGunzip();
        readStream
          .on("error", (error) => {
            reject(error);
          })
          .pipe(extractor)
          .pipe(tar.extract())
          .on("finish", () => {
            resolve();
          });
      });
  });
}

async function run() {
  try {
    const token = process.env["github_token"];
    if (!token) {
      core.setFailed("No GitHub token provided");
      exit(1);
    }

    const owner = process.env["GITHUB_REPOSITORY_OWNER"];
    const repo = process.env["GITHUB_REPOSITORY"];
    const uploadCodeScanning = (process.env["upload_code_scanning"] === "true");

    if (process.env["compile_legitify"] === "true") {
      console.log("using the compiled legitify version.");
    } else {
      const legitifyBaseVersion = process.env["legitify_base_version"];
      const fileUrl = await fetchLegitifyReleaseUrl(legitifyBaseVersion);
      const filePath = path.join(__dirname, "legitify.tar.gz");
      await downloadAndExtract(fileUrl, filePath);
    }

    const args = generateAnalyzeArgs(repo, owner);
    await executeLegitify(token, args, uploadCodeScanning);
  } catch (error) {
    core.setFailed(error.message);
    exit(1);
  }

  uploadErrorLog();
}

run();
