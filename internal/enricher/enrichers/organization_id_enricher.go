package enrichers

import (
	"context"
	"strconv"

	"github.com/Legit-Labs/legitify/internal/analyzers"
	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
)

const OrganizationId = "organizationId"

func NewOrganizationIdEnricher(ctx context.Context) Enricher {
	return organizationIdEnricher{
		newBasicEnricher(enrichOrgId),
	}
}

type organizationIdEnricher struct {
	basicEnricher
}

func enrichOrgId(data analyzers.AnalyzedData) (string, bool) {
	switch t := data.Entity.(type) {
	case githubcollected.OrganizationActions:
		return strconv.FormatInt(*t.Organization.ID, 10), true
	}
	return "", false
}
