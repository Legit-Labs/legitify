---
layout: default
title: Webhook Configured Without SSL
parent: Organization Policies
grand_parent: GitLab Policies
---


## Webhook Configured Without SSL
policy name: organization_webhook_doesnt_require_ssl

severity: LOW

### Description
Webhooks that are not configured with SSL enabled could expose your software to man in the middle attacks (MITM).


### Remediation
1. Go to the group Settings -> Webhooks page
2. Find the misconfigured webhook and press "Edit"
3. Toggle "Enable SSL verification"
4. Press "Save Changes"



