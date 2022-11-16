package actions

# METADATA
# scope: rule
# title: GitHub Actions Is Not Restricted To Selected Repositories
# description: By not limiting GitHub Actions to specific repositories, every user in the organization is able to run arbitrary workflows. This could enable malicious activity such as accessing organization secrets, crypto-mining, etc.
# custom:
#   requiredEnrichers: [organizationId]
#   remediationSteps: [Make sure you have admin permissions, Go to the org's settings page, Enter the "Actions - General" tab, Under "Policies", Change "All repositories" to "Selected repositories" and select repositories that should be able to run actions, Click "Save"]
#   severity: MEDIUM
#   requiredScopes: [admin:org]
#   threat:
#     - "This misconfiguration could lead to the following attack:"
#     - "1. Prerequisite: the attacker is part of your GitHub organization"
#     - "2. Attacker creates new repository in the organization"
#     - "3. Attacker creates a workflow file that reads all organization secrets and exfiltrate them"
#     - "4. Attacker trigger the workflow"
#     - "5. Attacker receives all organization secrets and uses them maliciously"
#   requiredScopes: [admin:org]
default all_repositories_can_run_github_actions = false
all_repositories_can_run_github_actions {
    input.actions_permissions.enabled_repositories == "all"
}

# METADATA
# scope: rule
# title: GitHub Actions Runs Are Not Limited To Verified Actions
# description: When using GitHub Actions, it is recommended to only use actions by Marketplace verified creators or explicitly trusted actions. By not restricting which actions are permitted allows your developers to use actions that were not audited and potentially malicious, thus exposing your pipeline to supply chain attacks.
# custom:
#   requiredEnrichers: [organizationId]
#   remediationSteps: [Make sure you have admin permissions, Go to the org's settings page, Enter "Actions - General" tab, Under "Policies", 'Select "Allow enterprise, and select non-enterprise, actions and reusable workflows"', Check "Allow actions created by GitHub" and "Allow actions by Marketplace verified creators", Set any other used trusted actions under "Allow specified actions and reusable workflows", Click "Save"]
#   severity: MEDIUM
#   requiredScopes: [admin:org]
#   threat:
#     - "This misconfiguration could lead to the following attack:"
#     - "1. Attacker creates a repository with a tempting but malicious custom GitHub Action"
#     - "2. An innocent developer / DevOps engineer uses this malicious action"
#     - "3. The malicious action has access to the developer repository and could steal its secrets or modify its content"
default all_github_actions_are_allowed  = false
all_github_actions_are_allowed {
    input.actions_permissions.allowed_actions == "all"
}

# METADATA
# scope: rule
# title: Default workflow token permission is not read only
# description: Your default GitHub Action workflow token permission is set to read-write. When creating workflow tokens, it is highly recommended to follow the Principle of Least Privilege and force workflow authors to specify explicitly which permissions they need.
# custom:
#   requiredEnrichers: [organizationId]
#   remediationSteps:
#     - 1. Make sure you have admin permissions
#     - 2. Go to the org's settings page
#     - 3. Enter "Actions - General" tab
#     - 4. Under 'Workflow permissions'
#     - 5. Select 'Read repository contents permission'
#     - 6. Click 'Save'
#   severity: MEDIUM
#   requiredScopes: [admin:org]
#   threat: In case of token compromise (due to a vulnerability or malicious third-party GitHub actions) attacker get access to a privileged access token.
default token_default_permissions_is_read_write  = false
token_default_permissions_is_read_write {
    input.token_permissions.default_workflow_permissions != "read"
}

# METADATA
# scope: rule
# title: Workflows Are Allowed To Approve Pull Requests
# description: Your default GitHub Actions configuration allows for workflows to approve pull requests. This could allows users to bypass code-review restrictions.
# custom:
#   requiredEnrichers: [organizationId]
#   remediationSteps:
#     - 1. Make sure you have admin permissions
#     - 2. Go to the org's settings page
#     - 3. Enter "Actions - General" tab
#     - 4. Under 'Workflow permissions'
#     - 5. Uncheck 'Allow GitHub actions to create and approve pull requests.
#     - 6. Click 'Save'
#   severity: MEDIUM
#   requiredScopes: [admin:org]
#   threat: Attackers can exploit this misconfiguration to bypass code-review restriction by creating a workflow that approve thier own pull request and then merge the pull request without anyone noticing.
default actions_can_approve_pull_requests  = false
actions_can_approve_pull_requests {
    input.token_permissions.can_approve_pull_request_reviews
}