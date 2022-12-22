package test

import (
	"flag"
	"github.com/thedevsaddam/gojsonq/v2"
	"testing"
)

var reportPath = flag.String("report_path", "/tmp/out.json", "legitify report output path")

func TestGitHub(t *testing.T) {
	tests := [][]testCase{
		testCasesGitHubOrganization,
		testCasesGitHubActions,
		testCasesGitHubRunnerGroup,
		testCasesGitHubRepository,
	}

	for _, testCases := range tests {
		AssertionLoop(t, testCases)
	}
}

func AssertionLoop(t *testing.T, tests []testCase) {
	jq := gojsonq.New(gojsonq.SetSeparator("->")).File(*reportPath)
	for _, test := range tests {
		t.Logf("Testing: %s", test.path)
		testFormattedPath := test.path + "->violations"

		if test.passedEntity != "" {
			res := jq.From(testFormattedPath).Where("aux->entityName", "=", test.passedEntity).Where("Status", "=", "PASSED").Count()

			if res != 1 {
				t.Logf("Failed on test %s, Entity %s did not pass", test.path, test.passedEntity)
				t.Fail()
			}
		}

		jq.Reset()

		if test.failedEntity != "" {
			res := jq.From(testFormattedPath).Where("aux->entityName", "=", test.failedEntity).Where("Status", "=", "FAILED").Count()

			if res != 1 {
				t.Logf("Failed on test: %s, Entity: %s did not failed", test.path, test.failedEntity)
				t.Fail()
			}
		}

		jq.Reset()

		if test.skippedEntity != "" {
			res := jq.From(testFormattedPath).Where("aux->entityName", "=", test.skippedEntity).Where("Status", "=", "SKIPPED").Count()

			if res != 1 {
				t.Logf("Failed on test: %s, Entity: %s did not skip", test.path, test.skippedEntity)
				t.Fail()
			}
		}

		jq.Reset()
	}
}

func TestGitLab(t *testing.T) {
	tests := []testCase{
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
			path:          "data.member.two_factor_authentication_is_disabled_for_a_collaborator",
			skippedEntity: "https://gitlab.com/legitify-test",
		},
		{
			path:          "data.member.two_factor_authentication_is_disabled_for_an_external_collaborator",
			skippedEntity: "https://gitlab.com/legitify-test",
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
	}
	AssertionLoop(t, tests)
}
