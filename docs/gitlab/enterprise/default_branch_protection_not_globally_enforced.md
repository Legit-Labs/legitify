---
layout: default
title: Branch Protection Should Be Globally Enabled By Default
parent: Enterprise Policies
grand_parent: GitLab Policies
---


## Branch Protection Should Be Globally Enabled By Default
policy name: default_branch_protection_not_globally_enforced

severity: MEDIUM

### Description
Branch protection is not enabled for the entire server or account by default. Any newly created repository will not be protected unless explicitly specified. Protecting branches ensures new code changes must go through a controlled merge process and it allows enforcement of code review and other security tests. It is recommended to turn it on by default.



### Remediation
2. Press Settings -> Repository
3. Expand 'Default Branch' section
4. Toggle 'Fully protected'
5. Press 'Save Changes'



