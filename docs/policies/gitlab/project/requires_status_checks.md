---
layout: default
title: Project Doesn’t Require All Pipelines to Succeed
parent: Project Policies
grand_parent: GitLab Policies
---


## Project Doesn’t Require All Pipelines to Succeed
policy name: requires_status_checks

severity: MEDIUM

### Description
Checks that validate the quality and security of the code are not required to pass before submitting new changes. It is advised to turn this flag on to ensure any existing or future check will be required to pass.


### Remediation
1. Make sure you can manage project merge requests permissions
2. Go to the project's settings page
3. Select "Merge Requests"
4. Press on the "Pipelines must succeed"
5. Click "Save changes"



