---
layout: default
title: Webhook configured insecurely
parent: Repository Policies
grand_parent: Policies
---


## Webhook configured insecurely
policy name: repository_webhook_is_not_secure

severity: HIGH

### Description
Webhooks that are not configured with a token and SSL verification cannot be authenticated to validate the origin of the request and could make your software vulnerable.


### Remediation
1. Make sure you can manage webhooks for the repository
2. Go to the repository settings page
3. Select "Webhooks"
4. Verify url starts with https
5. Press on the insecure webhook
6. Enable "SSL verfication" and confiure a secret
7. Click "Update webhook"



