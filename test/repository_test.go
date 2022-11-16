package test

import (
	"github.com/Legit-Labs/legitify/internal/clients/github/types"
	"testing"

	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/google/go-github/v44/github"
)

func repositoryTestTemplate(t *testing.T, name string, mockData interface{}, testedPolicyName string, expectFailure bool) {
	ns := namespace.Repository
	PolicyTestTemplate(t, name, mockData, ns, testedPolicyName, expectFailure)
}

var bools = []bool{true, false}

func makeRepo(repo githubcollected.GitHubQLRepository) githubcollected.Repository {
	return githubcollected.Repository{
		Repository: &repo,
	}
}

func makeRepoForBranch(branch githubcollected.GitHubQLBranch) githubcollected.Repository {
	return makeRepo(githubcollected.GitHubQLRepository{
		Name:             "REPO",
		DefaultBranchRef: &branch,
	})
}

func makeRepoForBranchProtection(prot githubcollected.GitHubQLBranchProtectionRule) githubcollected.Repository {
	return makeRepoForBranch(githubcollected.GitHubQLBranch{
		BranchProtectionRule: &prot,
	})
}

func TestRepositoryBranchProtection(t *testing.T) {
	name := "repository should have branch protection"
	testedPolicyName := "missing_default_branch_protection"
	makeMockData := makeRepoForBranch

	branches := []githubcollected.GitHubQLBranch{
		{},
		{
			BranchProtectionRule: &githubcollected.GitHubQLBranchProtectionRule{
				AllowsDeletions:   github.Bool(false),
				AllowsForcePushes: github.Bool(false),
			},
		},
	}

	for i, flag := range bools {
		repositoryTestTemplate(t, name, makeMockData(branches[i]), testedPolicyName, flag)
	}
}

func TestRepositoryForcePush(t *testing.T) {
	name := "repository should have branch protection: force push"
	testedPolicyName := "missing_default_branch_protection_force_push"
	makeMockData := func(flag bool) githubcollected.Repository {
		return makeRepoForBranchProtection(githubcollected.GitHubQLBranchProtectionRule{
			AllowsDeletions:   github.Bool(false),
			AllowsForcePushes: github.Bool(flag),
		})
	}

	for _, flag := range bools {
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, flag)
	}
}

func TestRepositoryAllowDeletion(t *testing.T) {
	name := "repository should have branch protection: delete"
	testedPolicyName := "missing_default_branch_protection_deletion"
	makeMockData := func(flag bool) githubcollected.Repository {
		return makeRepoForBranchProtection(githubcollected.GitHubQLBranchProtectionRule{
			AllowsDeletions:   github.Bool(flag),
			AllowsForcePushes: github.Bool(false),
		})
	}

	for _, flag := range bools {
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, flag)
	}
}

func TestRepositoryCodeReview(t *testing.T) {
	name := "repository should have code review required"
	testedPolicyName := "code_review_by_two_members_not_required"
	makeMockData := func(count int) githubcollected.Repository {
		return makeRepoForBranchProtection(githubcollected.GitHubQLBranchProtectionRule{
			RequiredApprovingReviewCount: github.Int(count),
		})
	}
	counts := []int{
		1,
		2,
	}
	for i, flag := range bools {
		repositoryTestTemplate(t, name, makeMockData(counts[i]), testedPolicyName, flag)
	}
}
func TestRepositoryCodeOwnersOnly(t *testing.T) {
	name := "repository should have code review limited to owners only"
	testedPolicyName := "code_review_not_limited_to_code_owners"
	makeMockData := func(flag bool) githubcollected.Repository {
		return makeRepoForBranchProtection(githubcollected.GitHubQLBranchProtectionRule{
			RequiresCodeOwnerReviews: github.Bool(flag),
		})
	}
	for _, flag := range bools {
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, !flag)
	}
}

func TestRepositoryLinearHistory(t *testing.T) {
	name := "repository should require linear history"
	testedPolicyName := "non_linear_history"
	makeMockData := func(flag bool) githubcollected.Repository {
		return makeRepoForBranchProtection(githubcollected.GitHubQLBranchProtectionRule{
			RequiresLinearHistory: github.Bool(flag),
		})
	}
	for _, flag := range bools {
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, !flag)
	}
}

func TestRepositoryReviewDismissal(t *testing.T) {
	name := "repository should restrict review dismissals"
	testedPolicyName := "review_dismissal_allowed"
	makeMockData := func(flag bool) githubcollected.Repository {
		return makeRepoForBranchProtection(githubcollected.GitHubQLBranchProtectionRule{
			RestrictsReviewDismissals: github.Bool(flag),
		})
	}
	for _, flag := range bools {
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, !flag)
	}
}

func TestRepositoryRestrictPushes(t *testing.T) {
	name := "repository should restrict pushes"
	testedPolicyName := "pushes_are_not_restricted"
	makeMockData := func(flag bool) githubcollected.Repository {
		return makeRepoForBranchProtection(githubcollected.GitHubQLBranchProtectionRule{
			RestrictsPushes: github.Bool(flag),
		})
	}
	for _, flag := range bools {
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, !flag)
	}
}

func TestRepositoryRequireConversationResolution(t *testing.T) {
	name := "repository should require all conversations resolved before merge"
	testedPolicyName := "no_conversation_resolution"
	makeMockData := func(flag bool) githubcollected.Repository {
		return makeRepoForBranchProtection(githubcollected.GitHubQLBranchProtectionRule{
			RequiresConversationResolution: github.Bool(flag),
		})
	}
	for _, flag := range bools {
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, !flag)
	}
}

func TestRepositoryStaleReviews(t *testing.T) {
	name := "repository should not dismiss stale reviews"
	testedPolicyName := "dismisses_stale_reviews"
	makeMockData := func(flag bool) githubcollected.Repository {
		return makeRepoForBranchProtection(githubcollected.GitHubQLBranchProtectionRule{
			DismissesStaleReviews: github.Bool(flag),
		})
	}
	for _, flag := range bools {
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, !flag)
	}
}

func TestRepositoryStatusChecks(t *testing.T) {
	name := "repository should require status checks"
	testedPolicyName := "requires_status_checks"
	makeMockData := func(flag bool) githubcollected.Repository {
		return makeRepoForBranchProtection(githubcollected.GitHubQLBranchProtectionRule{
			RequiresStatusChecks: github.Bool(flag),
		})
	}
	for _, flag := range bools {
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, !flag)
	}
}

func TestRepositoryBranchesUpToDate(t *testing.T) {
	name := "repository should require branches to be up to date before merging"
	testedPolicyName := "requires_branches_up_to_date_before_merge"
	makeMockData := func(flag bool) githubcollected.Repository {
		return makeRepoForBranchProtection(githubcollected.GitHubQLBranchProtectionRule{
			RequiresStrictStatusChecks: github.Bool(flag),
		})
	}
	for _, flag := range bools {
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, !flag)
	}
}
func TestRepositorySignedCommits(t *testing.T) {
	name := "signed commits should be enabled"
	testedPolicyName := "no_signed_commits"
	makeMockData := func(flag bool) githubcollected.Repository {
		return makeRepoForBranchProtection(githubcollected.GitHubQLBranchProtectionRule{
			RequiresCommitSignatures: github.Bool(flag),
		})
	}
	for _, flag := range bools {
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, !flag)
	}
}
func TestRepositoryVulnerabilityAlerts(t *testing.T) {
	name := "vulnerability alerts not enabled"
	testedPolicyName := "vulnerability_alerts_not_enabled"
	makeMockData := func(flag *bool) githubcollected.Repository {
		return githubcollected.Repository{
			VulnerabilityAlertsEnabled: flag,
		}
	}

	options := map[bool][]*bool{
		true:  {github.Bool(false)},
		false: {nil, github.Bool(true)},
	}

	for _, expectFailure := range bools {
		for _, flag := range options[expectFailure] {
			repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, expectFailure)
		}
	}
}

func TestRepositoryDepGraph(t *testing.T) {
	name := "repository should have github advanced security disabled"
	testedPolicyName := "ghas_dependency_review_not_enabled"
	makeMockData := func(count int) githubcollected.Repository {
		return makeRepo(githubcollected.GitHubQLRepository{
			Name: "REPO",
			DependencyGraphManifests: &githubcollected.GitHubQLDependencyGraphManifests{
				TotalCount: count,
			},
		})
	}

	counts := []int{
		0,
		4,
	}
	for i, flag := range bools {
		repositoryTestTemplate(t, name, makeMockData(counts[i]), testedPolicyName, flag)
	}
}

func TestRepositoryActionsSettingsDefaultTokenPermissions(t *testing.T) {
	name := "repository actions settings is set to read-write"
	testedPolicyName := "token_default_permissions_is_read_write"
	makeMockData := func(flag string) githubcollected.Repository {
		return githubcollected.Repository{
			ActionsTokenPermissions: &types.TokenPermissions{
				DefaultWorkflowPermissions: &flag,
			},
		}
	}

	options := map[bool]string{
		false: "read",
		true:  "write",
	}

	for _, expectFailure := range bools {
		flag := options[expectFailure]
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, expectFailure)
	}
}

func TestRepositoryActionsSettingsActionsCanApprovePullRequests(t *testing.T) {
	name := "repository actions can approve pull requests"
	testedPolicyName := "actions_can_approve_pull_requests"
	makeMockData := func(flag bool) githubcollected.Repository {
		return githubcollected.Repository{
			ActionsTokenPermissions: &types.TokenPermissions{
				CanApprovePullRequestReviews: &flag,
			},
		}
	}

	options := map[bool]bool{
		false: false,
		true:  true,
	}

	for _, expectFailure := range bools {
		flag := options[expectFailure]
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, expectFailure)
	}
}