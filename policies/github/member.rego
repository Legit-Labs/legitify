package member

# METADATA
# scope: rule
# title: Organization Has Too Many Owners
# description: Organization owners are highly privileged and could create great damage if being compromised, it's recommended to limit them to the minimum needed (recommended maximum 3 owners).
# custom:
#   remediationSteps: [Make sure you have admin permissions, Go to the organization People page, Select the unwanted owners, Using the "X members selected" - change role to member]
#   severity: MEDIUM
#   requiredScopes: [admin:org]
#   threat:
#     - "1. An organization has a permissive attitude and provides an owner role to all developers."
#     - "2. One of the developers has decided to collaborate with an evil ransomware gang, and uses his high privileges to add a malicious external collaborator"
#     - "3. The malicious collaborator, being an owner, has a wide range of destructive operations he can do (e.g. remove security settings)"
default organization_has_too_many_admins = false
organization_has_too_many_admins {
    admins := count({member | member := input.members[_]; member.is_admin == false})
    admins > 3
}

# METADATA
# scope: rule
# title: Stale Member Found
# description: A member didn't do any action in the last 6 months. Stale members can pose a potential risk if they are compromised. Consider removing the user's access completely.
# custom:
#   requiredEnrichers: [entityId, violatedUsers]
#   remediationSteps: [Make sure you have admin permissions, Go to the org's People page, Select all stale members, Using the "X members selected" - remove members from organization]
#   severity: LOW
#   requiredScopes: [admin:org]
#   prerequisites: [premium]
#   threat:
#     - "Stale members are most likely not managed and monitored, increasing the possibility of being compromised."
stale_member_found[mem] = true {
    some member
    mem := input.members[member]
    mem.is_admin == false
    isStale(mem.last_active, 6)
}

# METADATA
# scope: rule
# custom:
# title: Stale Admin Found
# description: A member with global admin permissions without any activity in the past 6 months. Admin users are extremely powerful and common compliance standards demand keeping the number of admins to a minimum. Consider revoking this member’s admin credentials by downgrading to regular user or removing the user completely.
# custom:
#   requiredEnrichers: [entityId, violatedUsers]
#   remediationSteps: [Make sure you have admin permissions, Go to the org's People page, Select all stale admins, Using the "X members selected" - remove members from organization]
#   severity: MEDIUM
#   requiredScopes: [admin:org]
#   prerequisites: [premium]
#   threat:
#     - "Stale admins are most likely not managed and monitored, increasing the possibility of being compromised."
stale_admin_found[mem] = true {
    some member
    mem := input.members[member]
    mem.is_admin == true
    isStale(mem.last_active, 6)
}

isStale(target_last_active, count_months) {
    now := time.now_ns()
    diff := time.diff(now, target_last_active)
    # diff[1] the months index
    diff[1] >= count_months
}
