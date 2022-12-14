---
layout: default
title: Project Doesn’t Require All Pipelines to Succeed
parent: Repository Policies
grand_parent: GitLab Policies
---


## Project Doesn’t Require All Pipelines to Succeed
policy name: requires_status_checks

severity: MEDIUM

### Description
the checks which validate the quality and security of the code are not required to pass before submitting new changes. It is advised to turn this control on to ensure any existing or future check will be required to pass

### Threat Example(s)
Users could merge its code without all required checks passes what could lead to insecure code reaching your main branch and production.



### Remediation
1. Make sure you can manage project merge requests permissions
2. Go to the project's settings page
3. Select "Merge Requests"
4. Press on the "Pipelines must succeed"
5. Click "Save changes"



