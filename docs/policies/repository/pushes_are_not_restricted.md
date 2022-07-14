---
layout: default
title: Default Branch Allows Pushes to Protected Branch
parent: Repository Policies
grand_parent: Policies
---


## Default Branch Allows Pushes to Protected Branch
policy name: pushes_are_not_restricted

severity: MEDIUM

### Description
By default, commits can be pushed directly to protected branches, without going through a Pull Request. Restrict pushes to protected branches so that commits can be added only via merges, which require Pull Request.


### Remediation
1. Make sure you have admin permissions
2. Go to the repo's settings page
3. Enter "Branches" tab
4. Under "Branch protection rules"
5. Click "Edit" on the default branch rule
6. Check "Restrict who can push to matching branches"
7. Click "Save changes"



