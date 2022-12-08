---
layout: default
title: Workflows Are Allowed To Approve Pull Requests
parent: Actions Policies
grand_parent: GitHub Policies
---


## Workflows Are Allowed To Approve Pull Requests
policy name: actions_can_approve_pull_requests

severity: HIGH

### Description
Your default GitHub Actions configuration allows for workflows to approve pull requests. This could allow users to bypass code-review restrictions.


### Remediation
1. Make sure you have admin permissions
2. Go to the org's settings page
3. Enter "Actions - General" tab
4. Under 'Workflow permissions'
5. Uncheck 'Allow GitHub actions to create and approve pull requests.
6. Click 'Save'



