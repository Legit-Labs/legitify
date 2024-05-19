---
layout: default
title: Creating Public Repositories Should Be Restricted To Admins
parent: Enterprise Policies
grand_parent: GitLab Policies
---


## Creating Public Repositories Should Be Restricted To Admins
policy name: all_users_are_allowed_to_create_public_projects

severity: HIGH

### Description
Your server allows non-admin members to create public repositories. Creating a public repository can be done by mistake, and may expose sensitive organization code, that once exposed may be copied, cached or stored by external parties. Therefore, it is highly recommended to restrict the option to create public repositories to admins only and reduce the risk of unintentional code exposure.



### Remediation
2. Press Settings -> General
3. Expand 'Visibility and access controls' section
4. Under 'Restricted visibility levels' toggle 'Public'
5. Press 'Save Changes'



