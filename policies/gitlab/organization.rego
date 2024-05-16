package organization

# METADATA
# scope: rule
# title: Two-Factor Authentication Should Be Enforced For The Group
# description: The two-factor authentication requirement is not enabled at the group level. Regardless of whether users are managed externally by SSO, it is highly recommended to enable this option, to reduce the risk of a deliberate or accidental user creation without MFA.
# custom:
#   severity: HIGH
#   remediationSteps:
#     - 1. Go to the group page
#     - 2. Press Settings -> General
#     - 3. Expand 'Permissions and group features'
#     - 4. Toggle 'Require all users in this group to set up two-factor authentication'
#     - 5. Press 'Save Changes'
#   threat:
#     - If an attacker gets valid credentials for one of the organization’s users, they can authenticate to your GitHub organization.
default two_factor_authentication_not_required_for_group := true

two_factor_authentication_not_required_for_group := false {
	input.require_two_factor_authentication
}

# METADATA
# scope: rule
# title: Forking of Repositories to External Namespaces Should Be Disabled.
# description: The ability to fork a project to external namespaces is turned on. Forking a repository can lead to loss of control and potential exposure of source code. If you do not need forking, it is recommended to turn it off in the project's configuration. The option to fork should be enabled only by owners deliberately when opting to create a fork.
# custom:
#   severity: MEDIUM
#   remediationSteps:
#     - 1. Go to the top-level groups Settings > General page
#     - 2. Expand the Permissions and group features section
#     - 3. Check 'Prevent project forking outside current group'
#     - 4. Select 'Save changes'
#   threat:
#     - Forking to external namespaces could result in loss of control over proprietary information and potentially expose the organization to security risks, such as data leaks.
default collaborators_can_fork_repositories_to_external_namespaces := true

collaborators_can_fork_repositories_to_external_namespaces := false {
	input.prevent_forking_outside_group
}

# METADATA
# scope: rule
# title: Webhooks Should Be Configured To Use SSL
# description: Webhooks that are not configured with SSL enabled could expose your software to man-in-the-middle attacks (MITM).
# custom:
#   severity: LOW
#   requiredEnrichers: [hooksList]
#   remediationSteps:
#     - 1. Go to the group Settings -> Webhooks page
#     - 2. Find the misconfigured webhook and press 'Edit'
#     - 3. Toggle 'Enable SSL verification'
#     - 4. Press 'Save Changes'
#   threat:
#     - Webhooks with SSL verification disabled can be exploited by any party with access to the target DNS domain, allowing them to masquerade as your designated payload URL and freely read and affect the response of any webhook request.
#     - In the case of GitLab Self-Managed, it may be sufficient only to control the DNS configuration of the network where the instance is deployed.
organization_webhook_doesnt_require_ssl[violation] := true {
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
#     - 1. Go to the group page
#     - 2. Press Settings -> Repository
#     - 3. Expand 'Default Branch' section
#     - 4. Toggle the required protection rule
#     - 5. Press 'Save Changes'
#   threat:
#     - A developer creates a repository without any branch protection rules
#     - An attacker that gains access to the repository can modify its main branch without any restrictions
default group_does_not_enforce_branch_protection_by_default := true

group_does_not_enforce_branch_protection_by_default := false {
	input.default_branch_protection > 0
}

# METADATA
# scope: rule
# title: Two-Factor Authentication Grace Period Should Not Be Longer Than One Week
# description: New members added to your group are allowed longer than a week to enable MFA. The time frame should be lowered to one week or less.
# custom:
#   severity: MEDIUM
#   remediationSteps:
#     - 1. Go to the group page
#     - 2. Press Settings -> General
#     - 3. Expand 'Permissions and group features'
#     - "4. In the box titled: 'Delay 2FA enforcement (hours)', enter a number under 168 (preferably 0)"
#     - 5. Press 'Save Changes'
#   threat:
#     - Any new group member effectively acts as an attack surface until two-factor authentication is enabled. The risk is compounded as new members may be more vulnerable to phishing and identity theft attacks.
default group_allows_excessive_mfa_grace_period := true

group_allows_excessive_mfa_grace_period := false{
	input.two_factor_grace_period <= 168
}
