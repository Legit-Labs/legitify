<div align="center">
 
[![Build & Test](https://github.com/Legit-Labs/legitify/actions/workflows/build_and_test.yaml/badge.svg)](https://github.com/Legit-Labs/legitify/actions/workflows/build_and_test.yaml)
[![Code Analysis](https://github.com/Legit-Labs/legitify/actions/workflows/code_analysis.yaml/badge.svg)](https://github.com/Legit-Labs/legitify/actions/workflows/code_analysis.yaml)
[![Version Releaser](https://github.com/Legit-Labs/legitify/actions/workflows/release.yaml/badge.svg)](https://github.com/Legit-Labs/legitify/actions/workflows/release.yaml)
[![Build Docs](https://github.com/Legit-Labs/legitify/actions/workflows/build_docs.yaml/badge.svg)](https://github.com/Legit-Labs/legitify/actions/workflows/build_docs.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Legit-Labs/legitify)](https://goreportcard.com/report/github.com/Legit-Labs/legitify)
<img referrerpolicy="no-referrer-when-downgrade" src="https://static.scarf.sh/a.png?x-pxid=6f4cbb25-54f4-4c47-b611-9b741732bb86" />
<br/>
 <img width="250" alt="Legitify Logo" src="https://user-images.githubusercontent.com/74864790/174815311-746a0c98-9a1f-44a9-808c-035788edfd4d.png">

Strengthen the security posture of your source-code management! <br/>
Detect and remediate misconfigurations, security and compliance issues across all your GitHub and GitLab assets with ease 🔥 <br/>
by [Legit Security](https://www.legitsecurity.com/).

<b>
Wonder what Legit Security does?
</b>

Legit Security is an application security posture management (ASPM) and software supply chain security solution.<br/>
For more information check out the [comparison table](#legitify-vs-the-legit-security-platform)

</div>


https://user-images.githubusercontent.com/107790206/210602039-2d022692-87ea-4005-b9c6-f091158de3ce.mov

## Installation

Installation is possible in several ways:

- For macOS (or linux) using homebrew:

```
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

## CI - Legitify Custom GitHub Action

You can run legitify as part of a CI process with the legitify Custom GitHub Actions:

```
name: Legitify Analyze
on:
    workflow_dispatch:
    schedule:
      - cron: '0 11 * * 1-5'

jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - name: Legitify Action
        uses: Legit-Labs/legitify@main
        with:
          github_token: ${{ secrets.PAT_FOR_LEGITIFY }}
          ignore-policies: |
             non_admins_can_create_public_repositories
             requires_status_checks
```

Checkout the [action file](https://github.com/Legit-Labs/legitify/blob/main/action.yml) for additional parameters
and configuration.

## Provenance

To enhance the software supply chain security of legitify's users, as of v0.1.6, every legitify release contains a [SLSA Level 3 Provenance](https://github.com/slsa-framework/slsa-github-generator/blob/main/internal/builders/generic/README.md) document.  
The provenance document refers to all artifacts in the release, as well as the generated docker image.  
You can use [SLSA framework's official verifier](https://github.com/slsa-framework/slsa-verifier) to verify the provenance.  
Example of usage for the darwin_arm64 architecture for the v0.1.6 release:

```
VERSION=0.1.6
ARCH=darwin_arm64
./slsa-verifier verify-artifact --source-branch main --builder-id 'https://github.com/slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@refs/tags/v1.2.2' --source-uri "git+https://github.com/Legit-Labs/legitify" --provenance-path multiple.intoto.jsonl ./legitify_${VERSION}_${ARCH}.tar.gz
```

## Commands

### analyze

```
SCM_TOKEN=<your_token> legitify analyze
```

By default, legitify will check the policies against all your resources (organizations, repositories, members, actions). Archived repositories are skipped.

You can control which resources will be analyzed with command-line flags namespace and org:

- `--namespace (-n)`: will analyze policies that relate to the specified resources
- `--org`: will limit the analysis to the specified GitHub organizations or GitLab group, excluding archived repositories
- `--repo`: will limit the analysis to the specified GitHub repositories or GitLab projects
- `--scm`: specify the source code management platform. Possible values are: `github` or `gitlab`. Defaults to `github`. Please note: when running on GitLab, `--scm gitlab` is required.
- `--enterprise`: will specify which enterprises should be analyzed. Please note: in order to analyze an enterprise, an enterprise slug must be provided.

```
SCM_TOKEN=<your_token> legitify analyze --org org1,org2 --namespace organization,member
```

The above command will test organization and member policies against org1 and org2.

### gpt-analysis

```
SCM_TOKEN=<your_token> OPENAI_TOKEN=<token> ./legitify gpt-analysis --repo org1/repo1 --org org1
```

GPT-3 based analysis of the security posture of the provided repository or organization.

**NOTE: The repository/organization metadata is sent to openai servers.**

Flags:

- `--org`: will limit the analysis to the specified GitHub organizations or GitLab group
- `--repo`: will limit the analysis to the specified GitHub repositories or GitLab projects
- `--scm`: specify the source code management platform. Possible values are: `github` or `gitlab`. Defaults to `github`.
- `--token`: token for the SCM (or set the SCM_TOKEN environment variable)
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
2. legitify requires a GitHub personal access token (PAT) to analyze your resources successfully, which can be either provided as an argument (`-t`) or as an environment variable (`SCM_TOKEN`).
   The PAT needs the following scopes for full analysis:

```
admin:org, read:enterprise, admin:org_hook, read:org, repo, read:repo_hook
```

See [Creating a Personal Access Token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) for more information.  
Fine-grained personal access tokens are currently not supported.

### GitHub Enterprise Server

You can run legitify against a GitHub Enterprise Server instance if you set the endpoint URL in the environment variable `SERVER_URL`:

```sh
export SERVER_URL="https://github.example.com/"
SCM_TOKEN=<your_token> legitify analyze --org org1,org2 --namespace organization,member
```

### GitLab Cloud/Server

1. As mentioned in the previous section, you need to be an owner of at least one GitLab group. Otherwise, you can still use the tool if you're an admin of at least one project inside a group, in which case you'll be able to see only project-related policies results.
2. legitify requires a GitLab personal access token (PAT) to analyze your resources successfully, which can be either provided as an argument (`-t`) or as an environment variable (`SCM_TOKEN`).
   The PAT needs the following scopes for full analysis:
   `   read_api, read_user, read_repository, read_registry`
   See [Creating a Personal Access Token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html) for more information.  
   To run legitify against GitLab Cloud set the scm flag to gitlab `--scm gitlab`, to run against GitLab Server you need to provide also a SERVER_URL:

```sh
export SERVER_URL="https://gitlab.example.com/"
SCM_TOKEN=<your_token> legitify analyze --namespace organization --scm gitlab
```

> **_NOTE 1:_** To ignore invalid server certificate, please pass the `ignore-invalid-certificate` flag

> **_NOTE 2:_** For non-premium GitLab accounts some policies (such as branch protection policies) will be skipped

## Namespaces

Namespaces in legitify are resources that are collected and run against the policies.
Currently, the following namespaces are supported:

1. `organization` - GitHub organization (or GitLab group) level policies (e.g., "Two-Factor Authentication Is Not Enforced for the Organization")
2. `actions` - organization GitHub Actions policies (e.g., "GitHub Actions Runs Are Not Limited To Verified Actions")
3. `member` - contributor level policies (e.g., "Stale Admin Found")
4. `repository` - GitHub repository (or GitLab Project) level policies (e.g., "Code Review By At Least Two Reviewers Is Not Enforced"). Note: Archived repositories are ignored unless specified directly via the `--repo` argument.
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
- Use the `--ignore-policies-path $PATH` and provide a file with the policies you want to ignore to skip specific policies.
  One policy per line, e.g.
  `no_conversation_resolution
requires_status_checks                                                     ─╯`

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

## Contribution

Thank you for considering contributing to Legitify! We encourage and appreciate any kind of contribution.
Here are some resources to help you get started:

- [Contribution Guide](https://github.com/Legit-Labs/legitify/blob/main/CONTRIBUTING.md)
- [Code of Conduct](https://github.com/Legit-Labs/legitify/blob/main/CODE_OF_CONDUCT.md)
- [Open an Issue](https://github.com/Legit-Labs/legitify/issues/new/choose)
- [Open a Pull Request](https://github.com/Legit-Labs/legitify/compare)

## Support

If you have questions about legitify or need any assistance with its operation, don't hesitate to [reach out](mailto:legitify@legitsecurity.com). Our team is committed to providing support and ensuring a smooth experience.

## Legitify vs. the Legit Security platform

If you liked Legitify, you are going to love the Legit Security Platform!

- It automates Legitify checks for the entire environments, discovers more systems and shows all results in a simple web app to manage at scale.
- Legit security is a complete CI/CD security solution together with Application Security Posture Management (ASPM) that covers application security end-to-end.
- It is a SaaS platform, built for engineering, DevOps and security teams and trusted by many leading organizaitons around the world.

Below is a comfeature parison between Legitify and Legit:

| **Capability** | **Legitify** | **Legit Security Platform** |
|---|---|---|
| Supported platforms | GitHub <br> GitLab | ALL major SCMs (incl. Azure DevOps, Bitbucket and more)<br> CI/CD systems (e.g. Jenkins) <br> Package registries (e.g. JFrog Artifactory)<br>Cloud providers (e.g. AWS) |
| Risk detection | SCM Misconfigurations only | SCMs Misconfigurations <br> CI Misconfigurations <br> CD Misconfigurations <br> Package Registries Misconfigurations <br> Pipeline risks <br> Secrets <br> IaC <br> Security Incidents <br> And more... |
| Compliance report | [OSSF SCM Best Practices](https://best.openssf.org/SCM-BestPractices/) | SSDF <br> SLSA <br> SOC2 <br> ISO 27001 <br> FedRAMP <br> And more... |
| Policy drifts detection | Can be detected periodically though Legitify's GitHub Action | Get real-time alerts when a misconfiguration is introduced |
| SDLC assets management | - | Yes |
| Issue & policy management | - | Yes |
| Code To Cloud context | - | Yes (contextualized information enables smarter prioritization) |
| Workspaces & product groups | - | Yes |
| Ticketing & alerting | - | Jira, Slack, and more |
| Ingest risk | - | Import APIs and integrations with SAST, SCA and other testing solutions |
| Rest APIs | - | Yes |

To check out Legit, visit our [website](https://www.legitsecurity.com/) or directly [book a demo](https://info.legitsecurity.com/book-a-demo)

<div align="center">
 <a href="https://www.legitsecurity.com">
  <img width="250" alt="Legitify Logo" src="https://github.com/Legit-Labs/legitify/assets/74864790/c76dc765-e8fd-498e-ab92-1228eb5a1f2d">
 </a>
</div>
