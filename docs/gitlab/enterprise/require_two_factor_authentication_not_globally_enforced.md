---
layout: default
title: Two-factor Authentication Should Be Globally Enforced
parent: Enterprise Policies
grand_parent: GitLab Policies
---


## Two-factor Authentication Should Be Globally Enforced
policy name: require_two_factor_authentication_not_globally_enforced

severity: HIGH

### Description
It is recommended to turn on MFA at the server or account level, and proactively prevent any new user created without MFA. Even if identities are managed externally using an SSO, it is highly recommended to keep this option on, for the 'admin' user or to be protected in the future for a deliberate or incidental creation of a user without MFA.



### Remediation
1. Press Settings -> General
2. Expand "Sign-in restrictions" section
3. Toggle "Two-factor authentication"
4. Press "Save Changes"



