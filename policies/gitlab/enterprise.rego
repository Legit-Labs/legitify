package enterprise

# METADATA
# scope: rule
# title: Two-factor Authentication Should Be Globally Enforced
# description: >
#    It is recommended to turn on MFA at the server or account level, and proactively prevent any new user created without MFA.
#    Even if identities are managed externally using an SSO,
#    it is highly recommended to keep this option on, for the 'admin' user or to be protected in the future for a deliberate or incidental creation of a user without MFA.
# custom:
#   severity: HIGH
#   remediationSteps:
#     - 1. Go to the admin page: Menu -> Admin
#     - 2. Press Settings -> General
#     - 3. Expand 'Sign-in restrictions' section
#     - 4. Toggle 'Two-factor authentication'
#     - 5. Press 'Save Changes'
default require_two_factor_authentication_not_globally_enforced := true

require_two_factor_authentication_not_globally_enforced := false {
	input.require_two_factor_authentication
}

# METADATA
# scope: rule
# title: Password Authentication For Git Over HTTP(S) Should Not Be Enabled
# description: >
#     Password authentication for Git protocol operations (pull / push) is discouraged in favor of using personal access tokens and keys.
#     From a security standpoint, it is recommended to disable password authentication completely.
# custom:
#   severity: LOW
#   remediationSteps:
#     - 1. Go to the admin page: Menu -> Admin
#     - 2. Press Settings -> General
#     - 3. Expand 'Sign-in restrictions' section
#     - 4. Un toggle 'Allow password authentication for Git over HTTP(S)'
#     - 5. Press 'Save Changes'
default password_authentication_enabled_for_git_is_enabled := true

password_authentication_enabled_for_git_is_enabled := false {
	not input.password_authentication_enabled_for_git
}

# METADATA
# scope: rule
# title: API Request Rate Limit Should Be Limited
# description: >
#     Enabling authenticated request rate-limit on APIs reduces volume of requests, and prevents incidents of reduced quality-of-service
#     or denial-of-service because of unnecessary or malicious high request volume.
# custom:
#   severity: LOW
#   remediationSteps:
#     - 1. Go to the admin page: Menu -> Admin
#     - 2. Press Settings -> Network
#     - 3. Expand 'User and IP rate limit' section
#     - 4. Toggle 'Enable authenticated API request rate limit'
#     - 5. Press 'Save Changes'
default throttle_authenticated_api_not_enabled := true

throttle_authenticated_api_not_enabled := false {
	input.throttle_authenticated_api_enabled
}

# METADATA
# scope: rule
# title: Branch Protection Should Be Globally Enabled By Default
# description: >
#     Branch protection is not enabled for the entire server or account by default. Any newly created repository will not be protected unless explicitly specified.
#     Protecting branches ensures new code changes must go through a controlled merge process and it allows enforcement of code review and other security tests. It is recommended to turn it on by default.
# custom:
#   severity: MEDIUM
#   remediationSteps:
#     - 1. Go to the admin page: Menu -> Admin
#     - 2. Press Settings -> Repository
#     - 3. Expand 'Default Branch' section
#     - 4. Toggle 'Fully protected'
#     - 5. Press 'Save Changes'
default default_branch_protection_not_globally_enforced := true

default_branch_protection_not_globally_enforced := false {
	input.default_branch_protection > 0
}

# METADATA
# scope: rule
# title: Creating Public Repositories Should Be Restricted To Admins
# description: >
#     Your server allows non-admin members to create public repositories.
#     Creating a public repository can be done by mistake, and may expose sensitive organization code, that once exposed may be copied, cached or stored by external parties.
#     Therefore, it is highly recommended to restrict the option to create public repositories to admins only and reduce the risk of unintentional code exposure.
# custom:
#   severity: HIGH
#   remediationSteps:
#     - 1. Go to the admin page: Menu -> Admin
#     - 2. Press Settings -> General
#     - 3. Expand 'Visibility and access controls' section
#     - 4. Under 'Restricted visibility levels' toggle 'Public'
#     - 5. Press 'Save Changes'
default all_users_are_allowed_to_create_public_projects := true

all_users_are_allowed_to_create_public_projects := false {
	public := [rule | rule := input.restricted_visibility_levels[_]; rule == "public"]
	count(public) > 0
}

# METADATA
# scope: rule
# title: Default Group Visibility Should Not Be Public
# description: >
#     Your server is configured by default that every new group is created public. This means that it will be accessible by any member and unauthenticated user. It is recommended to turn the default option to be private
#     and let admins control public visibility explicitly when needed.
# custom:
#   severity: HIGH
#   remediationSteps:
#     - 1. Go to the admin page: Menu -> Admin
#     - 2. Press Settings -> General
#     - 3. Expand 'Visibility and access controls' section
#     - 4. Under 'Default group visibility' toggle 'Private'
#     - 5. Press 'Save Changes'
default default_group_visibility_is_public := true

default_group_visibility_is_public := false {
	not input.default_group_visibility == "public"
}

# METADATA
# scope: rule
# title: Default Project Visibility Should Not Be Public
# description: >
#     Your server is configured by default that every new project is created public. This means that it will be accessible by any member and unauthenticated user. It is recommended to turn the default option to be private
#      and let admins control public visibility explicitly when needed.
# custom:
#   severity: HIGH
#   remediationSteps:
#     - 1. Go to the admin page: Menu -> Admin
#     - 2. Press Settings -> General
#     - 3. Expand 'Visibility and access controls' section
#     - 4. Under 'Default project visibility' toggle 'Private'
#     - 5. Press 'Save Changes'
default default_project_visibility_is_public := true

default_project_visibility_is_public := false {
	not input.default_project_visibility == "public"
}

# METADATA
# scope: rule
# title: Sign-Up Confirmation Email Should Be Mandatory
# description: >
#     This security option enforces any new user to confirm their email address, and ensures the user is verified. This increases the security and integrity of the identities created within the system, and prevents situations
#     where a user is spoofed by a malicious actor while using a legitimate corporate email address.
# custom:
#   severity: MEDIUM
#   remediationSteps:
#     - 1. Go to the admin page: Menu -> Admin
#     - 2. Press Settings -> General
#     - 3. Expand 'Sign-up restrictions' section
#     - 4. Toggle 'Send confirmation email on sign-up'
#     - 5. Press 'Save Changes'
default send_user_confirmation_email_not_enabled := true

send_user_confirmation_email_not_enabled := false {
	input.send_user_confirmation_email
}

# METADATA
# scope: rule
# title: Unauthenticated Requests Rate Limit Should Be Enabled
# description: >
#     The server allows restricting the limit of unauthenticated requests. It is recommended to turn it on as a security and reliability measure,
#     and to reduce request volume. If an attacker tries accessing the system, this will reduce the risk of brute-force and Denial-of-service to the end users caused by high request rate.
# custom:
#   severity: MEDIUM
#   remediationSteps:
#     - 1. Go to the admin page: Menu -> Admin
#     - 2. Press Settings -> Network
#     - 3. Expand 'User and IP rate limit' section
#     - 4. Toggle 'Enable unauthenticated API request rate limit' and 'Enable unauthenticated web request rate limit'
#     - 5. Press 'Save Changes'
default throttle_unauthenticated_request_not_enabled := true

throttle_unauthenticated_request_not_enabled := false {
    input.throttle_unauthenticated_api_enabled
    input.throttle_unauthenticated_web_enabled
}

# METADATA
# scope: rule
# title: Webhooks Should Not Be Allowed To Be Sent To The Local Network
# description: >
#     Webhooks sent by GitLab servers are authenticated, and can cause potential damage if sent uncontrollably. For example, a malicious user can plant a webhook address
#     that triggers an API call in GitLab itself which can delete resources. Therefore, as a security best practice, webhooks should be limited to external URLs only.
#     You can read more [here](https://docs.gitlab.com/ee/security/webhooks.html)
# custom:
#   severity: LOW
#   remediationSteps:
#     - 1. Go to the admin page: Menu -> Admin
#     - 2. Press Settings -> Network
#     - 3. Expand 'Outbound requests' section
#     - 4. Un toggle 'Allow requests to the local network from web hooks and services'
#     - 5. Press 'Save Changes'
default webhooks_are_allowed_to_be_sent_to_local_network := true

webhooks_are_allowed_to_be_sent_to_local_network := false {
    not input.allow_local_requests_from_web_hooks_and_services
}

# METADATA
# scope: rule
# title: Server Should Not Allow Access To Unauthenticated Users With Sign-Up
# description: The server allows any person with network access to sign up, create a user and access sensitive data. Turning this off will reduce the risk of attackers trying to infiltrate the server.
# custom:
#   severity: HIGH
#   remediationSteps:
#     - 1. Go to the admin page: Menu -> Admin
#     - 2. Press Settings -> General
#     - 3. Expand 'Sign-up restrictions' section
#     - 4. Un toggle 'Sign-up enabled'
#     - 5. Press 'Save Changes'
default unauthenticated_signup_enabled := true

unauthenticated_signup_enabled := false {
    not input.signup_enabled
}
