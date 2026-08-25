package test

import (
	"github.com/Legit-Labs/legitify/internal/common/scm_type"
	"testing"

	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/google/go-github/v53/github"
)

type runnerGroupMockConfiguration struct {
	allowedByPublic bool
	visibility      string
}

func newRunnerGroupMock(config runnerGroupMockConfiguration) githubcollected.RunnerGroup {
	return githubcollected.RunnerGroup{
		Organization: defaultOrg,
		RunnerGroup: &github.RunnerGroup{
			AllowsPublicRepositories: &config.allowedByPublic,
			Visibility:               &config.visibility,
		},
	}
}

func TestRunnerGroup(t *testing.T) {
	tests := []struct {
		name             string
		policyName       string
		scmType          scm_type.ScmType
		shouldBeViolated bool
		args             runnerGroupMockConfiguration
	}{
		{
			name:             "runner group is allowed to run by public repositories",
			policyName:       "runner_group_can_be_used_by_public_repositories",
			scmType:          scm_type.GitHub,
			shouldBeViolated: true,
			args: runnerGroupMockConfiguration{
				allowedByPublic: true,
			},
		},
		{
			name:             "runner group is not allowed to run by public repositoreis",
			policyName:       "runner_group_can_be_used_by_public_repositories",
			scmType:          scm_type.GitHub,
			shouldBeViolated: false,
			args: runnerGroupMockConfiguration{
				allowedByPublic: false,
			},
		},
		{
			name:             "runner group is limited to selected repositories",
			policyName:       "runner_group_not_limited_to_selected_repositories",
			scmType:          scm_type.GitHub,
			shouldBeViolated: false,
			args: runnerGroupMockConfiguration{
				visibility: "selected",
			},
		},
		{
			name:             "runner group is available to all repositories",
			policyName:       "runner_group_not_limited_to_selected_repositories",
			scmType:          scm_type.GitHub,
			shouldBeViolated: true,
			args: runnerGroupMockConfiguration{
				visibility: "all",
			},
		},
	}

	for _, test := range tests {
		PolicyTestTemplate(t, test.name, newRunnerGroupMock(test.args),
			namespace.RunnerGroup, test.policyName, test.shouldBeViolated, test.scmType)
	}
}
