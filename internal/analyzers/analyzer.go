package analyzers

import (
	"context"
	"github.com/Legit-Labs/legitify/internal/analyzers/parsing_utils"
	"github.com/Legit-Labs/legitify/internal/analyzers/skippers"
	githubcollected "github.com/Legit-Labs/legitify/internal/collected"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"log"

	"github.com/Legit-Labs/legitify/internal/collectors"
	"github.com/Legit-Labs/legitify/internal/common/group_waiter"
	"github.com/Legit-Labs/legitify/internal/common/severity"
	"github.com/Legit-Labs/legitify/internal/opa/opa_engine"
	"github.com/open-policy-agent/opa/ast"
)

type PolicyStatus = string

const (
	PolicyPassed  PolicyStatus = "PASSED"
	PolicyFailed  PolicyStatus = "FAILED"
	PolicySkipped PolicyStatus = "SKIPPED"
)

type AnalyzedData struct {
	Entity                   githubcollected.Entity
	Namespace                namespace.Namespace
	PolicyName               string
	FullyQualifiedPolicyName string
	Title                    string
	Description              string
	Annotations              *ast.Annotations
	RequiredEnrichers        []string
	RemediationSteps         []string
	Severity                 severity.Severity
	CanonicalLink            string
	ExtraData                interface{}
	Status                   PolicyStatus
}

type Analyzer interface {
	Analyze(dataChannel <-chan collectors.CollectedData) <-chan AnalyzedData
}

func NewAnalyzer(ctx context.Context, enginer opa_engine.Enginer, skipper skippers.Skipper) Analyzer {
	return &analyzer{
		context: ctx,
		engine:  enginer,
		skipper: skipper,
	}
}

type analyzer struct {
	context context.Context
	engine  opa_engine.Enginer
	skipper skippers.Skipper
}

func newAnalyzedData(collectedData collectors.CollectedData, result opa_engine.QueryResult, status PolicyStatus) AnalyzedData {
	return AnalyzedData{
		Entity:                   collectedData.Entity,
		Namespace:                collectedData.Namespace,
		PolicyName:               result.PolicyName,
		FullyQualifiedPolicyName: result.FullyQualifiedPolicyName,
		Annotations:              result.Annotations,
		Title:                    result.Annotations.Title,
		Description:              result.Annotations.Description,
		RequiredEnrichers:        parsing_utils.ResolveAnnotation(result.Annotations.Custom["requiredEnrichers"]),
		RemediationSteps:         parsing_utils.ResolveAnnotation(result.Annotations.Custom["remediationSteps"]),
		Severity:                 resolveSeverity(result),
		CanonicalLink:            collectedData.Entity.CanonicalLink(),
		ExtraData:                result.ExtraData,
		Status:                   status,
	}
}

func (a *analyzer) Analyze(dataChannel <-chan collectors.CollectedData) <-chan AnalyzedData {
	outputChannel := make(chan AnalyzedData)

	go func() {
		defer close(outputChannel)
		gw := group_waiter.New()
		for data := range dataChannel {
			data := data
			gw.Do(func() {
				results, err := a.engine.Query(a.context, data.Namespace, data.Entity)
				if err != nil {
					log.Printf("Failed to query opa %s: %s", data.Namespace, err)
					return
				}

				for _, result := range results {
					status := a.resolvePolicyStatus(data, result)
					outputChannel <- newAnalyzedData(data, result, status)
				}
			})
		}
		gw.Wait()
	}()

	return outputChannel
}

func (a *analyzer) resolvePolicyStatus(data collectors.CollectedData, opaResult opa_engine.QueryResult) PolicyStatus {
	if a.skipper.ShouldSkip(data, opaResult) {
		return PolicySkipped
	}

	if !opaResult.IsViolation {
		return PolicyPassed
	}

	return PolicyFailed
}

func resolveSeverity(qResult opa_engine.QueryResult) severity.Severity {
	s := severity.Unknown
	raw := qResult.Annotations.Custom["severity"]
	sRaw, ok := raw.(string)

	if !ok {
		log.Printf("Invalid severity type \"%T\" for policy %s\n", sRaw, qResult.FullyQualifiedPolicyName)
	} else if !severity.IsValid(sRaw) {
		log.Printf("Invalid severity value \"%s\" for policy %s\n", sRaw, qResult.FullyQualifiedPolicyName)
	} else {
		s = sRaw
	}

	return s
}
