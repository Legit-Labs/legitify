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
Require all Pull Request conversations to be resolved before merging. Check this to avoid bypassing/missing a Pull Reuqest comment.

### Threat Example(s)
Allowing the merging of code without resolving all conversations can promote poor and vulnerable code, as important comments may be forgotten or deliberately ignored when the code is merged.



### Remediation
1. Note: The remediation steps applys to legacy branch protections, rules set based protection should be updated from the rules set page
2. Make sure you have admin permissions
3. Go to the repo's settings page
4. Enter "Branches" tab
5. Under "Branch protection rules"
6. Click "Edit" on the default branch rule
7. Check "Require conversation resolution before merging"
8. Click "Save changes"



