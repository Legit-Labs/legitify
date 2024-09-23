---
layout: default
title: Overriding predefined CI/CD variables should be restricted.
parent: Project Policies
grand_parent: GitLab Policies
---


## Overriding predefined CI/CD variables should be restricted.
policy name: overriding_defined_variables_isnt_restricted

severity: LOW

### Description
It’s recommended to restrict users with low privileges from overriding predefined variables, as doing so could compromise the security and integrity of the CI/CD pipeline.

### Threat Example(s)
Allowing overrides of predefined variables can result in unintentional misconfigurations of the CI/CD pipeline or deliberate tampering.



### Remediation
1. Make sure you have owner or maintainer permissions
2. The remediation is available through the project's API (e.g., 'https://gitlab.com/api/v4/projects/<your-project-ID>')
3. Set the 'restrict_user_defined_variables' attribute to TRUE (this attribute is FALSE by default)



