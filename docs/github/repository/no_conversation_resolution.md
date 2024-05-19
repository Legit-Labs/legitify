---
layout: default
title: Default Branch Should Require All Conversations To Be Resolved Before Merge
parent: Repository Policies
grand_parent: GitHub Policies
---


## Default Branch Should Require All Conversations To Be Resolved Before Merge
policy name: no_conversation_resolution

severity: LOW

### Description
Require all Pull Request conversations to be resolved before merging. Check this to avoid bypassing/missing a Pull Request comment.

### Threat Example(s)
Allowing the merging of code without resolving all conversations can promote poor and vulnerable code, as important comments may be forgotten or deliberately ignored when the code is merged.



### Remediation
Note: The remediation steps apply to legacy branch protections, rules set-based protection should be updated from the rules set page
1. Make sure you have admin permissions
2. Go to the repo's settings page
3. Enter 'Branches' tab
4. Under 'Branch protection rules'
5. Click 'Edit' on the default branch rule
6. Check 'Require conversation resolution before merging'
7. Click 'Save changes'



