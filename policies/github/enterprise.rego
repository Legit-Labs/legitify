package enterprise

# METADATA
# scope: rule
# custom:
#   severity: MEDIUM
# title: Enterprise Should Not Allow Members To Change Repository Visibility
# description: The enterprise's Repository visibility change policy should be set to DISABLED. This will prevents users from creating private repositories and change them to be public. Malicous actors could leak code if enabled.
# custom:
#   remediationSteps: [Make sure you are an enterprise owner, Go to the policies page, Under the "Repository visibility change" section, choose the "Disabled" option]
#   requiredScopes: [admin:enterprise]
#   threat:
#     - "A member of the organization could inadvertently or maliciously make public an internal repository exposing confidential data."
default enterprise_not_using_visibility_change_disable_policy = true

enterprise_not_using_visibility_change_disable_policy = false {
	input.members_can_change_repository_visibility == "DISABLED"
}


# METADATA
# scope: rule
# custom:
#   severity: MEDIUM
# title: Enterprise Should Not Allow Members To Fork Internal And Private Repositories
# description: The enterprise's repository forking policy should be set to DISABLED. Forking a repository can lead to loss of control and potential exposure of source code. If you do not need forking, it is recommended to turn it off in the project's configuration. The option to fork should be enabled only by owners deliberately when opting to create a fork.
# custom:
#   remediationSteps: [Make sure you are an enterprise owner, Go to the policies page, Under the "Repository Forking" section, Choose the "Disabled" option]
#   requiredScopes: [admin:enterprise]
#   threat:
#     - Forking to external namespaces could result in loss of control over proprietary information and potentially expose the organization to security risks, such as data leaks.
default enterprise_allows_forking_repos = true

enterprise_allows_forking_repos = false {
	input.repositories_forking_policy == "DISABLED"
}


# METADATA
# scope: rule
# custom:
#   severity: MEDIUM
# title: Enterprise Should Not Allow Members To Invite Outside Collaborators
# description: The enterprise's external collaborators invite policy should be set to enterprise/organization owners only. Allowing members to invite external collaborators might result in unauthorized access to the internal projects.
# custom:
#   remediationSteps: [Make sure you are an enterprise owner, Go to the policies page, Under the "Repository outside collaborators" section - choose the "Enterprise Owners Only" or the "Organization Owners Only" option]
#   requiredScopes: [admin:enterprise]
#   threat:
#     - Inviting external collaborators could result in a loss of control over proprietary information and potentially expose the organization to security risks, such as data leaks.
default enterprise_allows_inviting_externals_collaborators = true

enterprise_allows_inviting_externals_collaborators = false {
	input.repositories_forking_policy == "DISABLED"
}

# METADATA
# scope: rule
# custom:
#   severity: HIGH
# title: Two-Factor Authentication Should Be Enforced For The Enterprise
# description: The two-factor authentication requirement should be enforced at the enterprise level. Regardless of whether users are managed externally by SSO, it is highly recommended to enable this option to reduce the risk of a deliberate or accidental user creation without MFA.
# custom:
#   remediationSteps: [Make sure you are an enterprise owner, Go to the Settings page, Go to the Authentication security tab, Check the "Require two-factor authentication for all organizations in the enterprise" checkbox]
#   requiredScopes: [admin:enterprise]
#   threat:
#     - If an attacker gets the valid credentials for one of the enterprise’s users they can authenticate to your GitHub enterprise.
default enterprise_enforce_two_factor_authentication = true

enterprise_enforce_two_factor_authentication = false {
	input.two_factor_required_setting == "ENABLED"
}



# METADATA
# scope: rule
# custom:
#   severity: MEDIUM
# title: Enterprise Should Use Single-Sign-On
# description: It is recommended to enable access to an enterprise via SAML single sign-on (SSO) by authenticating through an identity provider (IdP). This allows for central account control and for timely access revocations.
# custom:
#   remediationSteps: [Make sure you are an enterprise owner, Go to the Settings page, than, go to the Authentication security tab. Toggle on "Enable SAML authentication", Fill in the remaining SSO configuration as instructed on the screen, Click "Save"]
#   requiredScopes: [admin:enterprise]
#   threat: Not using an SSO solution makes it more difficult to track a potentially compromised user's actions accross different systems, prevents the organization from enterprise a common password policy, and makes it challenging to audit different aspects of the user's behavior.
default enterprise_not_using_single_sign_on = true

enterprise_not_using_single_sign_on = false {
	input.saml_enabled
}
