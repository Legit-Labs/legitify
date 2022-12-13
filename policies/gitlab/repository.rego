package repository

# METADATA
# scope: rule
# title: Repository not maintained
# description: The project was not active in the last 3 months. A project which is not active might not be patched against security issues within its code and dependencies, and is therefore at higher risk of including unpatched vulnerabilities.
# custom:
#   remediationSteps: [Make sure you have admin permissions, Either Delete or Archive the repository]
#   severity: HIGH
default repository_not_maintained = false

repository_not_maintained {
    input.archived == false
    ns := time.parse_rfc3339_ns(input.last_activity_at)
    now := time.now_ns()
    diff := time.diff(now, ns)
    monthsIndex := 1
    inactivityMonthsThreshold := 3
    diff[monthsIndex] >= inactivityMonthsThreshold
}

repository_not_maintained {
    input.archived == false
    ns := time.parse_rfc3339_ns(input.last_activity_at)
    now := time.now_ns()
    diff := time.diff(now, ns)
    yearIndex := 0
    diff[yearIndex] > 0
}



# METADATA
# scope: rule
# title: Project Has Too Many Owners
# description: Projects' owners are highly privileged and could create great damage if being compromised, it's recommeneded to limit them to the minimum required (recommended maximum 3 admins).
# custom:
#   severity: LOW
#   remediationSteps: [Make sure you have owner permissions, Go to the Project Information -> Members page, Select the unwanted owner users and remove the selected owners]
default repository_has_too_many_admins  = false
repository_has_too_many_admins {
    admins := [admin | admin := input.members[_]; admin.access_level == 50]
    count(admins) > 3
}

# METADATA
# scope: rule
# title: Forking Allowed for This Repository
# description: Forking a repository can lead to loss of control and potential exposure of the source code. The option to fork must be disabled by default and turned on only by owners deliberately when opting to create a fork. If you do not need forking, it is recommended to turn it off in the project's configuration.
# custom:
#   remediationSteps: [Make sure you have owner permissions, Go to the project's settings page, Enter "General" tab, Under "Visibility, project features, permissions", Toggle off "Forks"]
#   severity: LOW
#   threat:
#    - "A user with permissions to the repository could intentionally/accidentally fork a private repository, make it public and cause a code-leak incident"
default forking_allowed_for_repository = false
forking_allowed_for_repository {
    input.public == false
    input.forking_access_level == "enabled"
}

# METADATA
# scope: rule
# title: Default Branch Is Not Protected
# description: Branch protection is not enabled for this repository’s default branch. Protecting branches ensures new code changes must go through a controlled merge process and allows enforcement of code review as well as other security tests. This issue is raised if the default branch protection is turned off.
# custom:
#   remediationSteps: [Make sure you have owner permissions, Go to the projects's settings -> Repository page, Enter "Protected branches" tab, select the default branch. Set the allowed to merge to "maintainers" and the allowed to push to "No one"]
#   severity: MEDIUM
#   threat:
#    - "Users can merge code without being reviewed which can lead to insecure code reaching the main branch and production."
default missing_default_branch_protection = false
missing_default_branch_protection {
    default_protected_branches := [protected_branch | protected_branch := input.protected_branches[_]; protected_branch.name == input.default_branch]
    count(default_protected_branches) == 0
}


# METADATA
# scope: rule
# title: Default Branch Allows Force Pushes
# description: The history of the default branch is not protected against changes for this repository. Protecting branch history ensures every change that was made to code can be retained and later examined. This issue is raised if the default branch history can be modified using force push.
# custom:
#   remediationSteps: [Make sure you have owner permissions, Go to the projects's settings -> Repository page, Enter "Protected branches" tab, select the default branch. Set the allowed to merge to "maintainers" and the allowed to push to "No one"]
#   severity: MEDIUM
default missing_default_branch_protection_force_push = false

missing_default_branch_protection_force_push {
    missing_default_branch_protection
}

missing_default_branch_protection_force_push {
    default_protected_branches := [protected_branch | protected_branch := input.protected_branches[_]; protected_branch.name == input.default_branch]
    rules_allow_force_push := [rule_allow_force_push | rule_allow_force_push := default_protected_branches[_]; rule_allow_force_push.allow_force_push == true]
	count(rules_allow_force_push) > 0
}

# METADATA
# scope: rule
# title: Webhook Configured Without SSL Verification
# description: Webhooks that are not configured with SSL verification enabled could expose your sofware to man in the middle attacks (MITM).
# custom:
#   severity: LOW
#   remediationSteps: [Make sure you can manage webhooks for the repository, Go to the repository settings page, Select "Webhooks", Press on the "Enable SSL verfication", Click "Save changes"]
#   threat:
#   - "Attackers could take over a domain or perform a MiTM attack to get data and secrets sent to the webhook"
default repository_webhook_doesnt_require_ssl = false
repository_webhook_doesnt_require_ssl = true {
    webhooks_without_ssl_verification := [webhook_without_verification | webhook_without_verification := input.webhooks[_]; webhook_without_verification.enable_ssl_verification == false]
    count(webhooks_without_ssl_verification) > 0
}

# METADATA
# scope: rule
# title: Project Doesn’t Require All Pipelines to Succeed
# description: the checks which validate the quality and security of the code are not required to pass before submitting new changes. It is advised to turn this control on to ensure any existing or future check will be required to pass
# custom:
#   severity: MEDIUM
#   remediationSteps: [Make sure you can manage project merge requests permissions, Go to the project's settings page, Select "Merge Requests", Press on the "Pipelines must succeed", Click "Save changes"]
#   threat:
#     - "Users could merge its code without all required checks passes what could lead to insecure code reaching your main branch and production."
default requires_status_checks = false
requires_status_checks = true {
    input.only_allow_merge_if_pipeline_succeeds == false
}


# METADATA
# scope: rule
# title: Project Doesn't Require All Conversations To Be Resolved Before Merge
# description: Require all merge request conversations to be resolved before merging. Check this to avoid bypassing/missing a Pull Reuqest comment.
# custom:
#   severity: LOW
#   remediationSteps: [Make sure you can manage project merge requests permissions, Go to the project's settings page, Select "Merge Requests", Press on the "All threads must be resolved", Click "Save changes"]
default no_conversation_resolution = false
no_conversation_resolution = true {
    input.only_allow_merge_if_all_discussions_are_resolved == false
}
