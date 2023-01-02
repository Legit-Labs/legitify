---
layout: default
title: Default Branch Doesn’t Require All Checks To Pass Before Merge
parent: Repository Policies
grand_parent: GitHub Policies
---


## Default Branch Doesn’t Require All Checks To Pass Before Merge
policy name: requires_status_checks

severity: MEDIUM

### Description
Branch protection is enabled. However, the checks which validate the quality and security of the code are not required to pass before submitting new changes. The default check ensures code is up-to-date in order to prevent faulty merges and unexpected behaviors, as well as other custom checks that test security and quality. It is advised to turn this control on to ensure any existing or future check will be required to pass. This option is found in the branch protection setting for the repository.


### Remediation
1. Make sure you have admin permissions
2. Go to the repo's settings page
3. Enter "Branches" tab
4. Under "Branch protection rules"
5. Click "Edit" on the default branch rule
6. Check "Require status checks to pass before merging"
7. Add the required checks that must pass before merging (tests, lint, etc...)
8. Click "Save changes"



