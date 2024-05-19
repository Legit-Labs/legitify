---
layout: default
title: Unauthenticated Requests Rate Limit Should Be Enabled
parent: Enterprise Policies
grand_parent: GitLab Policies
---


## Unauthenticated Requests Rate Limit Should Be Enabled
policy name: throttle_unauthenticated_request_not_enabled

severity: MEDIUM

### Description
The server allows restricting the limit of unauthenticated requests. It is recommended to turn it on as a security and reliability measure, and to reduce request volume. If an attacker tries accessing the system, this will reduce the risk of brute-force and Denial-of-service to the end users caused by high request rate.



### Remediation
2. Press Settings -> Network
3. Expand 'User and IP rate limit' section
4. Toggle 'Enable unauthenticated API request rate limit' and 'Enable unauthenticated web request rate limit'
5. Press 'Save Changes'



