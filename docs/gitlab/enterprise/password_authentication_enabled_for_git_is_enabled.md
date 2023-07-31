---
layout: default
title: Password Authentication For Git Over HTTP(S) Should Not Be Enabled
parent: Enterprise Policies
grand_parent: GitLab Policies
---


## Password Authentication For Git Over HTTP(S) Should Not Be Enabled
policy name: password_authentication_enabled_for_git_is_enabled

severity: LOW

### Description
Password authentication for Git protocol operations (pull / push) is discouraged in favor of using personal access tokens and keys. From a security standpoint, it is recommended to disable password authentication completely.



### Remediation
1. Press Settings -> General
2. Expand "Sign-in restrictions" section
3. Un toggle "Allow password authentication for Git over HTTP(S)"
4. Press "Save Changes"



