---
layout: default
title: Webhook Configured Without A Secret
parent: Organization Policies
grand_parent: GitHub Policies
---


## Webhook Configured Without A Secret
policy name: organization_webhook_no_secret

severity: LOW

### Description
Webhooks that are not configured with a token authenticated to validate the origin of the request and could make your software vulnerable.


### Remediation
1. Make sure you have admin permissions
2. Go to the organization settings page
3. Select "Webhooks"
4. Press on the insecure webhook
5. Configure a secret
6. Click "Update webhook"



