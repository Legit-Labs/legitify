---
layout: default
title: Webhook Configured Without SSL
parent: Repository Policies
grand_parent: Policies
---


## Webhook Configured Without SSL
policy name: repository_webhook_doesnt_require_ssl

severity: LOW

### Description
Webhooks that are not configured with SSL enabled could expose your sofware to man in the middle attacks (MITM).


### Remediation
1. Make sure you can manage webhooks for the repository
2. Go to the repository settings page
3. Select "Webhooks"
4. Verify url starts with https
5. Press on the insecure webhook
6. Enable "SSL verfication"
7. Click "Update webhook"



