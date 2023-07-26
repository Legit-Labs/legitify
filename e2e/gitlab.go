package test

var testCasesGitLab = []testCase{
	{
		path:         "data.organization.two_factor_authentication_not_required_for_group",
		failedEntity: "Legitify-E2E-2",
		passedEntity: "Legitify-E2E",
	},
	{
		path:         "data.organization.collaborators_can_fork_repositories_to_external_namespaces",
		failedEntity: "Legitify-E2E-2",
		passedEntity: "Legitify-E2E",
	},
	{
		path:         "data.organization.organization_webhook_doesnt_require_ssl",
		failedEntity: "Legitify-E2E-2",
		passedEntity: "Legitify-E2E",
	},
	{
		path:         "data.organization.group_does_not_enforce_branch_protection_by_default",
		failedEntity: "Legitify-E2E-2",
		passedEntity: "Legitify-E2E",
	},
	{
		path:         "data.repository.missing_default_branch_protection",
		failedEntity: "failed_repo",
		passedEntity: "passed_repo",
	},
	{
		path:         "data.repository.code_review_by_two_members_not_required",
		failedEntity: "failed_repo",
		passedEntity: "passed_repo",
	},
	{
		path:         "data.repository.code_review_not_required",
		failedEntity: "failed_repo",
		passedEntity: "passed_repo",
	},
	{
		path:         "data.repository.repository_allows_review_requester_to_approve_their_own_request",
		failedEntity: "failed_repo",
		passedEntity: "passed_repo",
	},
	{
		path:         "data.repository.repository_allows_overriding_approvers",
		failedEntity: "failed_repo",
		passedEntity: "passed_repo",
	},
	{
		path:         "data.repository.repository_require_code_owner_reviews_policy",
		failedEntity: "failed_repo",
		passedEntity: "passed_repo",
	},
	{
		path:         "data.repository.repository_allows_committer_approvals_policy",
		failedEntity: "failed_repo",
		passedEntity: "passed_repo",
	},
	{
		path:         "data.repository.repository_dismiss_stale_reviews",
		failedEntity: "failed_repo",
		passedEntity: "passed_repo",
	},
	{
		path:          "data.member.two_factor_authentication_is_disabled_for_a_collaborator",
		skippedEntity: "legitify-test",
	},
	{
		path:          "data.member.two_factor_authentication_is_disabled_for_an_external_collaborator",
		skippedEntity: "legitify-test",
	},
}
