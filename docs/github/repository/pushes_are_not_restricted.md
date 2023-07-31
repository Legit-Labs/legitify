---
layout: default
title: Default Branch Should Restrict Who Can Push To It
parent: Repository Policies
grand_parent: GitHub Policies
---


## Default Branch Should Restrict Who Can Push To It
policy name: pushes_are_not_restricted

severity: LOW

### Description
By default, commits can be pushed directly to protected branches without going through a Pull Request. Restrict who can push commits to protected branches so that commits can be added only via merges, which require Pull Request.

### Threat Example(s)
An attacker with write credentials may introduce vulnerabilities to your code without your knowledge. Alternatively, contributors may commit unsafe code that is buggy or easy to exploit that could have been caught using a review process.



### Remediation
1. Note: The remediation steps applys to legacy branch protections, rules set based protection should be updated from the rules set page
2. Make sure you have admin permissions
3. Go to the repo's settings page
4. Enter "Branches" tab
5. Under "Branch protection rules"
6. Click "Edit" on the default branch rule
7. Check "Restrict who can push to matching branches"
8. Choose who should be allowed to push
9. Click "Save changes"



