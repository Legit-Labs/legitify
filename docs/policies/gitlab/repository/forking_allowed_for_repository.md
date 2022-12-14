---
layout: default
title: Forking Allowed for This Repository
parent: Repository Policies
grand_parent: GitLab Policies
---


## Forking Allowed for This Repository
policy name: forking_allowed_for_repository

severity: LOW

### Description
Forking a repository can lead to loss of control and potential exposure of the source code. The option to fork must be disabled by default and turned on only by owners deliberately when opting to create a fork. If you do not need forking, it is recommended to turn it off in the project's configuration.

### Threat Example(s)
A user with permissions to the repository could intentionally/accidentally fork a private repository, make it public and cause a code-leak incident



### Remediation
1. Make sure you have owner permissions
2. Go to the project's settings page
3. Enter "General" tab
4. Under "Visibility
5. project features
6. permissions"
7. Toggle off "Forks"



