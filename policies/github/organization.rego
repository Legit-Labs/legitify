package organization
import data.common.webhooks as webhookUtils

# METADATA
# scope: rule
# title: Webhook Configured Without A Secret
# description: Webhooks that are not configured with a token authenticated to validate the origin of the request and could make your software vulnerable.
# custom:
#   requiredEnrichers: [violatedHooks]
#   severity: LOW
#   remediationSteps: [Make sure you have admin permissions, Go to the organization settings page, Select "Webhooks", Press on the insecure webhook, Confiure a secret , Click "Update webhook"]
#   requiredScopes: [admin:org_hook]
organization_webhook_no_secret[hook] = true {
    some index
    hook := input.hooks[index]
    not webhookUtils.has_secret(hook)
}

# METADATA
# scope: rule
# title: Webhook Configured Without SSL
# description: Webhooks that are not configured with SSL enabled could expose your sofware to man in the middle attacks (MITM).
# custom:
#   requiredEnrichers: [violatedHooks]
#   severity: LOW
#   remediationSteps: [Make sure you have admin permissions, Go to the organization settings page, Select "Webhooks", Press on the insecure webhook, Verify url starts with https, Enable "SSL verfication" , Click "Update webhook"]
#   requiredScopes: [admin:org_hook]
organization_webhook_doesnt_require_ssl[hook] = true {
    some index
    hook := input.hooks[index]
    not webhookUtils.ssl_enabled(hook)
}

# METADATA
# scope: rule
# title: Two-Factor Authentication Is Not Enforced For The Organization
# description: The two-factor authentication requirement is not enabled at the organization level. Regardless of whether users are managed externally by SSO, it is highly recommended to enable this option, to reduce the risk of a deliberate or accidental user creation without MFA.
# custom:
#   severity: HIGH
#   remediationSteps: [Make sure you have admin permissions, Go to the organization settings page, Enter "Authentication security" tab, Under "Two-factor authentication", Toggle on "Require two-factor authentication for everyone in the <ORG> organization", Click "Save"]
#   requiredScopes: [admin:org]
#   threat:
#     - If an attacker gets the valid credentials for one of the organization’s users they can authenticate to your GitHub organization.
default two_factor_authentication_not_required_for_org  = false
two_factor_authentication_not_required_for_org {
    input.organization.two_factor_requirement_enabled == false
}

# METADATA
# scope: rule
# title: Non-Admins Can Create Public Repositories
# description: An organization allows non-admin members to create public repositories. Creating a public repository can be done by mistake, and may expose sensitive organization code, which, once exposed, may be copied, cached or stored by external parties. Therefore, it is highly recommended to restrict the option to create public repositories to admins only and reduce the risk of unintentional code exposure.
# custom:
#   severity: MEDIUM
#   remediationSteps: [Make sure you have admin permissions, Go to the organization settings page, Enter "Member privileges" tab, Under "Repository creation", Toggle off "Public", Click "Save"]
#   requiredScopes: [read:org]
#   threat:
#     - "A member of the organization could inadvertently or maliciously make public an internal repository exposing confidential data."
default non_admins_can_create_public_repositories = false
non_admins_can_create_public_repositories {
    input.organization.members_can_create_public_repositories == true
}

# METADATA
# scope: rule
# title: Permissive Default Member Permissions Exist For New Repositories
# description: Default repository permissions configuration is not set in the organization, thus every new repository will be accessible by default to all users. It is strongly recommended to remove the default permissions and assign them on demand.
# custom:
#   severity: HIGH
#   remediationSteps: [Make sure you have admin permissions, Go to the organization settings page, Enter "Member privileges" tab, Under "Base permissions", Set permissions to "No permissions", Click "Save"]
#   requiredScopes: [read:enterprise]
#   threat:
#     - "Organization members can see the content of freshly created repositories, even if they should be restricted."
default default_repository_permission_is_not_none = false
default_repository_permission_is_not_none {
    input.organization.default_repository_permission != "none"
}

# METADATA
# scope: rule
# custom:
#   severity: MEDIUM
# title: Organization Not Using Single-Sign-On
# description: It is recommended to enable access to an organization via SAML single sign-on (SSO) by authenticating through an identity provider (IdP).
# custom:
#   remediationSteps: [Make sure you have admin permissions, Go to the organization settings page, Enter "Authentication security" tab, Toggle on "Enable SAML authentication", Fill in the remaining SSO configuration as instructed on the screen, Click "Save"]
#   requiredScopes: [admin:org]
default organization_not_using_single_sign_on = false
organization_not_using_single_sign_on {
    input.saml_enabled == false
}
