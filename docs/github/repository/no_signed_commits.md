---
layout: default
title: Default Branch Should Require All Commits To Be Signed
parent: Repository Policies
grand_parent: GitHub Policies
---


## Default Branch Should Require All Commits To Be Signed
policy name: no_signed_commits

severity: LOW

### Description
Require all commits to be signed and verified

### Threat Example(s)
A commit containing malicious code may be crafted by a malicious actor that has acquired write access to the repository to initiate a supply chain attack. Commit signing provides another layer of defense that can prevent this type of compromise.



### Remediation
1. Note: The remediation steps applys to legacy branch protections, rules set based protection should be updated from the rules set page
2. Make sure you have admin permissions
3. Go to the repo's settings page
4. Enter "Branches" tab
5. Under "Branch protection rules"
6. Click "Edit" on the default branch rule
7. Check "Require signed commits"
8. Click "Save changes"



