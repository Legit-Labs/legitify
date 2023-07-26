package scheme_test

import (
	"encoding/json"

	"github.com/Legit-Labs/legitify/internal/collected"

	"github.com/Legit-Labs/legitify/internal/analyzers"
	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/enricher/enrichers"
	"github.com/google/go-github/v53/github"

	"github.com/Legit-Labs/legitify/internal/common/map_utils"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/common/severity"
	"github.com/Legit-Labs/legitify/internal/enricher"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
)

func auxSample() map[string]enrichers.Enrichment {
	aux := map[string]enrichers.Enrichment{
		"A": enrichers.NewBasicEnrichment("foo"),
		"B": enrichers.NewBasicEnrichment("42"),
	}

	return aux
}
func auxSample2() map[string]enrichers.Enrichment {
	aux := map[string]enrichers.Enrichment{
		"BAR":  enrichers.NewBasicEnrichment("xxx"),
		"BLUE": enrichers.NewBasicEnrichment("purple"),
	}

	return aux
}

func FullyQualifiedPolicyNameSample() string {
	return "full/policy1"
}
func FullyQualifiedPolicyNameSample2() string {
	return "full/policy2"
}

const (
	policy_1_name     = "policy1"
	policy_1_title    = "My policy example"
	policy_1_desc     = "This is an example policy that checks for a specific pattern in a file"
	policy_1_severity = severity.Low

	policy_2_name     = "policy2"
	policy_2_title    = "My other policy"
	policy_2_desc     = "This is a different example policy\nthat checks for multiline"
	policy_2_severity = severity.High
)

func arbitraryEntity() collected.Entity {
	var entityID int64 = 666
	var entityName = "arbitrary"
	var link = "link1"
	var login = "login"

	e := githubcollected.Organization{
		Organization: &githubcollected.ExtendedOrg{
			Organization: github.Organization{
				ID:      &entityID,
				Name:    &entityName,
				Login:   &login,
				HTMLURL: &link},
		},
	}

	return e
}

func arbitraryEntity2() collected.Entity {
	var entityID int64 = 667
	var entityName = "arbitrary2"
	var link = "link2"

	e := githubcollected.Repository{
		Repository: &githubcollected.GitHubQLRepository{
			DatabaseId: entityID,
			Name:       entityName,
			Url:        link,
		},
	}
	return e
}

var (
	policy_1_entity        = arbitraryEntity()
	RemediationStepsSample = []string{"do", "that"}

	policy_2_entity         = arbitraryEntity2()
	RemediationStepsSample2 = []string{"dont", "do", "that"}
)

func policyInfoSample() scheme.PolicyInfo {
	return scheme.PolicyInfo{
		PolicyName:               policy_1_name,
		Title:                    policy_1_title,
		Description:              policy_1_desc,
		FullyQualifiedPolicyName: FullyQualifiedPolicyNameSample(),
		RemediationSteps:         RemediationStepsSample,
		Severity:                 policy_1_severity,
		Namespace:                namespace.Organization,
	}
}
func policyInfoSample2() scheme.PolicyInfo {
	return scheme.PolicyInfo{
		PolicyName:               policy_2_name,
		Title:                    policy_2_title,
		Description:              policy_2_desc,
		FullyQualifiedPolicyName: FullyQualifiedPolicyNameSample2(),
		RemediationSteps:         RemediationStepsSample2,
		Severity:                 policy_2_severity,
		Namespace:                namespace.Repository,
	}
}

func enrichedDataSample(isFirst bool, withAux bool) enricher.EnrichedData {
	aux := auxSample()
	if !withAux {
		aux = nil
	}
	link := policy_1_entity.CanonicalLink()
	if isFirst {
		link = first(link)
	} else {
		link = second(link)
	}

	return enricher.EnrichedData{
		Entity:                   policy_1_entity,
		Namespace:                namespace.Organization,
		PolicyName:               policy_1_name,
		FullyQualifiedPolicyName: FullyQualifiedPolicyNameSample(),
		Title:                    policy_1_title,
		Description:              policy_1_desc,
		Enrichers:                aux,
		RemediationSteps:         RemediationStepsSample,
		Severity:                 policy_1_severity,
		CanonicalLink:            link,
		Status:                   analyzers.PolicyFailed,
	}
}

func enrichedDataSample2(isFirst bool, withAux bool) enricher.EnrichedData {
	aux := auxSample2()
	if !withAux {
		aux = nil
	}
	link := policy_2_entity.CanonicalLink()
	if isFirst {
		link = first(link)
	} else {
		link = second(link)
	}

	return enricher.EnrichedData{
		Entity:                   policy_2_entity,
		Namespace:                namespace.Repository,
		PolicyName:               policy_2_name,
		FullyQualifiedPolicyName: FullyQualifiedPolicyNameSample2(),
		Title:                    policy_2_title,
		Description:              policy_2_desc,
		Enrichers:                aux,
		RemediationSteps:         RemediationStepsSample2,
		Severity:                 policy_2_severity,
		CanonicalLink:            link,
		Status:                   analyzers.PolicyFailed,
	}
}

func first(x string) string {
	return x + "a"
}
func second(x string) string {
	return x + "b"
}

func SchemeSample() *scheme.Flattened {
	sample := scheme.NewFlattenedScheme()

	sample.AsOrderedMap().Set(FullyQualifiedPolicyNameSample(), scheme.OutputData{
		PolicyInfo: policyInfoSample(),
		Violations: []scheme.Violation{
			{
				ViolationEntityType: policy_1_entity.ViolationEntityType(),
				CanonicalLink:       first(policy_1_entity.CanonicalLink()),
				Aux:                 map_utils.ToKeySortedMap(auxSample()),
				Status:              analyzers.PolicyFailed,
			},
			{
				ViolationEntityType: policy_1_entity.ViolationEntityType(),
				CanonicalLink:       second(policy_1_entity.CanonicalLink()),
				Aux:                 nil,
				Status:              analyzers.PolicyFailed,
			},
		},
	})

	sample.AsOrderedMap().Set(FullyQualifiedPolicyNameSample2(), scheme.OutputData{
		PolicyInfo: policyInfoSample2(),
		Violations: []scheme.Violation{
			{
				ViolationEntityType: policy_2_entity.ViolationEntityType(),
				CanonicalLink:       first(policy_2_entity.CanonicalLink()),
				Aux:                 map_utils.ToKeySortedMap(auxSample2()),
				Status:              analyzers.PolicyFailed,
			},
			{
				ViolationEntityType: policy_2_entity.ViolationEntityType(),
				CanonicalLink:       second(policy_2_entity.CanonicalLink()),
				Aux:                 map_utils.ToKeySortedMap(auxSample2()),
				Status:              analyzers.PolicyFailed,
			},
		},
	})

	return sample
}

// corresponding to SchemeSample
func EnrichedDataSample() []enricher.EnrichedData {
	// Note: deliberarly mixing the order of different policies
	return []enricher.EnrichedData{
		enrichedDataSample(true, true),
		enrichedDataSample(false, false),
		enrichedDataSample2(true, true),
		enrichedDataSample2(false, true),
	}
}

// convert struct to map recursively (for simple nested comparison)
func StructToMap(obj interface{}) (map[string]interface{}, error) {
	var newMap map[string]interface{}
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &newMap)
	if err != nil {
		return nil, err
	}
	return newMap, nil
}

func CombineSchemes(a, b *scheme.Flattened) *scheme.Flattened {
	for _, policyName := range b.AsOrderedMap().Keys() {
		outputData := b.GetPolicyData(policyName)
		if _, ok := a.AsOrderedMap().Get(policyName); !ok {
			a.AsOrderedMap().Set(policyName, scheme.NewOutputData(outputData.PolicyInfo))
		}
		violation := a.GetPolicyData(policyName)

		a.AsOrderedMap().Set(policyName, scheme.AppendViolations(violation, outputData.Violations...))
	}
	return a
}
