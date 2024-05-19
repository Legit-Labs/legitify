package member

import data.common.members as memberUtils
# METADATA
# scope: rule
# title: Organization Should Have Fewer Than Three Owners
# description: Organization owners are highly privileged and could create great damage if they are compromised. It is recommended to limit the number of Organizational Admins to the minimum needed (recommended maximum 3 owners).
# custom:
#   remediationSteps:
#     - 1. Make sure you have admin permissions
#     - 2. Go to the organization People page
#     - 3. Select the unwanted owners
#     - 4. Using the 'X members selected' - change role to member
#   severity: MEDIUM
#   requiredScopes: [admin:org]
#   threat:
#     - 1. An organization has a permissive attitude and provides an owner role to all developers
#     - 2. One of the developers has decided to collaborate with an evil ransomware gang, and uses his high privileges to add a malicious external collaborator
#     - 3. The malicious collaborator, being an owner, has a wide range of destructive operations he can do (e.g. remove security settings)
default organization_has_too_many_admins := true

organization_has_too_many_admins := false {
	admins := count([member | member := input.members[_]; member.is_admin == true])
	admins <= 3
}

# METADATA
# scope: rule
# title: Organization Members Should Have Activity In The Last 6 Months
# description: A member did not perform any action in the last 6 months. Stale members can pose a potential risk if they are compromised. Consider removing the user's access completely.
# custom:
#   requiredEnrichers: [entityId, violatedUsers]
#   remediationSteps:
#     - 1. Make sure you have admin permissions
#     - 2. Go to the org's People page
#     - 3. Select all stale members
#     - 4. Using the 'X members selected' - remove members from organization
#   severity: LOW
#   requiredScopes: [admin:org]
#   prerequisites: [premium]
#   threat:
#     - Stale members are most likely not managed and monitored, increasing the possibility of being compromised.
stale_member_found[mem] := true {
	some member
	mem := input.members[member]
	mem.is_admin == false
	mem.last_active != -1
	memberUtils.isStale(mem.last_active, 6)
}

# METADATA
# scope: rule
# title: Organization Admins Should Have Activity In The Last 6 Months
# description: A member with organizational admin permissions did not perform any action in the last 6 months. Admin users are extremely powerful and common compliance standards demand keeping the number of admins to a minimum. Consider revoking this member’s admin credentials by downgrading to regular user or removing the user completely.
# custom:
#   requiredEnrichers: [entityId, violatedUsers]
#   remediationSteps:
#     - 1. Make sure you have admin permissions
#     - 2. Go to the org's People page
#     - 3. Select all stale admins
#     - 4. Using the 'X members selected' - remove members from organization
#   severity: MEDIUM
#   requiredScopes: [admin:org]
#   prerequisites: [premium]
#   threat:
#     - Stale admins are most likely not managed and monitored, increasing the possibility of being compromised.
stale_admin_found[mem] := true {
	some member
	mem := input.members[member]
	mem.is_admin == true
	mem.last_active != -1
	memberUtils.isStale(mem.last_active, 6)
}