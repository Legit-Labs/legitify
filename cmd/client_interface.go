package cmd

import (
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	"github.com/Legit-Labs/legitify/internal/common/types"
)

type Client interface {
	IsAnalyzable(repo types.RepositoryWithOwner) (bool, error)
	Scopes() permissions.TokenScopes
	Organizations() ([]types.Organization, error)
}
