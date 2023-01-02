package enricher

import (
	"context"
	githubcollected "github.com/Legit-Labs/legitify/internal/collected"

	"github.com/Legit-Labs/legitify/internal/analyzers"
	"github.com/Legit-Labs/legitify/internal/common/group_waiter"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/common/severity"
	"github.com/Legit-Labs/legitify/internal/enricher/enrichers"
	"github.com/open-policy-agent/opa/ast"
)

var (
	DefaultEnrichers = []string{
		enrichers.EntityId,
		enrichers.EntityName,
	}
)

type EnricherManager interface {
	Enrich(analyzedDataChannel <-chan analyzers.AnalyzedData) <-chan EnrichedData
}

type EnrichedData struct {
	Entity                   githubcollected.Entity
	Namespace                namespace.Namespace
	PolicyName               string
	FullyQualifiedPolicyName string
	Annotations              *ast.Annotations
	Title                    string
	Description              string
	Enrichers                map[string]enrichers.Enrichment
	Threat                   []string
	RemediationSteps         []string
	Severity                 severity.Severity
	CanonicalLink            string
	Status                   analyzers.PolicyStatus
}

func NewEnricherManager(ctx context.Context) EnricherManager {
	return &enricherManager{
		ctx: ctx,
	}
}

type enricherManager struct {
	ctx context.Context
}

type newEnricherFunc func(ctx context.Context) enrichers.Enricher

var enricherTextToEnricher = map[string]newEnricherFunc{
	enrichers.EntityId:       enrichers.NewEntityIdEnricher,
	enrichers.EntityName:     enrichers.NewEntityNameEnricher,
	enrichers.OrganizationId: enrichers.NewOrganizationIdEnricher,
	enrichers.Scorecard:      enrichers.NewScorecardEnricher,
	enrichers.MembersList:    enrichers.NewMembersListEnricher,
	enrichers.HooksList:      enrichers.NewHooksListEnricher,
}

func newEnrichedData(analyzed analyzers.AnalyzedData, enrichments map[string]enrichers.Enrichment) EnrichedData {
	return EnrichedData{
		Entity:                   analyzed.Entity,
		Namespace:                analyzed.Namespace,
		PolicyName:               analyzed.PolicyName,
		FullyQualifiedPolicyName: analyzed.FullyQualifiedPolicyName,
		Annotations:              analyzed.Annotations,
		Title:                    analyzed.Title,
		Description:              analyzed.Description,
		Enrichers:                enrichments,
		Threat:                   analyzed.Threat,
		Severity:                 analyzed.Severity,
		RemediationSteps:         analyzed.RemediationSteps,
		CanonicalLink:            analyzed.CanonicalLink,
		Status:                   analyzed.Status,
	}
}

func (e *enricherManager) Enrich(analyzedDataChannel <-chan analyzers.AnalyzedData) <-chan EnrichedData {
	outputChannel := make(chan EnrichedData)

	go func() {
		defer close(outputChannel)
		gw := group_waiter.New()
		for analyzedData := range analyzedDataChannel {
			func(analyzedData analyzers.AnalyzedData) {
				gw.Do(func() {
					requiredEnrichers := analyzedData.RequiredEnrichers
					requiredEnrichers = append(requiredEnrichers, DefaultEnrichers...)

					enrichments := make(map[string]enrichers.Enrichment)
					for _, requiredEnricher := range requiredEnrichers {
						createEnricher, ok := enricherTextToEnricher[requiredEnricher]
						if !ok {
							continue
						}

						enricher := createEnricher(e.ctx)

						enrichment, ok := enricher.Enrich(analyzedData)
						if !ok {
							continue
						}

						enrichments[enrichment.Name()] = enrichment
					}
					enrichedData := newEnrichedData(analyzedData, enrichments)

					outputChannel <- enrichedData
				})
			}(analyzedData)
			gw.Wait()
		}
	}()

	return outputChannel
}
