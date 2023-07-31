---
layout: default
title: Default Branch Should Require All Checks To Pass Before Merge
parent: Repository Policies
grand_parent: GitHub Policies
---


## Default Branch Should Require All Checks To Pass Before Merge
policy name: requires_status_checks

severity: MEDIUM

### Description
Branch protection is enabled. However, the checks which validate the quality and security of the code are not required to pass before submitting new changes. The default check ensures code is up-to-date in order to prevent faulty merges and unexpected behaviors, as well as other custom checks that test security and quality. It is advised to turn this control on to ensure any existing or future check will be required to pass.

### Threat Example(s)
Not defining a set of required status checks can make it easy for contributors to introduce buggy or insecure code as manual review, whether mandated or optional, is the only line of defense.



### Remediation
1. Note: The remediation steps applys to legacy branch protections, rules set based protection should be updated from the rules set page
2. Make sure you have admin permissions
3. Go to the repo's settings page
4. Enter "Branches" tab
5. Under "Branch protection rules"
6. Click "Edit" on the default branch rule
7. Check "Require status checks to pass before merging"
8. Add the required checks that must pass before merging (tests, lint, etc...)
9. Click "Save changes"



