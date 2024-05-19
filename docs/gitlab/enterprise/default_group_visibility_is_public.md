---
layout: default
title: Default Group Visibility Should Not Be Public
parent: Enterprise Policies
grand_parent: GitLab Policies
---


## Default Group Visibility Should Not Be Public
policy name: default_group_visibility_is_public

severity: HIGH

### Description
Your server is configured by default that every new group is created public. This means that it will be accessible by any member and unauthenticated user. It is recommended to turn the default option to be private and let admins control public visibility explicitly when needed.



### Remediation
2. Press Settings -> General
3. Expand 'Visibility and access controls' section
4. Under 'Default group visibility' toggle 'Private'
5. Press 'Save Changes'



