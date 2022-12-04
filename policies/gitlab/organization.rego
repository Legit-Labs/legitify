package organization

# METADATA
# scope: rule
# title: Two-Factor Authentication Is Not Enforced For The Group
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
default two_factor_authentication_not_required_for_group  = false
two_factor_authentication_not_required_for_group {
    input.require_two_factor_authentication == false
}

# METADATA
# scope: rule
# title: Collaborators Can Fork Repositories To External Namespaces
# description: The ability to fork project to external namespaces is turned on. Forking repositories poses security issues due to the loss of control over the code. It is recommended to disable this feature if it is not explicitly needed, in order to proactively prevent code leakage.
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
# title: Webhook Configured Without SSL
# description: Webhooks that are not configured with SSL enabled could expose your software to man in the middle attacks (MITM).
# custom:
#   severity: LOW
#   requiredEnrichers: [genericList]
#   remediationSteps:
#     - Go to the group Settings -> Webhooks page
#     - Find the misconfigured webhook and press "Edit"
#     - Toggle "Enable SSL verification"
#     - Press "Save Changes"
organization_webhook_doesnt_require_ssl[violation] = true {
    some index
    hook := input.hooks[index]
    hook.enable_ssl_verification == false
    violation := { "name": format_int(hook.id, 10), "url": hook.url }
}
