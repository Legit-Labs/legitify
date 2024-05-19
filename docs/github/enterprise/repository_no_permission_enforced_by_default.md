---
layout: default
title: Enterprise Should Define Base Permissions As 'No Permission' For All Members
parent: Enterprise Policies
grand_parent: GitHub Policies
---


## Enterprise Should Define Base Permissions As 'No Permission' For All Members
policy name: repository_no_permission_enforced_by_default

severity: MEDIUM

### Description
Collaborators in your organizations should receive access to specific organizations and repositories as necessary, and not have read and write access to all repositories across the enterprise.

### Threat Example(s)
An adversary will have access to all repositories in the enterprise, instead of just a part of them.



### Remediation
1. Make sure you are an enterprise owner
2. Go to the Settings page
3. Under the 'Policies' tab, choose 'Repositories'
4. Under 'Base Permission' choose 'No Permission'



