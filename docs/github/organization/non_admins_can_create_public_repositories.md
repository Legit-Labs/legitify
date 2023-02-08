---
layout: default
title: Only Admins Should Be Able To Create Public Repositories
parent: Organization Policies
grand_parent: GitHub Policies
---


## Only Admins Should Be Able To Create Public Repositories
policy name: non_admins_can_create_public_repositories

severity: MEDIUM

### Description
The organization should be configured to prevent non-admin members creating public repositories. Creating a public repository may expose sensitive organization code, which, once exposed, may be copied, cached or stored by external parties. Therefore, it is highly recommended to restrict the option to create public repositories to admins only and reduce the risk of unintentional code exposure.

### Threat Example(s)
A member of the organization could inadvertently or maliciously make public an internal repository exposing confidential data.



### Remediation
1. Make sure you have admin permissions
2. Go to the organization settings page
3. Enter "Member privileges" tab
4. Under "Repository creation"
5. Toggle off "Public"
6. Click "Save"



