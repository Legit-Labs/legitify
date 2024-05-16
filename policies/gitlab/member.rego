package member

# METADATA
# scope: rule
# title: Two Factor Authentication Should Be Enabled for Collaborators
# description: A collaborator's two factor authentication is disabled. Turn it on in the collaborator setting, or globally in the account, to prevent any access without MFA.
# custom:
#   severity: HIGH
#   prerequisites: [enterprise]
#   remediationSteps:
#     - 1. Login with the user credentials
#     - 2. Go to the user settings page
#     - 3. Select 'Account' on the left navigation bar
#     - 4. Press 'Enable two-factor authentication'
#   threat:
#     - Collaborators without two-factor authentication are prime targets for phishing and social engineering attacks, as compromise only requires acquiring the collaborator's password.
default two_factor_authentication_is_disabled_for_a_collaborator := true

two_factor_authentication_is_disabled_for_a_collaborator := false {
	input.two_factor_enabled
}

# METADATA
# scope: rule
# title: Two Factor Authentication Should Be Enabled for External Collaborators
# description: An external collaborator's two factor authentication is disabled. Turn it on in the collaborator setting, or globally in the account, to prevent any access without MFA.
# custom:
#   severity: HIGH
#   prerequisites: [enterprise]
#   remediationSteps:
#     - 1. Login with the user credentials
#     - 2. Go to the user settings page
#     - 3. Select 'Account' on the left navigation bar
#     - 4. Press 'Enable two-factor authentication'
#   threat:
#     - Collaborators without two-factor authentication are prime targets for phishing and social engineering attacks, as compromise only requires acquiring the collaborator's password.
#     - This is doubly important for external collaborators, as these are identities that aren't likely managed by you or your organization and may be easier to compromise.
default two_factor_authentication_is_disabled_for_an_external_collaborator := true

two_factor_authentication_is_disabled_for_an_external_collaborator := false {
	input.external
	input.two_factor_enabled
}

# METADATA
# scope: rule
# title: Administrators Should Have Activity in the Last 6 Months
# description: A collaborator with global admin permissions didn't perform any action in the last 6 months. Admin users are extremely powerful, and common compliance standards demand keeping the number of admins to a minimum. Consider revoking this collaborator's admin credentials (downgrade to regular user) or removing the user completely.
# custom:
#   severity: MEDIUM
#   remediationSteps:
#     - 1. Go to the admin menu
#     - 2. Select 'Overview -> Users' on the left navigation bar
#     - 3. Find the stale admin and either delete or block it
#   threat:
#     - Stale admins are most likely not managed and monitored, increasing the possibility of being compromised.
default stale_admin_found := true

stale_admin_found := false {
	input.is_admin == true
	not is_null(input.last_sign_in_at)
	ns := time.parse_rfc3339_ns(input.last_sign_in_at)
	not isStale(ns, 6)
}

isStale(target_last_active, count_months) {
	now := time.now_ns()
	diff := time.diff(now, target_last_active)

	# diff[1] the months index
	diff[1] >= count_months
}
