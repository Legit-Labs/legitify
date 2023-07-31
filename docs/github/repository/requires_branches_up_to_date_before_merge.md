---
layout: default
title: Default Branch Should Require Branches To Be Up To Date Before Merge
parent: Repository Policies
grand_parent: GitHub Policies
---


## Default Branch Should Require Branches To Be Up To Date Before Merge
policy name: requires_branches_up_to_date_before_merge

severity: MEDIUM

### Description
Status checks are required, but branches that are not up to date can be merged. This can result in previously remediated issues being merged in over fixes.

### Threat Example(s)
Required status checks may be failing on the latest version after passing on an earlier version of the code, making it easy to commit buggy or otherwise insecure code.



### Remediation
1. Note: The remediation steps applys to legacy branch protections, rules set based protection should be updated from the rules set page
2. Make sure you have admin permissions
3. Go to the repo's settings page
4. Enter "Branches" tab
5. Under "Branch protection rules"
6. Click "Edit" on the default branch rule
7. Check "Require status checks to pass before merging"
8. Check "Require branches to be up to date before merging"
9. Click "Save changes"



