---
layout: default
title: Webhook Configured Without SSL Verification
parent: Repository Policies
grand_parent: GitLab Policies
---


## Webhook Configured Without SSL Verification
policy name: repository_webhook_doesnt_require_ssl

severity: LOW

### Description
Webhooks that are not configured with SSL verification enabled could expose your sofware to man in the middle attacks (MITM).


### Remediation
1. Make sure you can manage webhooks for the repository
2. Go to the repository settings page
3. Select "Webhooks"
4. Press on the "Enable SSL verfication"
5. Click "Save changes"



