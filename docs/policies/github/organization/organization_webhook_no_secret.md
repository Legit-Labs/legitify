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
Webhooks are not configured with an authenticated token to validate the origin of the request and could make your software vulnerable.

### Threat Example(s)
Not using a webhook secret makes the service receiving the webhook unable to determine the authenticity of the request.
This allows attackers to masquerade as your organization, potentially creating an unstable or insecure state in other systems.



### Remediation
1. Make sure you have admin permissions
2. Go to the organization settings page
3. Select "Webhooks"
4. Press on the insecure webhook
5. Configure a secret
6. Click "Update webhook"



