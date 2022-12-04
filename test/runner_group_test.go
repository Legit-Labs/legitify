package test

import (
	"testing"

	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/google/go-github/v44/github"
)

type runnerGroupMockConfiguration struct {
	allowedByPublic bool
}

func newRunnerGroupMock(config runnerGroupMockConfiguration) githubcollected.RunnerGroup {
	return githubcollected.RunnerGroup{
		Organization: defaultOrg,
		RunnerGroup: &github.RunnerGroup{
			AllowsPublicRepositories: &config.allowedByPublic,
		},
	}
}

func TestRunnerGroup(t *testing.T) {
	tests := []struct {
		name             string
		policyName       string
		shouldBeViolated bool
		args             runnerGroupMockConfiguration
	}{
		{
			name:             "runner group is allowed to run by public repositories",
			policyName:       "runner_group_can_be_used_by_public_repositories",
			shouldBeViolated: true,
			args: runnerGroupMockConfiguration{
				allowedByPublic: true,
			},
		},
		{
			name:             "runner group is not allowed to run by public repositoreis",
			policyName:       "runner_group_can_be_used_by_public_repositories",
			shouldBeViolated: false,
			args: runnerGroupMockConfiguration{
				allowedByPublic: false,
			},
		},
	}

	for _, test := range tests {
		PolicyTestTemplateGitHub(t, test.name, newRunnerGroupMock(test.args),
			namespace.RunnerGroup, test.policyName, test.shouldBeViolated)
	}
}
