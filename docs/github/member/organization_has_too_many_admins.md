---
layout: default
title: Organization Should Have Fewer Than Three Owners
parent: Member Policies
grand_parent: GitHub Policies
---


## Organization Should Have Fewer Than Three Owners
policy name: organization_has_too_many_admins

severity: MEDIUM

### Description
Organization owners are highly privileged and could create great damage if they are compromised. It is recommended to limit the number of Organizational Admins to the minimum needed (recommended maximum 3 owners).

### Threat Example(s)
1. An organization has a permissive attitude and provides an owner role to all developers
2. One of the developers has decided to collaborate with an evil ransomware gang, and uses his high privileges to add a malicious external collaborator
3. The malicious collaborator, being an owner, has a wide range of destructive operations he can do (e.g. remove security settings)



### Remediation
1. Make sure you have admin permissions
2. Go to the organization People page
3. Select the unwanted owners
4. Using the 'X members selected' - change role to member



