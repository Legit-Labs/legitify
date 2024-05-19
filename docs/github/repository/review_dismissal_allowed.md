---
layout: default
title: Default Branch Should Restrict Who Can Dismiss Reviews
parent: Repository Policies
grand_parent: GitHub Policies
---


## Default Branch Should Restrict Who Can Dismiss Reviews
policy name: review_dismissal_allowed

severity: LOW

### Description
Any user with write access to the repository can dismiss pull-request reviews. Pull-request review contains essential information on the work that needs to be done and helps keep track of the changes. Dismissing it might cause a loss of this information and should be restricted to a limited number of users.

### Threat Example(s)
Allowing the dismissal of reviews can promote poor and vulnerable code, as important comments may be forgotten and ignored during the review process.



### Remediation
Note: The remediation steps apply to legacy branch protections, rules set-based protection should be updated from the rules set page
1. Make sure you have admin permissions
2. Go to the repo's settings page
3. Enter 'Branches' tab
4. Under 'Branch protection rules'
5. Click 'Edit' on the default branch rule
6. Check 'Restrict who can dismiss pull request reviews'
7. Click 'Save changes'



