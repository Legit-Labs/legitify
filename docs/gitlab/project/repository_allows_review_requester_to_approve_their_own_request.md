---
layout: default
title: Repository Should Not Allow Review Requester To Approve Their Own Request
parent: Project Policies
grand_parent: GitLab Policies
---


## Repository Should Not Allow Review Requester To Approve Their Own Request
policy name: repository_allows_review_requester_to_approve_their_own_request

severity: MEDIUM

### Description
To comply with separation of duties and enforce secure code practices, the repository should prohibit pull request owners from approving their own changes.

### Threat Example(s)
Users can merge code without being reviewed which can lead to insecure code reaching the main branch and production.



### Remediation
1. Make sure you have admin permissions
2. Go to the repo's settings page
3. Enter "Merge Requests" tab
4. Under "Approval settings"
5. Check "Prevent approval by author"
6. Click "Save Changes"



