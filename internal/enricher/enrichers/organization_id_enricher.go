package enrichers

import (
	"context"
	"strconv"

	"github.com/Legit-Labs/legitify/internal/analyzers"
	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
)

const OrganizationId = "organizationId"

func NewOrganizationIdEnricher(ctx context.Context) Enricher {
	return &organizationIdEnricher{}
}

type organizationIdEnricher struct {
}

func (e *organizationIdEnricher) Enrich(data analyzers.AnalyzedData) (Enrichment, bool) {
	switch t := data.Entity.(type) {
	case githubcollected.OrganizationActions:
		return NewBasicEnrichment(strconv.FormatInt(*t.Organization.ID, 10), OrganizationId), true
	}
	return nil, false
}

func (e *organizationIdEnricher) Name() string {
	return OrganizationId
}
