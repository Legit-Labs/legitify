---
layout: default
title: Enterprise Should Prevent Members From Creating Public Repositories
parent: Enterprise Policies
grand_parent: GitHub Policies
---


## Enterprise Should Prevent Members From Creating Public Repositories
policy name: enterprise_allows_creating_public_repos

severity: MEDIUM

### Description
The enterprise's repository creation policy should be set to private/internal repositories only. This will prevent non-admin users from creating public repositories and potentially exposing source code.

### Threat Example(s)
Users can accidentally create public repositories and expose source code.



### Remediation
1. Make sure you are an enterprise owner
2. Go to the policies page
3. Under the 'Repository creation' section, Choose the 'Members can create repositories' option and uncheck 'Public'



