package cmd

import (
	"context"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	"github.com/Legit-Labs/legitify/internal/common/types"
)

type Client interface {
	IsAnalyzable(ctx context.Context, repo types.RepositoryWithOwner) (bool, error)
	Scopes() permissions.TokenScopes
}
