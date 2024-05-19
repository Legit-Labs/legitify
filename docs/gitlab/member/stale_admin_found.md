---
layout: default
title: Administrators Should Have Activity in the Last 6 Months
parent: Member Policies
grand_parent: GitLab Policies
---


## Administrators Should Have Activity in the Last 6 Months
policy name: stale_admin_found

severity: MEDIUM

### Description
A collaborator with global admin permissions didn't perform any action in the last 6 months. Admin users are extremely powerful, and common compliance standards demand keeping the number of admins to a minimum. Consider revoking this collaborator's admin credentials (downgrade to regular user) or removing the user completely.

### Threat Example(s)
Stale admins are most likely not managed and monitored, increasing the possibility of being compromised.



### Remediation
1. Go to the admin menu
2. Select 'Overview -> Users' on the left navigation bar
3. Find the stale admin and either delete or block it



