---
layout: default
title: Forking Should Not Be Allowed for Private/Internal Projects
parent: Project Policies
grand_parent: GitLab Policies
---


## Forking Should Not Be Allowed for Private/Internal Projects
policy name: forking_allowed_for_repository

severity: LOW

### Description
Forking a project (repository) can lead to loss of control and potential exposure of source code. If you do not need forking, it is recommended to turn it off in the project or group configuration. The option to fork should be enabled only by owners deliberately when opting to create a fork.

### Threat Example(s)
Forked repositories may leak important code assets or sensitive secrets embedded in the code to anyone outside your organization, as the code becomes publicly accessible.



### Remediation
1. Make sure you have owner permissions
2. Go to the project's settings page
3. Enter 'General' tab
4. Under 'Visibility, project features, permissions', Toggle off 'Forks'



