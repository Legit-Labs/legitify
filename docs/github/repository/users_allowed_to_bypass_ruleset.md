---
layout: default
title: Users Are Allowed To Bypass Ruleset Rules
parent: Repository Policies
grand_parent: GitHub Policies
---


## Users Are Allowed To Bypass Ruleset Rules
policy name: users_allowed_to_bypass_ruleset

severity: MEDIUM

### Description
Rulesets rules are not enforced for some users. When defining rulesets it is recommended to make sure that no one is allowed to bypass these rules in order to avoid inadvertent or intentional alterations to critical code which can lead to potential errors or vulnerabilities in the software.

### Threat Example(s)
Attackers that gain access to a user that can bypass the ruleset rules can compromise the codebase without anyone noticing, introducing malicious code that would go straight ahead to production.



### Remediation
1. Go to the repository settings page
2. Under 'Code and automation', select 'Rules -> Rulesets'
3. Find the relevant ruleset
4. Empty the 'Bypass list'
5. Press 'Save Changes'



