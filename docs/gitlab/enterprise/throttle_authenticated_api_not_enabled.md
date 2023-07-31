---
layout: default
title: API Request Rate Limit Should Be Limited
parent: Enterprise Policies
grand_parent: GitLab Policies
---


## API Request Rate Limit Should Be Limited
policy name: throttle_authenticated_api_not_enabled

severity: LOW

### Description
Enabling authenticated request rate-limit on APIs reduces volume of requests, and prevents incidents of reduced quality-of-service or denial-of-service because of unnecessary or malicious high request volume.



### Remediation
1. Press Settings -> Network
2. Expand "User and IP rate limit" section
3. Toggle "Enable authenticated API request rate limit
4. Press "Save Changes"



