package enricher

import (
	"context"
	"fmt"
	"log"
	"sync"

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

var enricherTextToEnricher struct {
	once    sync.Once
	mapping map[string]enrichers.Enricher
}

func getEnricher(ctx context.Context, name string) (enrichers.Enricher, error) {
	enricherTextToEnricher.once.Do(func() {
		enricherTextToEnricher.mapping = map[string]enrichers.Enricher{
			enrichers.EntityId:       enrichers.NewEntityIdEnricher(ctx),
			enrichers.EntityName:     enrichers.NewEntityNameEnricher(ctx),
			enrichers.OrganizationId: enrichers.NewOrganizationIdEnricher(ctx),
			enrichers.Scorecard:      enrichers.NewScorecardEnricher(ctx),
			enrichers.MembersList:    enrichers.NewMembersListEnricher(ctx),
			enrichers.HooksList:      enrichers.NewHooksListEnricher(ctx),
		}
	})

	if e, ok := enricherTextToEnricher.mapping[name]; ok {
		return e, nil
	} else {
		return nil, fmt.Errorf("failed to find enricher %s", name)
	}
}

func ParseEnrichment(ctx context.Context, name string, data interface{}) (enrichers.Enrichment, error) {
	e, err := getEnricher(ctx, name)
	if err != nil {
		return nil, err
	}

	return e.Parse(data)
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
						enricher, err := getEnricher(e.ctx, requiredEnricher)
						if err != nil {
							log.Printf("failed to find enricher: %v", err)
							continue
						}

						enrichment, ok := enricher.Enrich(analyzedData)
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
