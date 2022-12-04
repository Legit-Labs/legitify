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
