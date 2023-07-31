---
layout: default
title: Default Branch Should Not Allow Force Pushes
parent: Repository Policies
grand_parent: GitHub Policies
---


## Default Branch Should Not Allow Force Pushes
policy name: missing_default_branch_protection_force_push

severity: MEDIUM

### Description
The history of the default branch is not protected against changes for this repository. Protecting branch history ensures every change that was made to code can be retained and later examined. This issue is raised if the default branch history can be modified using force push.

### Threat Example(s)
Rewriting project history can make it difficult to trace back when bugs or security issues were introduced, making them more difficult to remediate.



### Remediation
1. Note: The remediation steps applys to legacy branch protections, rules set based protection should be updated from the rules set page
2. Make sure you have admin permissions
3. Go to the repo's settings page
4. Enter "Branches" tab
5. Under "Branch protection rules"
6. Click "Edit" on the default branch rule
7. Uncheck "Allow force pushes"
8. Click "Save changes"



