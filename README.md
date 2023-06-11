<div align="center">
 
[![Build & Test](https://github.com/Legit-Labs/legitify/actions/workflows/build_and_test.yaml/badge.svg)](https://github.com/Legit-Labs/legitify/actions/workflows/build_and_test.yaml)
[![Code Analysis](https://github.com/Legit-Labs/legitify/actions/workflows/code_analysis.yaml/badge.svg)](https://github.com/Legit-Labs/legitify/actions/workflows/code_analysis.yaml)
[![Version Releaser](https://github.com/Legit-Labs/legitify/actions/workflows/release.yaml/badge.svg)](https://github.com/Legit-Labs/legitify/actions/workflows/release.yaml)
[![Build Docs](https://github.com/Legit-Labs/legitify/actions/workflows/build_docs.yaml/badge.svg)](https://github.com/Legit-Labs/legitify/actions/workflows/build_docs.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Legit-Labs/legitify)](https://goreportcard.com/report/github.com/Legit-Labs/legitify)
[![Maintained By Legit Security](https://img.shields.io/badge/Maintained%20by-Legit%20Security-blueviolet)](https://www.legitsecurity.com/)
<img referrerpolicy="no-referrer-when-downgrade" src="https://static.scarf.sh/a.png?x-pxid=6f4cbb25-54f4-4c47-b611-9b741732bb86" />
 <img width="250" alt="Legitify Logo" src="https://user-images.githubusercontent.com/74864790/174815311-746a0c98-9a1f-44a9-808c-035788edfd4d.png">
 
Strengthen the security posture of your source-code management! <br/>
Detect and remediate misconfigurations, security and compliance issues across all your GitHub and GitLab assets with ease 🔥 <br/>
 by [Legit Security](https://www.legitsecurity.com/).
</div>

https://user-images.githubusercontent.com/107790206/210602039-2d022692-87ea-4005-b9c6-f091158de3ce.mov

## Installation

Installation is possible in several ways:

- For macOS (or linux) using homebrew:
```
brew tap legit-labs/legit-labs
brew install legitify
```

- You can download the latest legitify release from https://github.com/Legit-Labs/legitify/releases, each archive contains:

  - Legitify binary for the desired platform
  - Built-in policies provided by Legit Security

- From source with the following steps:

```
git clone git@github.com:Legit-Labs/legitify.git
go run main.go analyze ...
```

- As a GitHub CLI extension (check out https://github.com/Legit-Labs/gh-legitify for more information)
```
gh extension install legit-labs/gh-legitify
gh legitify
```


## Provenance

To enhance the software supply chain security of legitify's users, as of v0.1.6, every legitify release contains a [SLSA Level 3 Provenance](https://github.com/slsa-framework/slsa-github-generator/blob/main/internal/builders/generic/README.md) document.  
The provenance document refers to all artifacts in the release, as well as the generated docker image.  
You can use [SLSA framework's official verifier](https://github.com/slsa-framework/slsa-verifier) to verify the provenance.  
Example of usage for the darwin_arm64 architecture for the v0.1.6 release:

```
VERSION=0.2.7
ARCH=darwin_arm64
./slsa-verifier verify-artifact --source-branch main --builder-id 'https://github.com/slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@refs/tags/v1.6.0' --source-uri "git+https://github.com/Legit-Labs/legitify" --provenance-path multiple.intoto.jsonl ./legitify_${VERSION}_${ARCH}.tar.gz
```


## Commands

### analyze

```
LEGITIFY_TOKEN=<your_token> legitify analyze
```

By default, legitify will check the policies against all your resources (organizations, repositories, members, actions).

You can control which resources will be analyzed with command-line flags namespace and org:
- `--namespace (-n)`: will analyze policies that relate to the specified resources 
- `--org`: will limit the analysis to the specified GitHub organizations or GitLab group
- `--repo`: will limit the analysis to the specified GitHub repositories or GitLab projects
- `--scm`: specify the source code management platform. Possible values are: `github` or `gitlab`. Defaults to `github`. Please note: when running on GitLab, `--scm gitlab` is required.
- `--enterprise`: will specifiy which enterprises should be analyzed. Please note: in order to analyze an enterprise, an enterprise slug must be provided. 


```
LEGITIFY_TOKEN=<your_token> legitify analyze --org org1,org2 --namespace organization,member
```

The above command will test organization and member policies against org1 and org2.

### gpt-analysis

```
LEGITIFT_TOKEN=<your_token> OPENAI_TOKEN=<token> ./legitify gpt-analysis --repo ORG_NAME/REPO_NAME --org ORG_NAME
```

GPT-3 based analysis to the security posture of the provided repository or organization.

**NOTE: The repository/organization metadata is sent to openai servers.**

Flags:
- `--org`: will limit the analysis to the specified GitHub organizations or GitLab group
- `--repo`: will limit the analysis to the specified GitHub repositories or GitLab projects
- `--scm`: specify the source code management platform. Possible values are: `github` or `gitlab`. Defaults to `github`.
- `--token`: token for the SCM (or set the LEGITIFY_TOKEN environment variable)
- `--openai-token`: token for openai API (or set OPENAI_TOKEN environment variable) 

Must provide either `--org` or `--repo` or both.

Generating openai token:
1. Go to https://beta.openai.com/signup and create an openai account
2. Under https://platform.openai.com/account/api-keys press "Create new secret key"

## GitHub Action Usage

You can also run legitify as a GitHub action in your workflows, see the **action_examples** directory for concrete examples.


## Requirements
### GitHub (Cloud and Enterprise Server)
1. To get the most out of legitify, you need to be an owner of at least one GitHub organization. Otherwise, you can still use the tool if you're an admin of at least one repository inside an organization, in which case you'll be able to see only repository-related policies results.
2. legitify requires a GitHub personal access token (PAT) to analyze your resources successfully, which can be either provided as an argument (`-t`) or as an environment variable (`LEGITIFY_TOKEN`).
   The PAT needs the following scopes for full analysis:
  ```
  admin:org, read:enterprise, admin:org_hook, read:org, repo, read:repo_hook
  ```
See [Creating a Personal Access Token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) for more information.  
Fine-grained personal access tokens are currently not supported because they do not support [GitHub's GraphQL](https://github.blog/2022-10-18-introducing-fine-grained-personal-access-tokens-for-github/).

### GitHub Enterprise Server
You can run legitify against a GitHub Enterprise Server instance if you set the endpoint URL in the environment variable ``SERVER_URL``:


```sh
export SERVER_URL="https://github.example.com/"
LEGITIFY_TOKEN=<your_token> legitify analyze --org org1,org2 --namespace organization,member
```
### GitLab Cloud/Server
1. As mentioned in the previous section, you need to be an owner of at least one GitLab group. Otherwise, you can still use the tool if you're an admin of at least one project inside a group, in which case you'll be able to see only project-related policies results.
2. legitify requires a GitLab personal access token (PAT) to analyze your resources successfully, which can be either provided as an argument (`-t`) or as an environment variable (`LEGITIFY_TOKEN`).
  The PAT needs the following scopes for full analysis:
    ```
    read_api, read_user, read_repository, read_registry
    ```
  See [Creating a Personal Access Token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html) for more information.  
To run legitify against GitLab Cloud set the scm flag to gitlab `--scm gitlab`, to run against GitLab Server you need to provide also a SERVER_URL:

```sh
export SERVER_URL="https://gitlab.example.com/"
LEGITIFY_TOKEN=<your_token> legitify analyze --namespace organization --scm gitlab
```

## Namespaces

Namespaces in legitify are resources that are collected and run against the policies.
Currently, the following namespaces are supported:
1. `organization` - GitHub organization (or GitLab group) level policies (e.g., "Two-Factor Authentication Is Not Enforced for the Organization")
2. `actions`      - organization GitHub Actions policies (e.g., "GitHub Actions Runs Are Not Limited To Verified Actions")
3. `member`       - contributor level policies (e.g., "Stale Admin Found")
4. `repository`   - GitHub repository (or GitLab Project) level policies (e.g., "Code Review By At Least Two Reviewers Is Not Enforced")
5. `runner_group` - runner group policies (e.g, "runner can be used by public repositories")

By default, legitify will analyze all namespaces. You can limit only to selected ones with the `--namespace` flag, and then a comma separated list of the selected namespaces.

## Output Options

By default, legitify will output the results in a human-readable format.
This includes the list of policy violations listed by severity,
as well as a summary table that is sorted by namespace.

### Output Formats

Using the `--output-format (-f)` flag, legitify supports outputting the results in the following formats:

1. `human-readable` - Human-readable text (default).
2. `json` - Standard JSON.
3. `sarif` - SARIF format ([info](https://sarifweb.azurewebsites.net/)).

### Output Schemes

Using the `--output-scheme` flag, legitify supports outputting the results in different grouping schemes.
Note: `--output-format=json` must be specified to output non-default schemes.

1. `flattened` - No grouping; A flat listing of the policies, each with its violations (default).
2. `group-by-namespace` - Group the policies by their namespace.
3. `group-by-resource` - Group the policies by their resource e.g. specific organization/repository.
4. `group-by-severity` - Group the policies by their severity.

### Output Destinations

- `--output-file` - full path of the output file (default: no output file, prints to stdout).
- `--error-file` - full path of the error logs (default: ./error.log).

### Coloring

When outputting in a human-readable format, legitify support the conventional `--color[=when]` flag, which has the following options:

- `auto` - colored output if stdout is a terminal, uncolored otherwise (default).
- `always` - colored output regardless of the output destination.
- `none` - uncolored output regardless of the output destination.

### Misc

- Use the `--failed-only` flag to filter-out passed/skipped checks from the result.

## Scorecard Support - Only for GitHub server/cloud repositories
[Scorecard](https://github.com/ossf/scorecard) is an OSSF's open-source project:

> Scorecards is an automated tool that assesses a number of important heuristics ("checks") associated with software security and assigns each check a score of 0-10. You can use these scores to understand specific areas to improve in order to strengthen the security posture of your project. You can also assess the risks that dependencies introduce, and make informed decisions about accepting these risks, evaluating alternative solutions, or working with the maintainers to make improvements.

legitify supports running scorecard for all of the organization's repositories, enforcing score policies and showing the results using the `--scorecard` flag:

- `no` - do not run scorecard (default).
- `yes` - run scorecard and employ a policy that alerts on each repo score below 7.0.
- `verbose` - run scorecard, employ a policy that alerts on each repo score below 7.0, and embed its output to legitify's output.

legitify runs the following scorecard checks:
|Check|Public Repository|Private Repository|
|--|--|--|
|Security-Policy|V||
|CII-Best-Practices|V||
|Fuzzing|V||
|License|V||
|Signed-Releases|V||
|Branch-Protection|V|V|
|Code-Review|V|V|
|Contributors|V|V|
|Dangerous-Workflow|V|V|
|Dependency-Update-Tool|V|V|
|Maintained|V|V|
|Pinned-Dependencies|V|V|
|SAST|V|V|
|Token-Permissions|V|V|
|Vulnerabilities|V|V|
|Webhooks|V|V|

## Policies
legitify comes with a set of policies for each SCM in the `policies/` directory.

These policies are documented [here](https://legitify.dev).

In addition, you can use the `--policies-path (-p)` flag to specify a custom directory for OPA policies.

## Contribution

Thank you for considering contributing to Legitify! We encourage and appreciate any kind of contribution.
Here are some resources to help you get started:

- [Contribution Guide](https://github.com/Legit-Labs/legitify/blob/main/CONTRIBUTING.md)
- [Code of Conduct](https://github.com/Legit-Labs/legitify/blob/main/CODE_OF_CONDUCT.md)
- [Open an Issue](https://github.com/Legit-Labs/legitify/issues/new/choose)
- [Open a Pull Request](https://github.com/Legit-Labs/legitify/compare)
