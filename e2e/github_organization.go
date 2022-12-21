package test

var testCasesGitHubOrganization = []testCase{
	{
		path:         "data.organization.non_admins_can_create_public_repositories",
		failedEntity: "Legitify-E2E-2",
		passedEntity: "Legitify-E2E",
	},
	{
		path:         "data.organization.default_repository_permission_is_not_none",
		failedEntity: "Legitify-E2E-2",
		passedEntity: "Legitify-E2E",
	},
	{
		path:         "data.organization.two_factor_authentication_not_required_for_org",
		failedEntity: "Legitify-E2E-2",
		passedEntity: "Legitify-E2E",
	},
	{
		path:         "data.organization.organization_webhook_no_secret",
		failedEntity: "Legitify-E2E-2",
		passedEntity: "Legitify-E2E",
	},
	{
		path:         "data.organization.organization_webhook_doesnt_require_ssl",
		failedEntity: "Legitify-E2E-2",
		passedEntity: "Legitify-E2E",
	},
}

var testCasesGitHubActions = []testCase{
	{
		path:         "data.actions.token_default_permissions_is_read_write",
		failedEntity: "Legitify-E2E-2",
		passedEntity: "Legitify-E2E",
	},
	{
		path:         "data.actions.all_repositories_can_run_github_actions",
		failedEntity: "Legitify-E2E-2",
		passedEntity: "Legitify-E2E",
	},
	{
		path:         "data.actions.all_github_actions_are_allowed",
		failedEntity: "Legitify-E2E-2",
		passedEntity: "Legitify-E2E",
	},
	{
		path:         "data.actions.actions_can_approve_pull_requests",
		failedEntity: "Legitify-E2E-2",
		passedEntity: "Legitify-E2E",
	},
}

var testCasesGitHubRunnerGroup = []testCase{
	{
		path:         "data.runner_group.runner_group_can_be_used_by_public_repositories",
		failedEntity: "test fail",
		passedEntity: "test pass",
	},
	{
		path:         "data.runner_group.runner_group_not_limited_to_selected_repositories",
		failedEntity: "test fail",
		passedEntity: "test pass",
	},
}
