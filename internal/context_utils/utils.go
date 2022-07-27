package context_utils

import (
	"context"
	"github.com/Legit-Labs/legitify/internal/common/types"

	"github.com/Legit-Labs/legitify/internal/common/permissions"
)

type contextKey string

const (
	organizationKey     contextKey = "org"
	repositoryKey       contextKey = "repo"
	tokenScopesKey      contextKey = "tokenScopes"
	scorecardEnabledKey contextKey = "scorecardEnabled"
	scorecardVerboseKey contextKey = "scorecardVerbose"
)

func NewContextWithRepos(repos []types.RepositoryWithOwner) context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, repositoryKey, repos)
}

func NewContextWithOrg(org []string) context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, organizationKey, org)
}

func NewContextWithScorecard(ctx context.Context, scorecardEnabled bool, scorecardVerbose bool) context.Context {
	c := context.WithValue(ctx, scorecardEnabledKey, scorecardEnabled)
	return context.WithValue(c, scorecardVerboseKey, scorecardVerbose)
}
func NewContextWithTokenScopes(ctx context.Context, tokenScopes permissions.TokenScopes) context.Context {
	return context.WithValue(ctx, tokenScopesKey, tokenScopes)
}

func GetTokenScopes(ctx context.Context) permissions.TokenScopes {
	return ctx.Value(tokenScopesKey).(permissions.TokenScopes)
}

func GetScorecardEnabled(ctx context.Context) bool {
	val, ok := ctx.Value(scorecardEnabledKey).(bool)
	return ok && val
}

func GetScorecardVerbose(ctx context.Context) bool {
	val, ok := ctx.Value(scorecardVerboseKey).(bool)
	return ok && val
}

func GetRepositories(ctx context.Context) ([]types.RepositoryWithOwner, bool) {
	val, ok := ctx.Value(repositoryKey).([]types.RepositoryWithOwner)
	return val, ok
}
