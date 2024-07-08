---
layout: default
title: Repository Should Have A Low Admin Count
parent: Repository Policies
grand_parent: GitHub Policies
---


## Repository Should Have A Low Admin Count
policy name: repository_has_too_many_admins

severity: LOW

### Description
Repository admins are highly privileged and could create great damage if they are compromised. It is recommended to limit the number of repository admins to the minimum required, and no more than 5% of the userbase (Up to 3 admins are always allowed).

### Threat Example(s)
A compromised user with admin permissions can initiate a supply chain attack in a plethora of ways.
Having many admin users increases the overall risk of user compromise, and makes it more likely to lose track of unused admin permissions given to users in the past.



### Remediation
1. Make sure you have admin permissions
2. Go to the repository settings page
3. Press 'Collaborators and teams'
4. Select the unwanted admin users
5. Select 'Change Role'



