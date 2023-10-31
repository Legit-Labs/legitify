package organization

import data.common.webhooks as webhookUtils

# METADATA
# scope: rule
# title: Webhooks Should Be Configured With A Secret
# description: Webhooks are not configured with a shared secret to validate the origin and content of the request. This could allow your webhook to be triggered by any bad actor with the URL.
# custom:
#   requiredEnrichers: [hooksList]
#   severity: LOW
#   remediationSteps: [Make sure you have admin permissions, Go to the organization settings page, Select "Webhooks", Press on the insecure webhook, Configure a secret , Click "Update webhook"]
#   requiredScopes: [admin:org_hook]
#   threat:
#     - "Not using a webhook secret makes the service receiving the webhook unable to determine the authenticity of the request."
#     - "This allows attackers to masquerade as your organization, potentially creating an unstable or insecure state in other systems."
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
#   remediationSteps: [Make sure you have admin permissions, Go to the organization settings page, Select "Webhooks", Press on the insecure webhook, Verify url starts with https, Enable "SSL verification" , Click "Update webhook"]
#   requiredScopes: [admin:org_hook]
#   threat:
#     - "If SSL verification is disabled, any party with access to the target DNS domain can masquerade as your designated payload URL, allowing it freely read and affect the response of any webhook request."
#     - "In the case of GitHub Enterprise Server instances, it may be sufficient only to control the DNS configuration of the network where the instance is deployed, as an attacker can redirect traffic to the target domain in your internal network directly to them, and this is often much easier than compromising an internet-facing domain."
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
#   remediationSteps: [Make sure you have admin permissions, Go to the organization settings page, Enter "Authentication security" tab, Under "Two-factor authentication", Toggle on "Require two-factor authentication for everyone in the <ORG> organization", Click "Save"]
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
#   The organization should be configured to prevent non-admin members creating public repositories. Creating a public repository may expose sensitive organization code, which, once exposed, may be copied, cached or stored by external parties. Therefore, it is highly recommended to restrict the option to create public repositories to admins only and reduce the risk of unintentional code exposure. NOTE: You should also verify that repositories owners can't change existing repositories visibility to be public. If allowed, a malicious user could create a private repo and change it to public. See: https://docs.github.com/en/enterprise-cloud@latest/organizations/managing-organization-settings/restricting-repository-visibility-changes-in-your-organization for further information
# custom:
#   severity: MEDIUM
#   remediationSteps: [Make sure you have admin permissions, Go to the organization settings page, Enter "Member privileges" tab, Under "Repository creation", Toggle off "Public", Click "Save"]
#   requiredScopes: [read:org]
#   threat:
#     - "A member of the organization could inadvertently or maliciously make public an internal repository exposing confidential data."
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
#   remediationSteps: [Make sure you have admin permissions, Go to the organization settings page, Enter "Member privileges" tab, Under "Base permissions", Set permissions to "No permissions", Click "Save"]
#   requiredScopes: [read:enterprise]
#   threat:
#     - "Organization members can see the content of freshly created repositories, even if they should be restricted."
default default_repository_permission_is_not_none := true

default_repository_permission_is_not_none := false {
	input.organization.default_repository_permission == "none"
}

# METADATA
# scope: rule
# custom:
#   severity: MEDIUM
# title: Organization Should Use Single-Sign-On
# description: It is recommended to enable access to an organization via SAML single sign-on (SSO) by authenticating through an identity provider (IdP). This allows for central account control and for timely access revocations.
# custom:
#   remediationSteps: [Make sure you have admin permissions, Go to the organization settings page, Enter "Authentication security" tab, Toggle on "Enable SAML authentication", Fill in the remaining SSO configuration as instructed on the screen, Click "Save"]
#   requiredScopes: [admin:org]
#   threat: Not using an SSO solution makes it more difficult to track a potentially compromised user's actions accross different systems, prevents the organization from defining a common password policy, and makes it challenging to audit different aspects of the user's behavior.
default organization_not_using_single_sign_on := true

organization_not_using_single_sign_on := false {
	input.saml_enabled
}
