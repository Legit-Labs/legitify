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
#     - "1. Prerequisites: the attacker is part of your GitHub organization"
#     - "2. Create new repository in the organization"
#     - "3. Create a workflow file that reads all organization secrets and exfiltrate them"
#     - "4. Trigger the workflow"
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
#     - "1. Attacker creates a repository with a malicious custom GitHub Action"
#     - "2. Developer / Devops engineer uses this malicious action"
#     - "3. The malicious action has access to the developer repository and could steal its secrets or modify its content"
default all_github_actions_are_allowed  = false
all_github_actions_are_allowed {
    input.actions_permissions.allowed_actions == "all"
}