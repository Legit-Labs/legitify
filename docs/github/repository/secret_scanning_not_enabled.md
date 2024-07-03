---
layout: default
title: Secret Scanning should be enabled
parent: Repository Policies
grand_parent: GitHub Policies
---


## Secret Scanning should be enabled
policy name: secret_scanning_not_enabled

severity: MEDIUM

### Description
Repository should have secret scanning enabled. Secret scanning helps prevent the exposure of sensitive information and ensures compliance.

### Threat Example(s)
Exposed secrets increases the risk of sensitive information such as API keys, passwords, and tokens being disclosed, leading to unauthorized access to systems and services, and data breaches.



### Remediation
1. Go to the repository settings page
2. Under the 'Security' title on the left, select 'Code security and analysis'
3. Under 'Secret scanning', click 'Enable'



