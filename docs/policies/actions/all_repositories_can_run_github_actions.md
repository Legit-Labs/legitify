---
layout: default
title: GitHub Actions Is Not Restricted To Selected Repositories
parent: Actions Policies
grand_parent: Policies
---


## GitHub Actions Is Not Restricted To Selected Repositories
policy name: all_repositories_can_run_github_actions

severity: MEDIUM

### Description
By not limiting GitHub Actions to specific repositories, every user in the organization is able to run arbitrary workflows. This could enable malicious activity such as accessing organization secrets, crypto-mining, etc.

### Threat Example(s)
This misconfiguration could lead to the following attack:
1. Prerequisites: the attacker is part of your GitHub organization
2. Create new repository in the organization
3. Create a workflow file that reads all organization secrets and exfiltrate them
4. Trigger the workflow



### Remediation
1. Make sure you have admin permissions
2. Go to the org's settings page
3. Enter the "Actions - General" tab
4. Under "Policies"
5. Change "All repositories" to "Selected repositories" and select repositories that should be able to run actions
6. Click "Save"



