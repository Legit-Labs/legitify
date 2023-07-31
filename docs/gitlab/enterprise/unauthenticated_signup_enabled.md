---
layout: default
title: Server Should Not Allow Access To Unauthenticated Users With Sign-Up
parent: Enterprise Policies
grand_parent: GitLab Policies
---


## Server Should Not Allow Access To Unauthenticated Users With Sign-Up
policy name: unauthenticated_signup_enabled

severity: HIGH

### Description
The server allows any person with network access to sign up, create a user and access sensitive data. Turning this off will reduce the risk of attackers trying to infiltrate the server.


### Remediation
1. Press Settings -> General
2. Expand "Sign-up restrictions" section
3. Un toggle "Sign-up enabled"
4. Press "Save Changes"



