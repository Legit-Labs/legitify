package enterprise

# METADATA
# scope: rule
# custom:
#   severity: MEDIUM
# title: Enterprise Should Not Allow Members To Change Repository Visibility
# description: It is recommended to change the enterprise's Repository visibility change policy. This will prevents users from creating private repositories and change them to be public. Malicous actors could leak code if enabled.
# custom:
#   remediationSteps: [Make sure you are an enterprise owner, Go to the policies page, Under the "Repository visibility change" section, choose the "Disabled" option]
#   requiredScopes: [admin:enterprise]
#   threat:
#     - "A member of the organization could inadvertently or maliciously make public an internal repository exposing confidential data."
default enterprise_not_using_visibility_change_disable_policy = true

enterprise_not_using_visibility_change_disable_policy = false {
	input.visibility_change_disabled == "DISABLED"
}