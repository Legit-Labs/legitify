---
layout: default
title: Forking Should Not Be Allowed for Private/Internal Repositories
parent: Repository Policies
grand_parent: GitHub Policies
---


## Forking Should Not Be Allowed for Private/Internal Repositories
policy name: forking_allowed_for_repository

severity: LOW

### Description
Forking private or internal repositories can lead to unauthorized spread and potential exposure of sensitive source code. It is recommended to disable forking for private repositories in the repository or the organization configuration to maintain control over the source code. If forking is necessary, it should be enabled selectively by admins for specific collaboration needs on private repositories.

### Threat Example(s)
Forked repositories cause more code and secret sprawl in the organization as forks are independent copies of the repository and need to be tracked separately, making it more difficult to keep track of sensitive assets and contain potential incidents.



### Remediation
1. Make sure you have admin permissions
2. Go to the repo's settings page
3. Enter 'General' tab
4. Under 'Features', Toggle off 'Allow forking'



