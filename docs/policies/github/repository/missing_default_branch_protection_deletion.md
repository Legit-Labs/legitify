---
layout: default
title: Default Branch Could Be Deleted
parent: Repository Policies
grand_parent: GitHub Policies
---


## Default Branch Could Be Deleted
policy name: missing_default_branch_protection_deletion

severity: MEDIUM

### Description
The history of the default branch is not protected against deletion for this repository.


### Remediation
1. Make sure you have admin permissions
2. Go to the repo's settings page
3. Enter "Branches" tab
4. Under "Branch protection rules"
5. Click "Edit" on the default branch rule
6. Uncheck "Allow deletions"
7. Click "Save changes"



