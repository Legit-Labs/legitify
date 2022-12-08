---
layout: default
title: Default Branch Doesn't Require Linear History
parent: Repository Policies
grand_parent: GitHub Policies
---


## Default Branch Doesn't Require Linear History
policy name: non_linear_history

severity: MEDIUM

### Description
Prevent merge commits from being pushed to protected branches.


### Remediation
1. Make sure you have admin permissions
2. Go to the repo's settings page
3. Enter "Branches" tab
4. Under "Branch protection rules"
5. Click "Edit" on the default branch rule
6. Check "Require linear history"
7. Click "Save changes"



