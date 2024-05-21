package runner_group

# METADATA
# scope: rule
# title: Runner Group Should Be Limited to Private Repositories
# description: |
#       Workflows from public repositories are allowed to run on GitHub Hosted Runners.
#       When using GitHub Hosted Runners,
#       it is recommended to allow only workflows from private repositories to run on these runners.
#       to avoid being vulnerable to malicious actors using workflows from public repositories to break into your private network.
#       In case of inadequate security measures implemented on the hosted runner,
#       malicious actors could fork your repository and then create a pwn-request (a pull-request from a forked repository to the base repository with malicious intentions)
#       that create a workflow that exploits these vulnerabilities and move laterally inside your network.
# custom:
#   severity: HIGH
#   requiredEnrichers: [organizationId]
#   requiredScopes: [admin:org]
#   remediationSteps:
#     - 1. Go to the organization settings page
#     - 2. Press Actions ➝ Runner groups
#     - 3. Select the violating repository
#     - 4. Uncheck Allow public repositories
#   threat:
#     - Hosted runners are usually part of the organization's private network and can be easily misconfigured.
#     - If the hosted runner is insecurely configured, any GitHub user could:
#     - 1. Create a workflow that runs on the public hosted runner
#     - 2. Exploit the misconfigurations to execute code inside the private network
default runner_group_can_be_used_by_public_repositories := true

runner_group_can_be_used_by_public_repositories := false {
	not input.runner_group.allows_public_repositories
}

# METADATA
# scope: rule
# title: Runner Group Should Be Limited to Selected Repositories
# description: |
#       Not limiting the runner group to selected repositories allows any user in the organization to execute workflows on the group's runners.
#       In case of inadequate security measures implemented on the hosted runner,
#       malicious insider could create a repository with a workflow that exploits the runner's vulnerabilities to move laterally inside your network.
# custom:
#   severity: MEDIUM
#   requiredEnrichers: [organizationId]
#   requiredScopes: [admin:org]
#   remediationSteps:
#     - 1. Go to the organization settings page
#     - 2. Go to Actions ➝ Runner groups
#     - 3. Under the 'Repository Access' section, select 'Selected repositories'
#     - 4. Select the required repositories
#   threat:
#     - Hosted runners are usually part of the organization's private network and can be easily misconfigured.
#     - If the hosted runner is insecurely configured, any user in the organization could:
#     - 1. Create a workflow that runs on the hosted runner
#     - 2. Exploit the runner misconfigurations/known CVE's to execute code inside the private network
default runner_group_not_limited_to_selected_repositories := true

runner_group_not_limited_to_selected_repositories := false {
	input.runner_group.visibility == "selected"
}
