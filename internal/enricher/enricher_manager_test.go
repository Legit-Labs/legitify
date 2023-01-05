package enricher_test

import (
	"context"
	"github.com/Legit-Labs/legitify/internal/collected"
	"testing"

	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/enricher"
	"github.com/google/go-github/v49/github"

	"github.com/Legit-Labs/legitify/internal/analyzers"
	"github.com/stretchr/testify/require"
)

type EnricherTestRequires struct {
	ctx context.Context
	e   enricher.EnricherManager
}

func arrangeEnricher(t *testing.T) EnricherTestRequires {
	ctx := context.Background()

	enricher := enricher.NewEnricherManager(ctx)
	require.NotNilf(t, enricher, "failed to create enricher")

	return EnricherTestRequires{
		ctx: ctx,
		e:   enricher,
	}
}

func arbitraryEntity() collected.Entity {
	var entityID int64 = 666
	var entityName string = "arbitrary"

	e := githubcollected.Organization{
		Organization: &githubcollected.ExtendedOrg{
			Organization: github.Organization{
				ID:      &entityID,
				Name:    &entityName,
				Login:   &entityName,
				HTMLURL: &entityName},
		},
	}

	return e
}

func TestEnricher_PolicyWithNoEnricher_DoesNotEnrich(t *testing.T) {
	enricherData := arrangeEnricher(t)
	data := make(chan analyzers.AnalyzedData, 3)
	outputChannel := enricherData.e.Enrich(data)

	go func() {

		someAnalyzedData := analyzers.AnalyzedData{
			Entity:                   arbitraryEntity(),
			PolicyName:               "A Policy",
			FullyQualifiedPolicyName: "A Full Policy",
			Annotations:              nil,
			RequiredEnrichers:        nil,
		}

		data <- someAnalyzedData
		data <- someAnalyzedData
		data <- someAnalyzedData

		close(data)
	}()
	for outgoingMessage := range outputChannel {
		require.Equalf(t, len(outgoingMessage.Enrichers), len(enricher.DefaultEnrichers), "A policy without enrichers should not enrich data")
	}
}

func TestEnricher_PolicyWithEnricher_EnrichData(t *testing.T) {
	enricherData := arrangeEnricher(t)
	data := make(chan analyzers.AnalyzedData, 3)

	entity := githubcollected.Repository{
		Repository: &githubcollected.GitHubQLRepository{
			Name:               "A Name",
			RebaseMergeAllowed: false,
			Url:                "",
			DatabaseId:         0,
			ForkingAllowed:     false,
			DefaultBranchRef:   nil,
		},
		VulnerabilityAlertsEnabled: nil,
	}
	outputChannel := enricherData.e.Enrich(data)
	go func() {
		someAnalyzedData := analyzers.AnalyzedData{
			Entity:                   entity,
			PolicyName:               "A Policy",
			FullyQualifiedPolicyName: "A Full Policy",
			Annotations:              nil,
		}

		data <- someAnalyzedData
		data <- someAnalyzedData
		data <- someAnalyzedData

		close(data)
	}()
	for outgoingMessage := range outputChannel {
		require.Equalf(t, len(outgoingMessage.Enrichers), 2, "A policy with no enrichers should enrich data twice (default enrichers)")
	}
}
