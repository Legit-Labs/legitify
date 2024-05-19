---
layout: default
title: Two-Factor Authentication Grace Period Should Not Be Longer Than One Week
parent: Group Policies
grand_parent: GitLab Policies
---


## Two-Factor Authentication Grace Period Should Not Be Longer Than One Week
policy name: group_allows_excessive_mfa_grace_period

severity: MEDIUM

### Description
New members added to your group are allowed longer than a week to enable MFA. The time frame should be lowered to one week or less.

### Threat Example(s)
Any new group member effectively acts as an attack surface until two-factor authentication is enabled. The risk is compounded as new members may be more vulnerable to phishing and identity theft attacks.



### Remediation
1. Go to the group page
2. Press Settings -> General
3. Expand 'Permissions and group features'
4. In the box titled: 'Delay 2FA enforcement (hours)', enter a number under 168 (preferably 0)
5. Press 'Save Changes'



