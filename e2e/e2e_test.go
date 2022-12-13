package test

import (
	"flag"
	"github.com/thedevsaddam/gojsonq/v2"
	"testing"
)

var reportPath = flag.String("report_path", "/tmp/out.json", "legitify report output path")

func TestGitHub(t *testing.T) {
	tests := []struct {
		path         string
		failedEntity string
		passedEntity string
	}{
		// actions tests
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
		// runner group tests
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
		// organization tests
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
	AssertionLoop(t, tests)
}

func AssertionLoop(t *testing.T, tests []struct {
	path         string
	failedEntity string
	passedEntity string
}) {
	for _, test := range tests {
		t.Logf("Testing: %s", test.path)
		testFormattedPath := test.path + "->violations"
		jq := gojsonq.New(gojsonq.SetSeparator("->")).File(*reportPath)
		res := jq.From(testFormattedPath).Where("aux->entityName", "=", test.passedEntity).Where("Status", "=", "PASSED").Count()

		if res != 1 {
			t.Logf("Failed on test %s, Entity %s did not pass", test.path, test.passedEntity)
			t.Fail()
		}
		jq = gojsonq.New(gojsonq.SetSeparator("->")).File(*reportPath)
		res = jq.From(testFormattedPath).Where("aux->entityName", "=", test.failedEntity).Where("Status", "=", "FAILED").Count()

		if res != 1 {
			t.Logf("Failed on test: %s, Entity: %s did not failed", test.path, test.failedEntity)
			t.Fail()
		}
	}
}

func TestGitLab(t *testing.T) {
	tests := []struct {
		path         string
		failedEntity string
		passedEntity string
	}{
		{
			path:         "data.member.two_factor_authentication_is_disabled_for_a_collaborator",
			failedEntity: "Legitify-E2E-2",
			passedEntity: "Legitify-E2E",
		},
		{
			path:         "data.member.two_factor_authentication_is_disabled_for_an_external_collaborator",
			failedEntity: "Legitify-E2E-2",
			passedEntity: "Legitify-E2E",
		},
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
	}
	AssertionLoop(t, tests)

}
