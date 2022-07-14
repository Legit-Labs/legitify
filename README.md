<div align="center">
 
[![golangci-lint](https://github.com/Legit-Labs/legitify/actions/workflows/build_and_lint.yaml/badge.svg)](https://github.com/Legit-Labs/legitify/actions/workflows/build_and_lint.yaml)
[![goreleaser](https://github.com/Legit-Labs/legitify/actions/workflows/release.yml/badge.svg)](https://github.com/Legit-Labs/legitify/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Legit-Labs/legitify)](https://goreportcard.com/report/github.com/Legit-Labs/legitify)

 <img width="250" alt="Legitify Logo" src="https://user-images.githubusercontent.com/74864790/174815311-746a0c98-9a1f-44a9-808c-035788edfd4d.png">
 
Strengthen the security posture of your GitHub organization! <br/>
Detect and remediate misconfigurations, security and compliance issues across all your GitHub assets with ease 🔥 <br/>
 by Legit Security
</div>

https://user-images.githubusercontent.com/74864790/178964716-825840a6-d714-4b1d-a41e-efa2728507a6.mp4

## Installation
1. You can download the latest legitify release from https://github.com/Legit-Labs/legitify/releases, each archive contains:
  * Legitify binary for the desired platform
  * Built-in policies provided by Legit Security
2. From source with the following steps:
```
git clone git@github.com:Legit-Labs/legitify.git
go run main.go analyze ...
```
## Requirements
legitify requires your GitHub personal access token (PAT) to analyze your organizations successfully, which can be either provided as an argument (-t) or as an environment variable ($GITHUB_ENV). Please Note:
Your user needs to be an admin or a member of at least one organization.
Additionally, to get repository-related policies results, your user should have write permissions on the repositories in the organization. The following table summarizes the user role requirements:

- legitify requires the following user role to operate successfully per resource:
  | Resource                                  | Organization |Repository |
  |-------------------------------------------|--------------|-----------|
  | Minimal role required (partial results)   | Member       | Admin     |
  | Full results and remediation capabilities | Owner        | Admin     |

- The PAT requires the following scopes for full analysis:
  ```
  admin:org, read:enterprise, admin:org_hook, read:org, repo, read:repo_hook
  ```  
See [Creating a Personal Access Token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) for more info.

## Usage
```
GITHUB_TOKEN=<your_token> legitify analyze
```
By default, legitify will check the policies against all your resources (organizations, repositories, members, actions).

You can control which resources will be analyzed with command-line flags namespace and org:
- `--namespace (-n)`: will analyze policies that relate to the specified resources 
- `--org`: will limit the analysis to the specified organizations

```
GITHUB_TOKEN=<your_token> legitify analyze --org org1,org2 --namespaces organization,member
```
The above command will test organization and member policies against org1 and org2.

## Namespaces
Namespaces in legitify are resources that are collected and run against the policies.
Currently, the following namespaces are supported:
1. `organization` - organization level policies (e.g., "Two-Factor Authentication Is Not Enforced for the Organization")
2. `actions`      - organization GitHub Actions policies (e.g., "GitHub Actions Runs Are Not Limited To Verified Actions")
3. `member`       - organization members policies (e.g., "Stale Admin Found")
4. `repository`   - repository level policies (e.g., "Code Review By At Least Two Reviewers Is Not Enforced")

By default, legitify will analyze all namespaces. You can limit only to selected ones with the `--namespace` flag, and then a comma seperated list of the selected namespaces.

## Output Options
By default, legitify will output the results in a human-readable format.
This includes the list of policy violations listed by severity,
as well as a summary table that is sorted by namespace.

### Output Formats
Using the `--output-format (-f)` flag, legitify supports outputing the results in the following formats:
1. `human-readable` - Human-readable text (default).
2. `json` - Standard JSON.

### Output Schemes
Using the `--output-scheme` flag, legitify supports outputing the results in different grouping schemes. 
Note: `--output-format=json` must be specified to output non-default schemes.
1. `flattened` - No grouping; A flat listing of the policies, each with its violations (default).
2. `group-by-namespace` - Group the policies by their namespace.
3. `group-by-resource` - Group the policies by their resource e.g. specific organization/repository.
4. `group-by-severity` - Group the policies by their severity.

### Output Destinations
- `--output-file (-o)` - full path of the output file (default: no output file, prints to stdout).
- `--error-file (-e)` - full path of the error logs (default: ./errors.log).

### Coloring
When outputing in a human-readable format, legitify support the conventional `--color[=when]` flag, which has the following options:
- `auto` - colored output if stdout is a terminal, uncolored otherwise (default).
- `always` - colored output regardless of the output destination.
- `none` - uncolored output regardless of the output destination.

### Misc
- Use the `--failed-only` flag to filter-out passed/skipped checks from the result.

## Scorecard Support
[scorecard](https://github.com/ossf/scorecard) is an OSSF's open-source project:
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

## Custom Policies
legitify comes with a set of policies in the `policies/github` directory.
In addition, you can use the `--policies-path (-p)` flag to specificy a custom directory for OPA policies.

## Contribution
Thank you for considering contributing to Legitify! We encourage and appreciate any kind of contribution.
