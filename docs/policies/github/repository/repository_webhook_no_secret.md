---
layout: default
title: Webhook Configured Without A Secret
parent: Repository Policies
grand_parent: GitHub Policies
---


## Webhook Configured Without A Secret
policy name: repository_webhook_no_secret

severity: LOW

### Description
Webhooks that are not configured with a token authenticated to validate the origin of the request and could make your software vulnerable.


### Remediation
1. Make sure you can manage webhooks for the repository
2. Go to the repository settings page
3. Select "Webhooks"
4. Press on the insecure webhook
5. Confiure a secret
6. Click "Update webhook"



