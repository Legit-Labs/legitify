package organization

# METADATA
# scope: rule
# title: Two-Factor Authentication Must Be Enforced For The Group
# description: The two-factor authentication requirement is not enabled at the group level. Regardless of whether users are managed externally by SSO, it is highly recommended to enable this option, to reduce the risk of a deliberate or accidental user creation without MFA.
# custom:
#   severity: HIGH
#   remediationSteps:
#     - Go to the group page
#     - Press Settings -> General
#     - Expand "Permissions and group features"
#     - Toggle "Require all users in this group to set up two-factor authentication"
#     - Press "Save Changes"
#   threat:
#     - If an attacker gets the valid credentials for one of the organization’s users they can authenticate to your GitHub organization.
default two_factor_authentication_not_required_for_group = false

two_factor_authentication_not_required_for_group {
	input.require_two_factor_authentication == false
}

# METADATA
# scope: rule
# title: Forking of Repositories to External Namespaces must be disabled.
# description: The ability to fork project to external namespaces is turned on. Forking a repository can lead to loss of control and potential exposure of source code. If you do not need forking, it is recommended to turn it off in the project's configuration. The option to fork should be enabled only by owners deliberately when opting to create a fork. 
# custom:
#   severity: MEDIUM
#   remediationSteps:
#     - "Go to the top-level groups Settings > General page"
#     - "Expand the Permissions and group features section"
#     - "Check Prevent project forking outside current group"
#     - "Select Save changes"
#   threat:
#     - Forking to external namespaces could result in loss of control over proprietary information and potentially expose the organization to security risks, such as data leaks.
default collaborators_can_fork_repositories_to_external_namespaces = false

collaborators_can_fork_repositories_to_external_namespaces {
	input.prevent_forking_outside_group == false
}

# METADATA
# scope: rule
# title: Webhooks Should Be Configured To Use SSL
# description: Webhooks that are not configured with SSL enabled could expose your sofware to man in the middle attacks (MITM).
# custom:
#   severity: LOW
#   requiredEnrichers: [hooksList]
#   remediationSteps:
#     - Go to the group Settings -> Webhooks page
#     - Find the misconfigured webhook and press "Edit"
#     - Toggle "Enable SSL verification"
#     - Press "Save Changes"
#   threat:
#     - "If SSL verification is disabled, any party with access to the target DNS domain can masquerade as your designated payload URL, allowing it freely read and affect the response of any webhook request."
#     - "In the case of GitLab Self-Managed, it may be sufficient only to control the DNS configuration of the network where the instance is deployed."
organization_webhook_doesnt_require_ssl[violation] = true {
	some index
	hook := input.hooks[index]
	hook.enable_ssl_verification == false
	violation := {"id": format_int(hook.id, 10), "url": hook.url}
}

# METADATA
# scope: rule
# title: Group Should Enforce Branch Protection
# description: The default branch should be protected in each group so that any new repository will be created with a protected default branch by default. In fully protected level, developers cannot push new commits, and no one can force push or delete the branch. Protecting branches ensures new code changes must go through a controlled merge process and it allows enforcement of code review and other security tests.
# custom:
#   severity: MEDIUM
#   remediationSteps:
#     - Go to the group page
#     - Press Settings -> Repository
#     - Expand "Default Branch" section
#     - Toggle the required protection rule
#     - Press "Save Changes"
#   threat:
#     - A developer creates a repository without any branch protection rules
#     - Attacker that get access to the repository can modify its main branch without any restrictions
default group_does_not_enforce_branch_protection_by_default = false

group_does_not_enforce_branch_protection_by_default {
	input.default_branch_protection == 0
}
