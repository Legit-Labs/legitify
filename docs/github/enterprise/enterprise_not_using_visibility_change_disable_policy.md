---
layout: default
title: Enterprise Should Prevent Repository Admins From Changing Repository Visibility
parent: Enterprise Policies
grand_parent: GitHub Policies
---


## Enterprise Should Prevent Repository Admins From Changing Repository Visibility
policy name: enterprise_not_using_visibility_change_disable_policy

severity: MEDIUM

### Description
The enterprise's Repository visibility change policy should be set to DISABLED. This will prevent users from creating private repositories and changing them to be public. Malicious actors could leak code if enabled.

### Threat Example(s)
A member of the organization could inadvertently or maliciously make public an internal repository exposing confidential data.



### Remediation
1. Make sure you are an enterprise owner
2. Go to the policies page
3. Under the 'Repository visibility change' section, choose the 'Disabled' option



