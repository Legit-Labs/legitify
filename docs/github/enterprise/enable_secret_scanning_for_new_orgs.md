---
layout: default
title: Enterprise Should Automatically Enable Secret Scanning Across All Organizations/Repositories
parent: Enterprise Policies
grand_parent: GitHub Policies
---


## Enterprise Should Automatically Enable Secret Scanning Across All Organizations/Repositories
policy name: enable_secret_scanning_for_new_orgs

severity: MEDIUM

### Description
Enable GitHub Advanced Security secret scanning to alert on sensitive data that exists in your enterprise. Secrets shouldn’t be hard-coded in to your repositories as they will be retrievable by anyone with access to the repository.


### Remediation
1. Make sure you are an enterprise owner
2. Go to the Enterprise Settings page
3. Under the ‘Settings’ tab choose ‘Code security and analysis’
4. Check 'Automatically enable for new repositories with Advanced Security enabled'



