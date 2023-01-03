---
layout: default
title: Forking Allowed for This Repository
parent: Repository Policies
grand_parent: GitHub Policies
---


## Forking Allowed for This Repository
policy name: forking_allowed_for_repository

severity: LOW

### Description
Forking a repository can lead to loss of control and potential exposure of the source code. The option to fork must be disabled by default and turned on only by admins deliberately when opting to create a fork. If you do not need forking, it is recommended to turn it off in the repository configuration.

### Threat Example(s)
Forked repositories cause more code and secret sprawl in the organization as forks are independent copies of the repository and need to be tracked separately, making it more difficult to keep track of sensitive assets and contain potential incidents.



### Remediation
1. Make sure you have admin permissions
2. Go to the repo's settings page
3. Enter "General" tab
4. Under "Features"
5. Toggle off "Allow forking"



