# Welcome to Legitify contributing guide

Thank you for investing your time in contributing to our project! Any contribution you make will be reflected on the [Changelog](https://github.com/Legit-Labs/legitify/blob/main/CHANGELOG.md).

Read our  [Code of Conduct](https://github.com/Legit-Labs/legitify/blob/main/CODE_OF_CONDUCT.md)  to keep our community approachable and respectable.

In this guide you will get an overview of the contribution workflow from opening an issue, creating a PR, reviewing, and merging the PR.

## Opening Issues
If you found a bug or you have new ideas, we'd be glad to hear from you!
* Please go over existing issues before you open a new one to avoid duplicates.
* Please follow the relevant template when openning new issues.
* Please do not create a Github issue for security vulnerabilities. In such a case, please refer to [SECURITY.md](https://github.com/Legit-Labs/legitify/blob/main/SECURITY.md).

## Contributing Code
We happily welcome the contribution of code to all the components in the project!
To get an overview of the project, read the  [README](https://github.com/Legit-Labs/legitify/blob/main/README.md). Here are some resources to help you get started with contributions to Legitify:
1. Get familiar with [Open Policy Agent](https://www.openpolicyagent.org/docs/latest/) to implement policies.
2. Get familiar with [GitHub API and features](https://docs.github.com/en) to suggest new policies and implement collectors.
3. Get familiar with [GoLang idioms](https://go.dev/doc/effective_go) and conventions to contribute code.

## Getting started

Clone the repository to get started. Then you should:
* Run `make build` to build from source.

Once your changes are ready, run `make test` to test your changes.

### Commit Changes

Please keep your commits:
* Standalone - The code must compile and run successfully after each commit (no breaking commits!).
* Minimal - Break your code into minimal, logically-complete chunks.
* Self-Reviewed - Always double-check yourself before submitting.
* Tested - All tests must run. New features should come along with new tests.

All commit messages (and PRs) must start with one of the following prefixes:
* **build** - build-system changes.
* **build(deps)** - build dependencies changes (used by Dependabot).
* **ci** - Continuous Integration changes (GitHub Workflow).
* **docs** - Documentation changes.
* **feat** - New features.
* **fix** - Bug fixes.
* **perf** - Performance improvements.
* **refactor** - Code refactoring.
* **style** - Code styling changes.
* **test** - Testing changes. 

Commit messages should be:
* **Limit the first line to 72 characters or less**
* Clear - Phrase your messages in a clear and grammatically-correct language.
* Concise - Brief, but comprehensive.
* Present Tense - Use the present tense ("Add feature" not "Added feature")
* Imperative - Use the imperative mood ("Move cursor to..." not "Moves cursor to...")

### Create a Pull Request
If you're unfamiliar with open-source contributions on GitHub, follow the [Creating a pull request](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request) guide.
 
### Submit a Pull Request

When you're finished with the changes, create a pull request, also known as a PR.

-   Fill the "Ready for review" template so that we can review your PR. This template helps reviewers understand your changes as well as the purpose of your pull request.
-   Don't forget to  [link PR to issue](https://docs.github.com/en/issues/tracking-your-work-with-issues/linking-a-pull-request-to-an-issue)  if you are solving one.
-   Enable the checkbox to  [allow maintainer edits](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/allowing-changes-to-a-pull-request-branch-created-from-a-fork)  so the branch can be updated for a merge. Once you submit your PR, a Docs team member will review your proposal. We may ask questions or request for additional information.
-   We may ask for changes to be made before a PR can be merged, either using  [suggested changes](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/incorporating-feedback-in-your-pull-request)  or pull request comments. You can apply suggested changes directly through the UI. You can make any other changes in your fork, then commit them to your branch.
-   As you update your PR and apply changes, mark each conversation as  [resolved](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/commenting-on-a-pull-request#resolving-conversations).
-   If you run into any merge issues, checkout this  [git tutorial](https://lab.github.com/githubtraining/managing-merge-conflicts)  to help you resolve merge conflicts and other issues.

### Your PR is merged!
Thank you for helping to improve Legitify!
Now that you're officially part of the community, your name will be publicly visible in the Changelog!

### Special Thanks
This Contributing guide was inspired by GitHub Docs's Contributing guide.
