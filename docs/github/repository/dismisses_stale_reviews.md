---
layout: default
title: Default Branch Should Require New Code Changes After Approval To Be Re-Approved
parent: Repository Policies
grand_parent: GitHub Policies
---


## Default Branch Should Require New Code Changes After Approval To Be Re-Approved
policy name: dismisses_stale_reviews

severity: LOW

### Description
This security control prevents merging code that was approved but later on changed. Turning it on ensures any new changes must be reviewed again. This setting is part of the branch protection and code-review settings, and hardens the review process. If turned off - a developer can change the code after approval, and push code that is different from the one that was previously allowed. This option is found in the branch protection setting for the repository.

### Threat Example(s)
Buggy or insecure code may be committed after approval and will reach the main branch without review. Alternatively, an attacker can attempt a just-in-time attack to introduce dangerous code just before merge.



### Remediation
1. Note: The remediation steps applys to legacy branch protections, rules set based protection should be updated from the rules set page
2. Make sure you have admin permissions
3. Go to the repo's settings page
4. Enter "Branches" tab
5. Under "Branch protection rules"
6. Click "Edit" on the default branch rule
7. Check "Require a pull request before merging"
8. Check "Dismiss stale pull request approvals when new commits are pushed"
9. Click "Save changes"



