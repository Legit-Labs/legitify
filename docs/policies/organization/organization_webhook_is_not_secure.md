---
layout: default
title: Webhook configured insecurely
parent: Organization Policies
grand_parent: Policies
---


## Webhook configured insecurely
policy name: organization_webhook_is_not_secure

severity: MEDIUM

### Description
Webhooks that are not configured with a token and SSL verification cannot be authenticated to validate the origin of the request and could make the webhooks receiver service vulnerable.


### Remediation
1. Make sure you have admin permissions
2. Go to the organization settings page
3. Select "Webhooks"
4. Press on the insecure webhook
5. Verify url starts with https
6. Enable "SSL verfication" and confiure a secret
7. Click "Update webhook"



