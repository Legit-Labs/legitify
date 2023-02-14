package scheme

import (
	"fmt"

	"github.com/Legit-Labs/legitify/internal/analyzers"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/common/severity"
	"github.com/Legit-Labs/legitify/internal/common/utils"
	"github.com/Legit-Labs/legitify/internal/enricher"
	"github.com/iancoleman/orderedmap"
)

type PolicyInfo struct {
	Title                    string              `json:"title"`
	Description              string              `json:"description"`
	PolicyName               string              `json:"policyName"`
	FullyQualifiedPolicyName string              `json:"fullyQualifiedPolicyName"`
	Severity                 severity.Severity   `json:"severity"`
	Threat                   []string            `json:"threat"`
	RemediationSteps         []string            `json:"remediationSteps"`
	Namespace                namespace.Namespace `json:"namespace"`
}

func NewPolicyInfoFromMap(m map[string]interface{}) (*PolicyInfo, error) {
	var p PolicyInfo
	err := utils.ShalloUnmarshalMap(m, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

type Violation struct { // Must be exported for json marshal
	ViolationEntityType string                 `json:"violationEntityType"`
	CanonicalLink       string                 `json:"canonicalLink"`
	Aux                 *orderedmap.OrderedMap `json:"aux"`
	Status              analyzers.PolicyStatus `json:"status"`
}

func newAuxFromMap(m *orderedmap.OrderedMap) (*orderedmap.OrderedMap, error) {
	newM := orderedmap.New()
	for _, name := range m.Keys() {
		v := utils.UnsafeGetUntyped(m, name)
		if v == nil {
			newM.Set(name, nil)
			continue
		}
		enrichment, err := enricher.NewEnricherManager().Parse(name, v)
		if err != nil {
			return nil, fmt.Errorf("failed to enrich %v: %v", name, err)
		}
		newM.Set(name, enrichment)
	}

	return newM, nil
}

func NewViolationFromMap(m *orderedmap.OrderedMap) (*Violation, error) {
	var p Violation
	err := utils.ShalloUnmarshalOrderedMap(m, &p)
	if err != nil {
		return nil, err
	}

	p.Aux, err = newAuxFromMap(p.Aux)
	if err != nil {
		return nil, fmt.Errorf("failed to parse aux for violation: %v", err)
	}

	return &p, nil
}

type OutputData struct { // Must be exported for json marshal
	PolicyInfo PolicyInfo  `json:"policyInfo"`
	Violations []Violation `json:"violations"`
}

func NewOutputData(policyInfo PolicyInfo) OutputData {
	return OutputData{
		PolicyInfo: policyInfo,
		Violations: []Violation{},
	}
}

func NewOutputDataFromMap(m *orderedmap.OrderedMap) (*OutputData, error) {
	_, okP := m.Get("policyInfo")
	_, okV := m.Get("violations")
	if !okP || !okV {
		return nil, fmt.Errorf("output data missing fields")
	}

	infoMap := utils.UnsafeGet[orderedmap.OrderedMap](m, "policyInfo")
	var policyInfo PolicyInfo
	err := utils.ShalloUnmarshalOrderedMap(&infoMap, &policyInfo)
	if err != nil {
		return nil, err
	}

	outputData := NewOutputData(policyInfo)
	violationMaps := utils.UnsafeGet[[]interface{}](m, "violations")
	for _, v := range violationMaps {
		asMap := v.(orderedmap.OrderedMap)
		violation, err := NewViolationFromMap(&asMap)
		if err != nil {
			return nil, err
		}
		outputData.Violations = append(outputData.Violations, *violation)
	}

	return &outputData, nil
}

func (o OutputData) Clone() OutputData {
	clone := NewOutputData(o.PolicyInfo)
	clone = AppendViolations(clone, o.Violations...)
	return clone
}

func AppendViolations(o OutputData, violations ...Violation) OutputData {
	o.Violations = append(o.Violations, violations...)
	return o
}
