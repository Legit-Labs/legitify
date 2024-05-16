package actions

# METADATA
# scope: rule
# title: GitHub Actions Should Be Restricted To Selected Repositories
# description: By not limiting GitHub Actions to specific repositories, every user in the organization is able to run arbitrary workflows. This could enable malicious activity such as accessing organization secrets, crypto-mining, etc.
# custom:
#   requiredEnrichers: [organizationId]
#   remediationSteps: 
#     - 1. Make sure you have admin permissions
#     - 2. Go to the org's settings page
#     - 3. Enter the 'Actions - General' tab
#     - 4. Under 'Policies', Change 'All repositories' to 'Selected repositories' and select repositories that should be able to run actions
#     - 5. Click 'Save'
#   severity: MEDIUM
#   requiredScopes: [admin:org]
#   threat: 
#     - This misconfiguration could lead to the following attack:
#     - 1. Prerequisite: the attacker is part of your GitHub organization
#     - 2. Attacker creates new repository in the organization
#     - 3. Attacker creates a workflow file that reads all organization secrets and exfiltrate them
#     - 4. Attacker trigger the workflow
#     - 5. Attacker receives all organization secrets and uses them maliciously
#   requiredScopes: [admin:org]
default all_repositories_can_run_github_actions := true

all_repositories_can_run_github_actions := false {
	input.actions_permissions.enabled_repositories != "all"
}

# METADATA
# scope: rule
# title: GitHub Actions Should Be Limited To Verified or Explicitly Trusted Actions
# description: It is recommended to only use GitHub Actions by Marketplace verified creators or explicitly trusted actions. By not restricting which actions are permitted, developers may use actions that were not audited and may be malicious, thus exposing your pipeline to supply chain attacks.
# custom:
#   requiredEnrichers: [organizationId]
#   remediationSteps: 
#     - 1. Make sure you have admin permissions
#     - 2. Go to the org's settings page
#     - 3. Enter 'Actions - General' tab
#     - 4. Under 'Policies', Select 'Allow enterprise, and select non-enterprise, actions and reusable workflows'
#     - 5. Check 'Allow actions created by GitHub' and 'Allow actions by Marketplace verified creators'
#     - 6. Set any other used trusted actions under 'Allow specified actions and reusable workflows'
#     - 7. Click 'Save'
#   severity: MEDIUM
#   requiredScopes: [admin:org]
#   threat:
#     - This misconfiguration could lead to the following attack:
#     - 1. Attacker creates a repository with a tempting but malicious custom GitHub Action
#     - 2. An innocent developer / DevOps engineer uses this malicious action
#     - 3. The malicious action has access to the developer repository and could steal its secrets or modify its content
default all_github_actions_are_allowed := true

all_github_actions_are_allowed := false {
	input.actions_permissions.allowed_actions != "all"
}

# METADATA
# scope: rule
# title: Default Workflow Token Permission Should Be Read Only
# description: The default GitHub Action workflow token permission is set to read-write. When creating workflow tokens, it is highly recommended to follow the Principle of Least Privilege and force workflow authors to specify explicitly which permissions they need.
# custom:
#   requiredEnrichers: [organizationId]
#   remediationSteps:
#     - 1. Make sure you have admin permissions
#     - 2. Go to the org's settings page
#     - 3. Enter 'Actions - General' tab
#     - 4. Under 'Workflow permissions'
#     - 5. Select 'Read repository contents permission'
#     - 6. Click 'Save'
#   severity: MEDIUM
#   requiredScopes: [admin:org]
#   threat: In case of token compromise (due to a vulnerability or malicious third-party GitHub actions), an attacker can use this token to sabotage various assets in your CI/CD pipeline, such as packages, pull-requests, deployments, and more.
default token_default_permissions_is_read_write := true

token_default_permissions_is_read_write := false {
	input.token_permissions.default_workflow_permissions == "read"
}

# METADATA
# scope: rule
# title: Workflows Should Not Be Allowed To Approve Pull Requests
# description: Your organization allows workflows to open and approve pull requests. This could allow users to bypass code-review restrictions.
# custom:
#   requiredEnrichers: [organizationId]
#   remediationSteps:
#     - 1. Make sure you have admin permissions
#     - 2. Go to the org's settings page
#     - 3. Enter 'Actions - General' tab
#     - 4. Under 'Workflow permissions'
#     - 5. Uncheck 'Allow GitHub actions to create and approve pull requests'
#     - 6. Click 'Save'
#   severity: HIGH
#   requiredScopes: [admin:org]
#   threat: Attackers can exploit this misconfiguration to bypass code-review restrictions by creating a workflow that approves their own pull request and then merging the pull request without anyone noticing, introducing malicious code that would go straight ahead to production.
default actions_can_approve_pull_requests := true

actions_can_approve_pull_requests := false {
	not input.token_permissions.can_approve_pull_request_reviews
}
