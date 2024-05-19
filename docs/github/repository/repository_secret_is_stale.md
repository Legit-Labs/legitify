---
layout: default
title: Repository Secrets Should Be Updated At Least Yearly
parent: Repository Policies
grand_parent: GitHub Policies
---


## Repository Secrets Should Be Updated At Least Yearly
policy name: repository_secret_is_stale

severity: MEDIUM

### Description
Some of the repository secrets have not been updated for over a year. It is recommended to refresh secret values regularly in order to minimize the risk of breach in case of an information leak.

### Threat Example(s)
Sensitive data may have been inadvertently made public in the past, and an attacker who holds this data may gain access to your current CI and services. In addition, there may be old or unnecessary tokens that have not been inspected and can be used to access sensitive information.



### Remediation
1. Enter your repository's landing page
2. Go to the settings tab
3. Under the 'Security' title on the left, choose 'Secrets and variables'
4. Click 'Actions'
5. Sort secrets by 'Last Updated'
6. Regenerate every secret older than one year and add the new value to GitHub's secret manager



