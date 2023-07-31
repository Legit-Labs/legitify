---
layout: default
title: Default Branch Should Require Linear History
parent: Repository Policies
grand_parent: GitHub Policies
---


## Default Branch Should Require Linear History
policy name: non_linear_history

severity: MEDIUM

### Description
Prevent merge commits from being pushed to protected branches.

### Threat Example(s)
Having a non-linear history makes it harder to reverse changes, making recovery from bugs and security risks slower and more difficult.



### Remediation
1. Note: The remediation steps applys to legacy branch protections, rules set based protection should be updated from the rules set page
2. Make sure you have admin permissions
3. Go to the repo's settings page
4. Enter "Branches" tab
5. Under "Branch protection rules"
6. Click "Edit" on the default branch rule
7. Check "Require linear history"
8. Click "Save changes"



