---
layout: default
title: Webhooks Should Not Be Allowed To Be Sent To The Local Network
parent: Enterprise Policies
grand_parent: GitLab Policies
---


## Webhooks Should Not Be Allowed To Be Sent To The Local Network
policy name: webhooks_are_allowed_to_be_sent_to_local_network

severity: LOW

### Description
Webhooks sent by GitLab servers are authenticated, and can cause potential damage if sent uncontrollably. For example, a malicious user can plant a webhook address that triggers an API call in GitLab itself which can delete resources. Therefore, as a security best practice, webhooks should be limited to external URLs only. You can read more [here](https://docs.gitlab.com/ee/security/webhooks.html)



### Remediation
1. Press Settings -> Network
2. Expand "Outbound requests" section
3. Un toggle "Allow requests to the local network from web hooks and services"
4. Press "Save Changes"



