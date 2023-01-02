package member

# METADATA
# scope: rule
# title: Two Factor Authentication Is Disabled For a Collaborator
# description: A collaborator's two factor authentication is disabled at the source code management system. Turn it on in the collaborator setting, or globally in the account, to prevent any access without MFA.
# custom:
#   severity: HIGH
#   prerequisites: [enterprise]
#   remediationSteps:
#     - Login with the user credentials
#     - Go to the user settings page
#     - Select "Account" on the left navigation bar
#     - Press "Enable two-factor authentication"
#   threat: Collaborators without two-factor authentication are prime targets for phising and social engineering attacks, as compromise only requires acquiring the collaborator's password.
default two_factor_authentication_is_disabled_for_a_collaborator = false
two_factor_authentication_is_disabled_for_a_collaborator {
    input.two_factor_enabled == false
}

# METADATA
# scope: rule
# title: Two Factor Authentication Is Disabled For an External Collaborator
# description: An external collaborator's two factor authentication is disabled at the source code management system. Turn it on in the collaborator setting, or globally in the account, to prevent any access without MFA.
# custom:
#   severity: HIGH
#   prerequisites: [enterprise]
#   remediationSteps:
#     - Login with the user credentials
#     - Go to the user settings page
#     - Select "Account" on the left navigation bar
#     - Press "Enable two-factor authentication"
#   threat: 
#     - "Collaborators without two-factor authentication are prime targets for phising and social engineering attacks, as compromise only requires acquiring the collaborator's password."
#     - "This is doubly important for external collaborators, as these are identities that aren't likely managed by you or your organization and may be easier to compromise."
default two_factor_authentication_is_disabled_for_an_external_collaborator = false
two_factor_authentication_is_disabled_for_an_external_collaborator {
    input.external == true
    input.two_factor_enabled == false
}

# METADATA
# scope: rule
# title: Stale admin detected
# description: A collaborator with global admin permissions didn't do any action in the last 6 months. Admin users are extremely powerful and common compliance standards demand keeping the number of admins at minimum. Consider revoking this collaborator admin credentials (downgrade to regular user), or remove the user completely.
# custom:
#   severity: MEDIUM
#   remediationSteps:
#     - Go to admin menu
#     - Select "Overview -> Users" on the left navigation bar
#     - Find the stale admin and either delete of block it
#   threat:
#     - "Stale admins are most likely not managed and monitored, increasing the possibility of being compromised."
default stale_admin_found = false
stale_admin_found {
    input.is_admin == true
    not is_null(input.last_sign_in_at)
    ns := time.parse_rfc3339_ns(input.last_sign_in_at)
    isStale(ns, 6)
}

isStale(target_last_active, count_months) {
    now := time.now_ns()
    diff := time.diff(now, target_last_active)

    # diff[1] the months index
    diff[1] >= count_months
}

is_null(x) {
    x == null
}
