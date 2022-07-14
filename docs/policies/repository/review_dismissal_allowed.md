---
layout: default
title: Default Branch Allows Dismissal Of Reviews
parent: Repository Policies
grand_parent: Policies
---


## Default Branch Allows Dismissal Of Reviews
policy name: review_dismissal_allowed

severity: LOW

### Description
Users might be able to dismiss comments in a review, ignoring Pull Request comments and silently bypassing the conversation. Restrict dismissal of pull request reviews.


### Remediation
1. Make sure you have admin permissions
2. Go to the repo's settings page
3. Enter "Branches" tab
4. Under "Branch protection rules"
5. Click "Edit" on the default branch rule
6. Check "Restrict who can dismiss pull request reviews"
7. Click "Save changes"



