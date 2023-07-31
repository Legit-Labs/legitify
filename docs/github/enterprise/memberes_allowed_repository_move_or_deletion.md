---
layout: default
title: Enterprise Should Prevent Repository Admins From Deleting Or Transferring Repositories
parent: Enterprise Policies
grand_parent: GitHub Policies
---


## Enterprise Should Prevent Repository Admins From Deleting Or Transferring Repositories
policy name: memberes_allowed_repository_move_or_deletion

severity: MEDIUM

### Description
The enterprise’s Repository deletion and transfer policy should be set to DISABLED. This will prevent repository admins from deleting a repo or transferring it to a different owner or organization. Malicious actors could leak code if enabled.

### Threat Example(s)
A member of the organization could inadvertently or maliciously transfer a repository to an external namespace and expose confidential data.



### Remediation
1. Make sure you are an enterprise owner
2. Go to the Enterprise Settings page
3. Under the ‘Policies’ tab choose ‘Repositories’
4. Go to the ‘Admin repository permissions' section
5. under ‘Repository deletion and transfer' and select 'Disabled'



