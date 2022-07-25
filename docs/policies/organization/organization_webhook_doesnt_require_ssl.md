---
layout: default
title: Webhook Configured Without SSL
parent: Organization Policies
grand_parent: Policies
---


## Webhook Configured Without SSL
policy name: organization_webhook_doesnt_require_ssl

severity: LOW

### Description
Webhooks that are not configured with SSL enabled could expose your sofware to man in the middle attacks (MITM).


### Remediation
1. Make sure you have admin permissions
2. Go to the organization settings page
3. Select "Webhooks"
4. Press on the insecure webhook
5. Verify url starts with https
6. Enable "SSL verfication"
7. Click "Update webhook"



