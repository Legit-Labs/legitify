---
layout: default
title: Default Branch Should Require Code Review
parent: Project Policies
grand_parent: GitLab Policies
---


## Default Branch Should Require Code Review
policy name: code_review_not_required

severity: HIGH

### Description
In order to comply with separation of duties principle and enforce secure code practices, a code review should be mandatory using the source-code-management built-in enforcement.

### Threat Example(s)
Users can merge code without being reviewed which can lead to insecure code reaching the main branch and production.



### Remediation
1. Make sure you have admin permissions
2. Go to the repo's settings page
3. Enter 'Merge Requests' tab
4. Under 'Merge request approvals', Click 'Add approval rule' on the default branch rule
5. Select 'Approvals required' and enter at least 1 approvers
6. Select 'Add approvers' and select the desired members
7. Click 'Add approval rule'



