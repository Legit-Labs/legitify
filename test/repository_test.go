package test

import (
	"testing"
	"time"

	"github.com/Legit-Labs/legitify/internal/clients/github/types"
	"github.com/Legit-Labs/legitify/internal/common/scm_type"
	gitlab2 "github.com/xanzy/go-gitlab"

	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	gitlabcollected "github.com/Legit-Labs/legitify/internal/collected/gitlab_collected"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/google/go-github/v53/github"
)

func repositoryTestTemplate(t *testing.T, name string, mockData interface{}, testedPolicyName string, expectFailure bool, scmType scm_type.ScmType) {
	ns := namespace.Repository
	PolicyTestTemplate(t, name, mockData, ns, testedPolicyName, expectFailure, scmType)
}

var bools = []bool{true, false}

func makeRepoWithDeps(repo githubcollected.GitHubQLRepository, deps *githubcollected.GitHubQLDependencyGraphManifests) githubcollected.Repository {
	return githubcollected.Repository{
		Repository:               &repo,
		DependencyGraphManifests: deps,
	}
}
func makeRepo(repo githubcollected.GitHubQLRepository) githubcollected.Repository {
	return makeRepoWithDeps(repo, &githubcollected.GitHubQLDependencyGraphManifests{})
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
		repositoryTestTemplate(t, name, makeMockData(branches[i]), testedPolicyName, flag, scm_type.GitHub)
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
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, flag, scm_type.GitHub)
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
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, flag, scm_type.GitHub)
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
		repositoryTestTemplate(t, name, makeMockData(counts[i]), testedPolicyName, flag, scm_type.GitHub)
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
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, !flag, scm_type.GitHub)
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
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, !flag, scm_type.GitHub)
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
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, !flag, scm_type.GitHub)
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
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, !flag, scm_type.GitHub)
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
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, !flag, scm_type.GitHub)
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
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, !flag, scm_type.GitHub)
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
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, !flag, scm_type.GitHub)
	}
}

func TestRepositoryBranchesUpToDate(t *testing.T) {
	name := "repository should require branches to be up to date before merging"
	testedPolicyName := "requires_branches_up_to_date_before_merge"
	makeMockData := func(flag bool) githubcollected.Repository {
		return makeRepoForBranchProtection(githubcollected.GitHubQLBranchProtectionRule{
			RequiresStatusChecks:       github.Bool(flag),
			RequiresStrictStatusChecks: github.Bool(flag),
		})
	}
	for _, flag := range bools {
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, !flag, scm_type.GitHub)
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
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, !flag, scm_type.GitHub)
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
			repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, expectFailure, scm_type.GitHub)
		}
	}
}

func TestRepositoryDepGraph(t *testing.T) {
	name := "repository should have github advanced security disabled"
	testedPolicyName := "ghas_dependency_review_not_enabled"
	makeMockData := func(count int) githubcollected.Repository {
		return makeRepoWithDeps(githubcollected.GitHubQLRepository{Name: "REPO"},
			&githubcollected.GitHubQLDependencyGraphManifests{TotalCount: count},
		)
	}

	counts := []int{
		0,
		4,
	}
	for i, flag := range bools {
		repositoryTestTemplate(t, name, makeMockData(counts[i]), testedPolicyName, flag, scm_type.GitHub)
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
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, expectFailure, scm_type.GitHub)
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
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, expectFailure, scm_type.GitHub)
	}
}

func TestRepositoryWithNoStaleSecrets(t *testing.T) {
	name := "repository has no secrets"
	testedPolicyName := "repository_secret_is_stale"
	makeMockData := func() githubcollected.Repository {
		return githubcollected.Repository{
			RepoSecrets: []*githubcollected.RepositorySecret{
				{
					Name:      "test1",
					UpdatedAt: int(time.Now().UnixNano()) - 2628000000000000, // one month
				},
				{
					Name:      "test2",
					UpdatedAt: int(time.Now().UnixNano()) - (3 * 2628000000000000), // three months
				},
				{
					Name:      "test3",
					UpdatedAt: int(time.Now().UnixNano()) - (6 * 2628000000000000), // six month
				},
			},
		}
	}
	expectFailure := false
	repositoryTestTemplate(t, name, makeMockData(), testedPolicyName, expectFailure, scm_type.GitHub)
}

func TestRepositoryWithStaleSecrets(t *testing.T) {
	name := "repository has no secrets"
	testedPolicyName := "repository_secret_is_stale"
	makeMockData := func() githubcollected.Repository {
		return githubcollected.Repository{
			RepoSecrets: []*githubcollected.RepositorySecret{
				{
					Name:      "test1",
					UpdatedAt: 1652020546000000000, //08.05.2022
				},
				{
					Name:      "test2",
					UpdatedAt: 957796546000000000, //08.05.2000
				},
				{
					Name:      "test3",
					UpdatedAt: int(time.Now().UnixNano()),
				},
			},
		}
	}
	expectFailure := true
	repositoryTestTemplate(t, name, makeMockData(), testedPolicyName, expectFailure, scm_type.GitHub)
}

func TestRepositoryWithNoSecrets(t *testing.T) {
	name := "repository has no secrets"
	testedPolicyName := "repository_secret_is_stale"
	makeMockData := func() githubcollected.Repository {
		return githubcollected.Repository{
			RepoSecrets: nil,
		}
	}
	expectFailure := false
	repositoryTestTemplate(t, name, makeMockData(), testedPolicyName, expectFailure, scm_type.GitHub)
}

func TestRepositorySecretScanning(t *testing.T) {
	name := "repository secret scanning is disabled"
	testedPolicyName := "secret_scanning_not_enabled"
	makeMockData := func(flag string) githubcollected.Repository {
		return githubcollected.Repository{
			SecurityAndAnalysis: &github.SecurityAndAnalysis{
				SecretScanning: &github.SecretScanning{Status: &flag},
			},
		}
	}

	options := map[bool]string{
		false: "enabled",
		true:  "disabled",
	}

	for _, expectFailure := range bools {
		flag := options[expectFailure]
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, expectFailure, scm_type.GitHub)
	}
}

func TestGitlabRepositoryTooManyAdmins(t *testing.T) {
	name := "Project Has Too Many Owners"
	testedPolicyName := "project_has_too_many_admins"

	makeMockData := func(flag []*gitlab2.ProjectMember) gitlabcollected.Repository {
		return gitlabcollected.Repository{
			Members: flag,
		}
	}

	tmpAdminMember := &gitlab2.ProjectMember{
		AccessLevel: 50,
	}
	tmpRegMember := &gitlab2.ProjectMember{
		AccessLevel: 20,
	}
	trueCase := []*gitlab2.ProjectMember{tmpAdminMember, tmpAdminMember, tmpAdminMember, tmpAdminMember}
	for i := 0; i < 10; i++ {
		trueCase = append(trueCase, tmpRegMember)
	}
	falseCase := []*gitlab2.ProjectMember{tmpAdminMember, tmpAdminMember, tmpAdminMember, tmpAdminMember}
	for i := 0; i < 57; i++ {
		falseCase = append(falseCase, tmpRegMember)
	}

	options := map[bool][]*gitlab2.ProjectMember{
		false: falseCase,
		true:  trueCase,
	}

	for _, expectFailure := range bools {
		flag := options[expectFailure]
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, expectFailure, scm_type.GitLab)
	}
}

func TestGitlabRepositoryAllowForking(t *testing.T) {
	name := "Forking Allowed for This Repository"
	testedPolicyName := "forking_allowed_for_repository"

	makeMockData := func(flag *gitlab2.Project) gitlabcollected.Repository {
		return gitlabcollected.Repository{Project: flag}
	}

	falseCase := &gitlab2.Project{Public: true, ForkingAccessLevel: "disabled"}
	trueCase := &gitlab2.Project{Public: false, ForkingAccessLevel: "enabled"}
	options := map[bool]*gitlab2.Project{
		false: falseCase,
		true:  trueCase,
	}

	for _, expectFailure := range bools {
		flag := options[expectFailure]
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, expectFailure, scm_type.GitLab)
	}
}

func TestGitlabRepositoryNotMaintained(t *testing.T) {
	name := "Repository not maintained"
	testedPolicyName := "project_not_maintained"

	makeMockData := func(flag *gitlab2.Project) gitlabcollected.Repository {
		return gitlabcollected.Repository{Project: flag}
	}
	nowTime := time.Now()
	// Creating a mock for a project, last active more than 10 years ago.
	archivedFewYearsTime := nowTime.AddDate(-10, 0, 0)
	// Creating a mock for a project, last active more than 3 month ago.
	archived5MonthTime := nowTime.AddDate(0, -5, 0)
	falseCase := []*gitlab2.Project{{Archived: false, LastActivityAt: &nowTime}}
	trueCase := []*gitlab2.Project{
		{Public: false, LastActivityAt: &archivedFewYearsTime},
		{Public: false, LastActivityAt: &archived5MonthTime},
	}
	options := map[bool][]*gitlab2.Project{
		false: falseCase,
		true:  trueCase,
	}

	for _, expectFailure := range bools {
		for _, testCase := range options[expectFailure] {
			repositoryTestTemplate(t, name, makeMockData(testCase), testedPolicyName, expectFailure, scm_type.GitLab)
		}
	}
}

func TestGitlabRepositoryMissingBranchProtection(t *testing.T) {
	name := "Default Branch Is Not Protected"
	testedPolicyName := "missing_default_branch_protection"

	makeMockData := func(flag []*gitlab2.ProtectedBranch) gitlabcollected.Repository {
		return gitlabcollected.Repository{Project: &gitlab2.Project{DefaultBranch: "default_branch_name"}, ProtectedBranches: flag}
	}

	defaultBranchProtectedMock := &gitlab2.ProtectedBranch{Name: "default_branch_name"}
	nonDefaultBranchProtectedMock := &gitlab2.ProtectedBranch{Name: "fooBar"}
	falseCase := []*gitlab2.ProtectedBranch{defaultBranchProtectedMock}
	trueCase := []*gitlab2.ProtectedBranch{nonDefaultBranchProtectedMock}
	options := map[bool][]*gitlab2.ProtectedBranch{
		false: falseCase,
		true:  trueCase,
	}

	for _, expectFailure := range bools {
		flag := options[expectFailure]
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, expectFailure, scm_type.GitLab)
	}
}

func TestGitlabRepositoryMissingForcePushProtection(t *testing.T) {
	name := "Default Branch Allows Force Pushes"
	testedPolicyName := "missing_default_branch_protection_force_push"

	makeMockData := func(flag []*gitlab2.ProtectedBranch) gitlabcollected.Repository {
		return gitlabcollected.Repository{Project: &gitlab2.Project{DefaultBranch: "default_branch_name"}, ProtectedBranches: flag}
	}

	defaultBranchProtectedMock := &gitlab2.ProtectedBranch{Name: "default_branch_name", AllowForcePush: false}
	DefaultNotBranchProtectedMock := &gitlab2.ProtectedBranch{Name: "sss", AllowForcePush: false}
	falseCase := []*gitlab2.ProtectedBranch{defaultBranchProtectedMock}
	trueCase := []*gitlab2.ProtectedBranch{DefaultNotBranchProtectedMock}
	options := map[bool][]*gitlab2.ProtectedBranch{
		false: falseCase,
		true:  trueCase,
	}

	for _, expectFailure := range bools {
		flag := options[expectFailure]
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, expectFailure, scm_type.GitLab)
	}
}

func TestGitlabWebhookSSL(t *testing.T) {
	name := "Webhook Configured Without SSL Verification"
	testedPolicyName := "project_webhook_doesnt_require_ssl"

	makeMockData := func(flag []*gitlab2.ProjectHook) gitlabcollected.Repository {
		return gitlabcollected.Repository{Webhooks: flag}
	}

	sslNotVerifiedHookMock := &gitlab2.ProjectHook{EnableSSLVerification: false}
	sslVerifiedHookMock := &gitlab2.ProjectHook{EnableSSLVerification: true}
	falseCase := []*gitlab2.ProjectHook{sslVerifiedHookMock}
	trueCase := []*gitlab2.ProjectHook{sslNotVerifiedHookMock}
	options := map[bool][]*gitlab2.ProjectHook{
		false: falseCase,
		true:  trueCase,
	}

	for _, expectFailure := range bools {
		flag := options[expectFailure]
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, expectFailure, scm_type.GitLab)
	}
}

func TestGitlabPipelineStatusCheck(t *testing.T) {
	name := "Project Does’nt Require All Pipelines to Succeed"
	testedPolicyName := "requires_status_checks"

	makeMockData := func(flag bool) gitlabcollected.Repository {
		return gitlabcollected.Repository{Project: &gitlab2.Project{OnlyAllowMergeIfPipelineSucceeds: flag}}
	}
	options := map[bool]bool{
		false: true,
		true:  false,
	}
	for _, expectFailure := range bools {
		flag := options[expectFailure]
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, expectFailure, scm_type.GitLab)
	}
}

func TestGitlabResolvedThreads(t *testing.T) {
	name := "Project Doesn't Require All Conversations To Be Resolved Before Merge"
	testedPolicyName := "no_conversation_resolution"

	makeMockData := func(flag bool) gitlabcollected.Repository {
		return gitlabcollected.Repository{Project: &gitlab2.Project{OnlyAllowMergeIfAllDiscussionsAreResolved: flag}}
	}
	options := map[bool]bool{
		false: true,
		true:  false,
	}
	for _, expectFailure := range bools {
		flag := options[expectFailure]
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, expectFailure, scm_type.GitLab)
	}
}

func TestGitlabRepositoryMissingSignedCommitsVerifications(t *testing.T) {
	name := "Unsigned Commits Are Not Allowed"
	testedPolicyName := "no_signed_commits"

	makeMockData := func(flag *gitlab2.ProjectPushRules) gitlabcollected.Repository {
		return gitlabcollected.Repository{PushRules: flag}
	}

	falseCase := []*gitlab2.ProjectPushRules{
		{RejectUnsignedCommits: true},
	}
	trueCase := []*gitlab2.ProjectPushRules{
		{RejectUnsignedCommits: false},
		nil,
	}
	options := map[bool][]*gitlab2.ProjectPushRules{
		false: falseCase,
		true:  trueCase,
	}

	for _, expectFailure := range bools {
		for _, testCase := range options[expectFailure] {
			repositoryTestTemplate(t, name, makeMockData(testCase), testedPolicyName, expectFailure, scm_type.GitLab)
		}
	}
}

func TestGitlabRepositoryRequiredReview(t *testing.T) {
	name := "Project Doesn't Require Code Review"
	testedPolicyName := "code_review_not_required"

	makeMockData := func(flag int) gitlabcollected.Repository {
		return gitlabcollected.Repository{MinimumRequiredApprovals: flag}
	}

	falseCase := []int{1, 2}
	trueCase := []int{0}
	options := map[bool][]int{
		false: falseCase,
		true:  trueCase,
	}

	for _, expectFailure := range bools {
		for _, testCase := range options[expectFailure] {
			repositoryTestTemplate(t, name, makeMockData(testCase), testedPolicyName, expectFailure, scm_type.GitLab)
		}
	}
}

func TestGitlabRepositoryMinimum2RequiredReview(t *testing.T) {
	name := "Project Doesn't Require Code Review By At Least Two Reviewers"
	testedPolicyName := "code_review_by_two_members_not_required"

	makeMockData := func(flag int) gitlabcollected.Repository {
		return gitlabcollected.Repository{MinimumRequiredApprovals: flag}
	}

	falseCase := []int{2, 3}
	trueCase := []int{0, 1}
	options := map[bool][]int{
		false: falseCase,
		true:  trueCase,
	}

	for _, expectFailure := range bools {
		for _, testCase := range options[expectFailure] {
			repositoryTestTemplate(t, name, makeMockData(testCase), testedPolicyName, expectFailure, scm_type.GitLab)
		}
	}
}

func TestGitlabRepositoryAllowsReviewRequesterToApproveTheirOwnRequest(t *testing.T) {
	name := "Repository Allows Review Requester To Approve Their Own Request"
	testedPolicyName := "repository_allows_review_requester_to_approve_their_own_request"

	makeMockData := func(flag bool) gitlabcollected.Repository {
		return gitlabcollected.Repository{ApprovalConfiguration: &gitlab2.ProjectApprovals{MergeRequestsAuthorApproval: flag}}
	}

	options := map[bool]bool{
		false: false,
		true:  true,
	}

	for _, expectFailure := range bools {
		flag := options[expectFailure]
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, expectFailure, scm_type.GitLab)
	}
}

func TestGitlabRepositoryAllowsOverridingApproversPolicy(t *testing.T) {
	name := "Merge request authors may override the approvers list"
	testedPolicyName := "repository_allows_overriding_approvers"

	makeMockData := func(flag bool) gitlabcollected.Repository {
		return gitlabcollected.Repository{ApprovalConfiguration: &gitlab2.ProjectApprovals{DisableOverridingApproversPerMergeRequest: flag}}
	}

	options := map[bool]bool{
		false: true,
		true:  false,
	}

	for _, expectFailure := range bools {
		flag := options[expectFailure]
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, expectFailure, scm_type.GitLab)
	}
}

func TestGitlabRepositoryAllowsCommitterApprovalsPolicy(t *testing.T) {
	name := "Repository Allows Committer Approvals Policy"
	testedPolicyName := "repository_allows_committer_approvals_policy"

	makeMockData := func(flag bool) gitlabcollected.Repository {
		return gitlabcollected.Repository{ApprovalConfiguration: &gitlab2.ProjectApprovals{MergeRequestsDisableCommittersApproval: flag}}
	}

	options := map[bool]bool{
		false: true,
		true:  false,
	}

	for _, expectFailure := range bools {
		flag := options[expectFailure]
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, expectFailure, scm_type.GitLab)
	}
}

func TestGitlabRepositoryMissingRequiredCodeOwnersReview(t *testing.T) {
	name := "Code review is not limited to code-owners only"
	testedPolicyName := "repository_require_code_owner_reviews_policy"

	makeMockData := func(flag []*gitlab2.ProtectedBranch) gitlabcollected.Repository {
		return gitlabcollected.Repository{Project: &gitlab2.Project{DefaultBranch: "default_branch_name"}, ProtectedBranches: flag}
	}

	defaultBranchProtectedMock := &gitlab2.ProtectedBranch{Name: "default_branch_name", CodeOwnerApprovalRequired: true}
	defaultNotBranchProtectedMock := &gitlab2.ProtectedBranch{Name: "default_branch_name", CodeOwnerApprovalRequired: false}
	falseCase := []*gitlab2.ProtectedBranch{defaultBranchProtectedMock}
	trueCase := []*gitlab2.ProtectedBranch{defaultNotBranchProtectedMock}
	options := map[bool][]*gitlab2.ProtectedBranch{
		false: falseCase,
		true:  trueCase,
	}

	for _, expectFailure := range bools {
		flag := options[expectFailure]
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, expectFailure, scm_type.GitLab)
	}
}

func TestGitlabRepositoryDismissStaleReviews(t *testing.T) {
	name := "Repository Dismiss Stale Reviews"
	testedPolicyName := "repository_dismiss_stale_reviews"

	makeMockData := func(flag bool) gitlabcollected.Repository {
		return gitlabcollected.Repository{ApprovalConfiguration: &gitlab2.ProjectApprovals{ResetApprovalsOnPush: flag}}
	}

	options := map[bool]bool{
		false: true,
		true:  false,
	}

	for _, expectFailure := range bools {
		flag := options[expectFailure]
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, expectFailure, scm_type.GitLab)
	}
}

func TestGitlabRepositoryRestrictsOverrideVariables(t *testing.T) {
	name := "Restrict Override Of Defined Variables"
	testedPolicyName := "overriding_defined_variables_isnt_restricted"

	makeMockData := func(flag bool) gitlabcollected.Repository {
		return gitlabcollected.Repository{Project: &gitlab2.Project{RestrictUserDefinedVariables: flag}}
	}
	options := map[bool]bool{
		false: true,
		true:  false,
	}
	for _, expectFailure := range bools {
		flag := options[expectFailure]
		repositoryTestTemplate(t, name, makeMockData(flag), testedPolicyName, expectFailure, scm_type.GitLab)
	}
}
