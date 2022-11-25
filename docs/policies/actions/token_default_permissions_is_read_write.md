---
layout: default
title: Default workflow token permission is not read only
parent: Actions Policies
grand_parent: Policies
---


## Default workflow token permission is not read only
policy name: token_default_permissions_is_read_write

severity: MEDIUM

### Description
Your default GitHub Action workflow token permission is set to read-write. When creating workflow tokens, it is highly recommended to follow the Principle of Least Privilege and force workflow authors to specify explicitly which permissions they need.


### Remediation
1. Make sure you have admin permissions
2. Go to the org's settings page
3. Enter "Actions - General" tab
4. Under 'Workflow permissions'
5. Select 'Read repository contents permission'
6. Click 'Save'



