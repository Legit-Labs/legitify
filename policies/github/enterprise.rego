package enterprise

# METADATA
# scope: rule
# title: Enterprise Should Prevent Repository Admins From Changing Repository Visibility
# description: The enterprise's Repository visibility change policy should be set to DISABLED. This will prevent users from creating private repositories and changing them to be public. Malicious actors could leak code if enabled.
# custom:
#   severity: MEDIUM
#   remediationSteps:
#     - 1. Make sure you are an enterprise owner
#     - 2. Go to the policies page
#     - 3. Under the "Repository visibility change" section, choose the "Disabled" option
#   requiredScopes: [admin:enterprise]
#   threat:
#     - 'A member of the organization could inadvertently or maliciously make public an internal repository exposing confidential data.'
default enterprise_not_using_visibility_change_disable_policy := true

enterprise_not_using_visibility_change_disable_policy := false {
	input.members_can_change_repository_visibility == "DISABLED"
}


# METADATA
# scope: rule
# title: Enterprise Should Prevent Members From Forking Internal And Private Repositories
# description: The enterprise's repository forking policy should be set to DISABLED. Forking a repository can lead to loss of control and potential exposure of source code. If you do not need forking, it is recommended to turn it off in the project's configuration. The option to fork should be enabled only by owners deliberately when opting to create a fork.
# custom:
#   severity: LOW
#   remediationSteps:
#     - 1. Make sure you are an enterprise owner
#     - 2. Go to the policies page
#     - 3. Under the "Repository Forking" section, Choose the "Disabled" option
#   requiredScopes: [admin:enterprise]
#   threat:
#     - Forking to external namespaces could result in loss of control over proprietary information and potentially expose the organization to security risks, such as data leaks.
default enterprise_allows_forking_repos := true

enterprise_allows_forking_repos := false {
	input.repositories_forking_policy == "DISABLED"
}


# METADATA
# scope: rule
# title: Enterprise Should Prevent Members From Creating Public Repositories
# description: The enterprise's repository creation policy should be set to private/internal repositories only. This will prevent non-admin users from creating public repositories and potentially exposing source code.
# custom:
#   severity: MEDIUM
#   remediationSteps:
#     - 1. Make sure you are an enterprise owner
#     - 2. Go to the policies page
#     - 3. Under the "Repository creation" section, Choose the "Members can create repositories" option and uncheck 'Public'
#   requiredScopes: [admin:enterprise]
#   threat:
#     - Users can accidentally create public repositories and expose source code.
default enterprise_allows_creating_public_repos := true

enterprise_allows_creating_public_repos := false {
	input.members_can_create_public_repositories == false
}


# METADATA
# scope: rule
# title: Enterprise Should Prevent Members From Inviting Outside Collaborators
# description: The enterprise's external collaborators invite policy should be set to enterprise/organization owners only. Allowing members to invite external collaborators might result in unauthorized access to internal projects.
# custom:
#   severity: MEDIUM
#   remediationSteps:
#     - 1. Make sure you are an enterprise owner
#     - 2. Go to the policies page
#     - 3. Under the "Repository outside collaborators" section, choose the "Enterprise Owners Only" or the "Organization Owners Only" option
#   requiredScopes: [admin:enterprise]
#   threat:
#     - Inviting external collaborators could result in a loss of control over proprietary information and potentially expose the organization to security risks, such as data leaks.
default enterprise_allows_inviting_externals_collaborators := true

enterprise_allows_inviting_externals_collaborators := false {
	input.external_collaborators_invite_policy == "DISABLED"
}

# METADATA
# scope: rule
# title: Two-Factor Authentication Should Be Enforced For The Enterprise
# description: The two-factor authentication requirement should be enforced at the enterprise level. Regardless of whether users are managed externally by SSO, it is highly recommended to enable this option to reduce the risk of deliberate or accidental user creation without MFA.
# custom:
#   severity: HIGH
#   remediationSteps:
#     - 1. Make sure you are an enterprise owner
#     - 2. Go to the Settings page
#     - 3. Go to the Authentication security tab
#     - 4. Check the "Require two-factor authentication for all organizations in the enterprise" checkbox
#   requiredScopes: [admin:enterprise]
#   threat:
#     - If an attacker gets valid credentials for one of the enterprise’s users they can authenticate to your GitHub enterprise.
default enterprise_enforce_two_factor_authentication := true

enterprise_enforce_two_factor_authentication := false {
	input.two_factor_required_setting == "ENABLED"
}


# METADATA
# scope: rule
# title: Enterprise Should Use Single-Sign-On
# description: It is recommended to enable access to an enterprise via SAML single sign-on (SSO) by authenticating through an identity provider (IdP). This allows for central account control and timely access revocations.
# custom:
#   severity: MEDIUM
#   remediationSteps:
#     - 1. Make sure you are an enterprise owner
#     - 2. Go to the Settings page
#     - 3. Go to the Authentication security tab
#     - 4. Toggle on "Enable SAML authentication"
#     - 5. Fill in the remaining SSO configuration as instructed on the screen
#     - 6. Click "Save"
#   requiredScopes: [admin:enterprise]
#   threat: Not using an SSO solution makes it more difficult to track a potentially compromised user's actions across different systems, prevents common password policy throughout the enterprise, and makes it challenging to audit different aspects of the user's behavior.
default enterprise_not_using_single_sign_on := true

enterprise_not_using_single_sign_on := false {
	input.saml_enabled
}

# METADATA
# scope: rule
# title: Enterprise Should Define Base Permissions As "No Permission" For All Members
# description: Collaborators in your organizations should receive access to specific organizations and repositories as necessary, and not have read and write access to all repositories across the enterprise.
# custom:
#   severity: MEDIUM
#   remediationSteps:
#     - 1. Make sure you are an enterprise owner
#     - 2. Go to the Settings page
#     - 3. Under the 'Policies' tab, choose 'Repositories'
#     - 4. Under 'Base Permission' choose 'No Permission'
#   requiredScopes: [admin:enterprise]
#   threat: An adversary will have access to all repositories in the enterprise, instead of just a part of them.
default repository_no_permission_enforced_by_default := true

repository_no_permission_enforced_by_default := false {
	input.default_repository_no_permission_enforced == "NONE"
}

# METADATA
# scope: rule
# title: Enterprise Should Prevent Repository Admins From Deleting Or Transferring Repositories
# description: The enterprise’s Repository deletion and transfer policy should be set to DISABLED. This will prevent repository admins from deleting a repo or transferring it to a different owner or organization. Malicious actors could leak code if enabled.
# custom:
#   severity: MEDIUM
#   remediationSteps:
#     - 1. Make sure you are an enterprise owner
#     - 2. Go to the Enterprise Settings page
#     - 3. Under the ‘Policies’ tab choose ‘Repositories’
#     - 4. Go to the ‘Admin repository permissions' section, under ‘Repository deletion and transfer' and select 'Disabled'
#   requiredScopes: [admin:enterprise]
#   threat: A member of the organization could inadvertently or maliciously transfer a repository to an external namespace and expose confidential data.
default memberes_allowed_repository_move_or_deletion := true

memberes_allowed_repository_move_or_deletion := false {
	input.member_can_delete_repository == "DISABLED"
}

# METADATA
# scope: rule
# title: Enterprise Should Automatically Enable Advanced Security Across All Organizations/Repositories
# description: Advanced Security includes code scanning, secret scanning and dependency review. These features protect your repositories from containing vulnerable data. Prevents the risk of unauthorized access or exploitation of vulnerabilities.
# custom:
#   severity: MEDIUM
#   remediationSteps:
#     - 1. Make sure you are an enterprise owner
#     - 2. Go to the Enterprise Settings page
#     - 3. Under the ‘Settings’ tab choose ‘Code security and analysis’
#     - 4. Check 'Automatically enable for new repositories'
#   requiredScopes: [admin:enterprise]
default enable_ghas_for_new_orgs := true

enable_ghas_for_new_orgs := false {
	input.code_analysis_and_security_policies.advanced_security_enabled_for_new_repositories == true
}

# METADATA
# scope: rule
# title: Enterprise Should Automatically Enable Secret Scanning Across All Organizations/Repositories
# description: Enable GitHub Advanced Security secret scanning to alert on sensitive data that exists in your enterprise. Secrets shouldn’t be hard-coded in to your repositories as they will be retrievable by anyone with access to the repository.
# custom:
#   severity: MEDIUM
#   remediationSteps:
#     - 1. Make sure you are an enterprise owner
#     - 2. Go to the Enterprise Settings page
#     - 3. Under the ‘Settings’ tab choose ‘Code security and analysis’
#     - 4. Check 'Automatically enable for new repositories with Advanced Security enabled'
#   requiredScopes: [admin:enterprise]
default enable_secret_scanning_for_new_orgs := true

enable_secret_scanning_for_new_orgs := false {
	input.code_analysis_and_security_policies.secret_scanning_enabled_for_new_repositories == true
}

# METADATA
# scope: rule
# title: Enterprise Should Automatically Enable Secret Scanning Push Protection Across All Organizations/Repositories
# description: The enterprise should prevent sensitive data from being pushed to all repositories, to prevent it from being exposed to anyone with access to the repository.
# custom:
#   severity: MEDIUM
#   remediationSteps:
#     - 1. Make sure you are an enterprise owner
#     - 2. Go to the Enterprise Settings page
#     - 3. Under the ‘Settings’ tab choose ‘Code security and analysis’
#     - 4. Check 'Automatically enable for repositories added to secret scanning'
#   requiredScopes: [admin:enterprise]
default enable_push_protection_secret_scanning := true

enable_push_protection_secret_scanning := false {
	input.code_analysis_and_security_policies.secret_scanning_push_protection_enabled_for_new_repositories == true
}

# METADATA
# scope: rule
# title: Enterprise Should Send Email Notifications Only To Verified Domains
# description: The enterprise should mitigate the leakage of sensitive data by allowing email notifications to be sent only to verified or approved domains.
# custom:
#   severity: MEDIUM
#   remediationSteps:
#     - 1. Make sure you are an enterprise owner
#     - 2. Go to the Enterprise Landing page
#     - 3. Under the ‘Settings’ tab on the left click ‘Verified & approved domains’
#     - 4. Press the 'Add a domain' button and follow the instructions in the menu
#     - 5. Check the 'Restrict email notifications to only approved or verified domains' box
#     - 6. Press 'Save'
#   requiredScopes: [admin:enterprise]
default enable_email_notification_to_verified_domains := true

enable_email_notification_to_verified_domains := false {
	input.notification_delivery_restriction_enabled == "ENABLED"
}
