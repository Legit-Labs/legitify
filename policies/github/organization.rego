package organization

import data.common.webhooks as webhookUtils
import data.common.secrets as secretUtils

# METADATA
# scope: rule
# title: Webhooks Should Be Configured With A Secret
# description: Webhooks are not configured with a shared secret to validate the origin and content of the request. This could allow your webhook to be triggered by any bad actor with the URL.
# custom:
#   requiredEnrichers: [hooksList]
#   severity: LOW
#   remediationSteps:
#     - 1. Make sure you have admin permissions
#     - 2. Go to the organization settings page
#     - 3. Select 'Webhooks'
#     - 4. Press on the insecure webhook
#     - 5. Configure a secret
#     - 6. Click 'Update webhook'
#   requiredScopes: [admin:org_hook]
#   threat:
#     - Not using a webhook secret makes the service receiving the webhook unable to determine the authenticity of the request.
#     - This allows attackers to masquerade as your organization, potentially creating an unstable or insecure state in other systems.
organization_webhook_no_secret[violated] := true {
	some index
	hook := input.hooks[index]
	not webhookUtils.has_secret(hook)
	violated := {
		"name": hook.name,
		"url": hook.url,
	}
}

# METADATA
# scope: rule
# title: Webhooks Should Be Configured To Use SSL
# description: Webhooks that are not configured with SSL enabled could expose your software to man in the middle attacks (MITM).
# custom:
#   requiredEnrichers: [hooksList]
#   severity: LOW
#   remediationSteps:
#     - 1. Make sure you have admin permissions
#     - 2. Go to the organization settings page
#     - 3. Select 'Webhooks'
#     - 4. Press on the insecure webhook
#     - 5. Verify URL starts with https
#     - 6. Enable 'SSL verification'
#     - 7. Click 'Update webhook'
#   requiredScopes: [admin:org_hook]
#   threat:
#     - If SSL verification is disabled, any party with access to the target DNS domain can masquerade as your designated payload URL, allowing it freely read and affect the response of any webhook request.
#     - In the case of GitHub Enterprise Server instances, it may be sufficient only to control the DNS configuration of the network where the instance is deployed, as an attacker can redirect traffic to the target domain in your internal network directly to them, and this is often much easier than compromising an internet-facing domain.
organization_webhook_doesnt_require_ssl[violated] := true {
	some index
	hook := input.hooks[index]
	not webhookUtils.ssl_enabled(hook)
	violated := {
		"name": hook.name,
		"url": hook.url,
	}
}

# METADATA
# scope: rule
# title: Two-Factor Authentication Should Be Enforced For The Organization
# description: The two-factor authentication requirement is not enabled at the organization level. Regardless of whether users are managed externally by SSO, it is highly recommended to enable this option to reduce the risk of a deliberate or accidental user creation without MFA.
# custom:
#   severity: HIGH
#   remediationSteps:
#     - 1. Make sure you have admin permissions
#     - 2. Go to the organization settings page
#     - 3. Enter 'Authentication security' tab
#     - 4. Under 'Two-factor authentication'
#     - 5. Toggle on 'Require two-factor authentication for everyone in the <ORG> organization'
#     - 6. Click 'Save'
#   requiredScopes: [admin:org]
#   threat:
#     - If an attacker gets the valid credentials for one of the organization’s users they can authenticate to your GitHub organization.
default two_factor_authentication_not_required_for_org := true

two_factor_authentication_not_required_for_org := false {
	input.organization.two_factor_requirement_enabled
}

# METADATA
# scope: rule
# title: Only Admins Should Be Able To Create Public Repositories
# description: >-
#   The organization should be configured to prevent non-admin members from creating public repositories. Creating a public repository may expose sensitive organization code, which, once exposed, may be copied, cached, or stored by external parties. Therefore, it is highly recommended to restrict the option to create public repositories to admins only and reduce the risk of unintentional code exposure. NOTE: You should also verify that repository owners can't change existing repository visibility to be public. If allowed, a malicious user could create a private repo and change it to public. See: https://docs.github.com/en/enterprise-cloud@latest/organizations/managing-organization-settings/restricting-repository-visibility-changes-in-your-organization for further information
# custom:
#   severity: MEDIUM
#   remediationSteps:
#     - 1. Make sure you have admin permissions
#     - 2. Go to the organization settings page
#     - 3. Enter 'Member privileges' tab
#     - 4. Under 'Repository creation'
#     - 5. Toggle off 'Public'
#     - 6. Click 'Save'
#   requiredScopes: [read:org]
#   threat:
#     - A member of the organization could inadvertently or maliciously make public an internal repository exposing confidential data.
#   compliance:
#     - soc2: [CC4.1, CC4.2, CC5.2, CC5.3, CC6.1, CC6.2, CC6.3, CC6.6, CC6.7, CC7.1, CC7.2, CC7.3, CC8.1, CC9.2]
#     - iso27002: [5.2, 5.3, 5.12, 5.15, 5.18, 5.21, 5.24, 5.25, 5.26, 5.28, 5.33, 5.36, 5.37, 8.2, 8.3, 8.4, 8.6, 8.8, 8.9, 8.10, 8.12, 8.15, 8.16, 8.25, 8.28, 8.29, 8.32]
default non_admins_can_create_public_repositories := true

non_admins_can_create_public_repositories := false {
	not input.organization.members_can_create_public_repositories
}

# METADATA
# scope: rule
# title: Default Member Permissions Should Be Restricted
# description: Default repository permissions configuration is not set in the organization, thus every new repository will be accessible by default to all users. It is strongly recommended to remove the default permissions and assign them on demand.
# custom:
#   severity: HIGH
#   remediationSteps:
#     - 1. Make sure you have admin permissions
#     - 2. Go to the organization settings page
#     - 3. Enter 'Member privileges' tab
#     - 4. Under 'Base permissions'
#     - 5. Set permissions to 'No permissions'
#     - 6. Click 'Save'
#   requiredScopes: [read:enterprise]
#   threat:
#     - Organization members can see the content of freshly created repositories, even if they should be restricted.
default default_repository_permission_is_not_none := true

default_repository_permission_is_not_none := false {
	input.organization.default_repository_permission == "none"
}

# METADATA
# scope: rule
# title: Organization Should Use Single-Sign-On
# description: It is recommended to enable access to an organization via SAML single sign-on (SSO) by authenticating through an identity provider (IdP). This allows for central account control and for timely access revocations.
# custom:
#   severity: MEDIUM
#   remediationSteps:
#     - 1. Make sure you have admin permissions
#     - 2. Go to the organization settings page
#     - 3. Enter 'Authentication security' tab
#     - 4. Toggle on 'Enable SAML authentication'
#     - 5. Fill in the remaining SSO configuration as instructed on the screen
#     - 6. Click 'Save'
#   requiredScopes: [admin:org]
#   threat:
#     - Not using an SSO solution makes it more difficult to track a potentially compromised user's actions across different systems, prevents the organization from defining a common password policy, and makes it challenging to audit different aspects of the user's behavior.
default organization_not_using_single_sign_on := true

organization_not_using_single_sign_on := false {
	input.saml_enabled
}

# METADATA
# scope: rule
# title: Organization Secrets Should Be Updated At Least Yearly
# description: Some of the organizations secrets have not been updated for over a year. It is recommended to refresh secret values regularly in order to minimize the risk of breach in case of an information leak.
# custom:
#   requiredEnrichers: [secretsList]
#   remediationSteps:
#      - 1. Enter your organization's landing page
#      - 2. Go to the settings tab
#      - 3. Under the 'Security' title on the left, choose 'Secrets and variables'
#      - 4. Click 'Actions'
#      - 5. Sort secrets by 'Last Updated'
#      - 6. Regenerate every secret older than one year and add the new value to GitHub's secret manager
#   severity: MEDIUM
#   requiredScopes: [admin:org, repo]
#   threat: Sensitive data may have been inadvertently made public in the past, and an attacker who holds this data may gain access to your current CI and services. In addition, there may be old or unnecessary tokens that have not been inspected and can be used to access sensitive information.
organization_secret_is_stale[stale] := true{
    some index
    secret := input.organization_secrets[index]
    secretUtils.is_stale(secret.updated_at)
    stale := {
    "name" : secret.name,
    "update date" : time.format(secret.updated_at),
    }
}

# METADATA
# scope: rule
# title: Web-Based Commit Signoff Should Be Required
# description: The organization does not require contributors to sign off on web-based commits. Enabling this setting enforces Developer Certificate of Origin (DCO) compliance, ensuring that contributors certify they have the right to submit their contributions under the project's license.
# custom:
#   severity: MEDIUM
#   remediationSteps:
#     - 1. Make sure you have admin permissions
#     - 2. Go to the organization settings page
#     - 3. Scroll down to the 'Commit signoff' section
#     - 4. Check the box for 'Require contributors to sign off on web-based commits'
#     - 5. Click 'Save'
#   requiredScopes: [admin:org]
#   threat:
#     - Without requiring commit signoff, there is no explicit certification that contributors have the legal right to submit their code under the project's license.
#     - This could lead to intellectual property issues and compliance violations.
#     - The Developer Certificate of Origin (DCO) provides a lightweight mechanism for tracking contribution provenance.
#   compliance:
#     - title: SOC 2
#       controls: [CC4.1, CC4.2, CC5.2, CC5.3, CC6.1, CC6.2, CC6.3, CC6.6, CC6.7, CC7.1, CC7.2, CC7.3, CC8.1, CC9.2]
#     - title: ISO 27002:2022
#       controls: [5.2, 5.3, 5.12, 5.15, 5.18, 5.21, 5.24, 5.25, 5.26, 5.28, 5.33, 5.36, 5.37, 8.2, 8.3, 8.4, 8.6, 8.8, 8.9, 8.10, 8.12, 8.15, 8.16, 8.25, 8.28, 8.29, 8.32]
default web_commit_signoff_required := false

web_commit_signoff_required := true {
	not input.web_commit_signoff_required
}
