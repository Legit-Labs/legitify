package githubcollected

import (
	"github.com/Legit-Labs/legitify/internal/clients/github/types"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/scorecard"
	"github.com/google/go-github/v44/github"
	"github.com/shurcooL/githubv4"
)

type GitHubQLPageInfo struct {
	EndCursor       *githubv4.String
	HasNextPage     bool
	HasPreviousPage bool
	StartCursor     *githubv4.String
}

type GitHubQLDependencyGraphManifests struct {
	TotalCount int `json:"total_count"`
}

type GitHubQLRepositoryCollaboratorsEdge struct {
	Permission *githubv4.String `json:"permission"`
}

type GitHubQLRepositoryCollaborators struct {
	Edges []GitHubQLRepositoryCollaboratorsEdge `json:"edges" graphql:"edges"`
}

type GitHubQLRepository struct {
	Name                     string `json:"name"`
	RebaseMergeAllowed       bool
	Url                      string
	DatabaseId               int64
	IsPrivate                bool                              `json:"is_private"`
	ForkingAllowed           bool                              `json:"allow_forking"`
	IsArchived               bool                              `json:"is_archived"`
	DefaultBranchRef         *GitHubQLBranch                   `json:"default_branch"`
	DependencyGraphManifests *GitHubQLDependencyGraphManifests `json:"dependency_graph_manifests" graphql:"dependencyGraphManifests(first: 1)"`
	PushedAt                 *githubv4.DateTime                `json:"pushed_at"`
	ViewerPermission         string                            `json:"viewerPermission"`
}

type GitHubQLBranchProtectionRule struct {
	AllowsDeletions                *bool `json:"allows_deletions,omitempty"`
	AllowsForcePushes              *bool `json:"allows_force_pushes,omitempty"`
	BlocksCreations                *bool `json:"blocks_creations,omitempty"`
	DismissesStaleReviews          *bool `json:"dismisses_stale_reviews,omitempty"`
	IsAdminEnforced                *bool `json:"is_admin_enforced,omitempty"`
	RequiredApprovingReviewCount   *int  `json:"required_approving_review_count,omitempty"`
	RequiresStatusChecks           *bool `json:"requires_status_checks,omitempty"`
	RequiresStrictStatusChecks     *bool `json:"requires_strict_status_checks,omitempty"`
	RestrictsPushes                *bool `json:"restricts_pushes,omitempty"`
	RequiresCodeOwnerReviews       *bool `json:"requires_code_owner_reviews,omitempty"`
	RequiresLinearHistory          *bool `json:"requires_linear_history,omitempty"`
	RequiresConversationResolution *bool `json:"requires_conversation_resolution,omitempty"`
	RequiresCommitSignatures       *bool `json:"requires_commit_signatures,omitempty"`
	RestrictsReviewDismissals      *bool `json:"restricts_review_dismissals,omitempty"`
}

type GitHubQLBranch struct {
	Name                 *string
	BranchProtectionRule *GitHubQLBranchProtectionRule `json:"branch_protection_rule"`
}

type Repository struct {
	Repository                   *GitHubQLRepository     `json:"repository"`
	VulnerabilityAlertsEnabled   *bool                   `json:"vulnerability_alerts_enabled"`
	NoBranchProtectionPermission bool                    `json:"no_branch_protection_permission"`
	Scorecard                    *scorecard.Result       `json:"scorecard,omitempty"`
	Hooks                        []*github.Hook          `json:"hooks"`
	Collaborators                []*github.User          `json:"collaborators"`
	ActionsTokenPermissions      *types.TokenPermissions `json:"actions_token_permissions"`
}

func (r Repository) ViolationEntityType() string {
	return namespace.Repository
}

func (r Repository) CanonicalLink() string {
	return r.Repository.Url
}

func (r Repository) Name() string {
	return r.Repository.Name
}

func (r Repository) ID() int64 {
	// Deliberately using the Org; see membersList enricher
	return r.Repository.DatabaseId
}
