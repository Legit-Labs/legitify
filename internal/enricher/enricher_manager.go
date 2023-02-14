package enricher

import (
	"context"
	"fmt"
	"log"

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
	Enrich(ctx context.Context, analyzedDataChannel <-chan analyzers.AnalyzedData) <-chan EnrichedData
	Parse(name string, data interface{}) (enrichers.Enrichment, error)
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

var mapping = map[string]enrichers.Enricher{
	enrichers.EntityId:       enrichers.NewEntityIdEnricher(),
	enrichers.EntityName:     enrichers.NewEntityNameEnricher(),
	enrichers.OrganizationId: enrichers.NewOrganizationIdEnricher(),
	enrichers.Scorecard:      enrichers.NewScorecardEnricher(),
	enrichers.MembersList:    enrichers.NewMembersListEnricher(),
	enrichers.HooksList:      enrichers.NewHooksListEnricher(),
}

func NewEnricherManager() EnricherManager {
	return &enricherManager{}
}

type enricherManager struct {
}

func (e *enricherManager) Enrich(ctx context.Context, analyzedDataChannel <-chan analyzers.AnalyzedData) <-chan EnrichedData {
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
						enricher, err := e.getEnricher(requiredEnricher)
						if err != nil {
							log.Printf("failed to find enricher: %v", err)
							continue
						}

						enrichment, ok := enricher.Enrich(ctx, analyzedData)
						if !ok {
							continue
						}

						enrichments[requiredEnricher] = enrichment
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
func (e *enricherManager) Parse(name string, data interface{}) (enrichers.Enrichment, error) {
	enricher, err := e.getEnricher(name)
	if err != nil {
		return nil, err
	}

	return enricher.Parse(data)
}

func (e *enricherManager) getEnricher(name string) (enrichers.Enricher, error) {
	if e, ok := mapping[name]; ok {
		return e, nil
	} else {
		return nil, fmt.Errorf("failed to find enricher %s", name)
	}
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
