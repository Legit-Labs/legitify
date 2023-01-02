package scheme

import (
	"sort"

	"github.com/Legit-Labs/legitify/internal/common/namespace"

	"github.com/Legit-Labs/legitify/internal/analyzers"
	"github.com/Legit-Labs/legitify/internal/enricher/enrichers"

	"github.com/Legit-Labs/legitify/internal/common/severity"
	"github.com/Legit-Labs/legitify/internal/common/utils"
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

type Violation struct { // Must be exported for json marshal
	ViolationEntityType string                          `json:"violationEntityType"`
	CanonicalLink       string                          `json:"canonicalLink"`
	Aux                 map[string]enrichers.Enrichment `json:"aux"`
	Status              analyzers.PolicyStatus
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

func (o OutputData) Clone() OutputData {
	clone := NewOutputData(o.PolicyInfo)
	clone = AppendViolations(clone, o.Violations...)
	return clone
}

func AppendViolations(o OutputData, violations ...Violation) OutputData {
	o.Violations = append(o.Violations, violations...)
	return o
}

// FlattenedScheme maps policy-qualified-names to OutputData
// It is used as the default output, from which the converter selects the desired scheme
type FlattenedScheme struct {
	orderedmap.OrderedMap
} // Must be exported for json marshal

func NewFlattenedScheme() FlattenedScheme {
	return FlattenedScheme{*orderedmap.New()}
}

func (s *FlattenedScheme) Clone() FlattenedScheme {
	clone := NewFlattenedScheme()
	for _, k := range s.Keys() {
		v := s.GetPolicyData(k).Clone()
		clone.Set(k, v)
	}
	return clone
}

func (s FlattenedScheme) AsOrderedMap() *orderedmap.OrderedMap {
	return &s.OrderedMap
}

func (s FlattenedScheme) GetPolicyData(policyName string) OutputData {
	return utils.UnsafeGet(s.AsOrderedMap(), policyName).(OutputData)
}

// ByTypeScheme maps Entity Type to the default scheme
type ByTypeScheme = *orderedmap.OrderedMap // Must be exported for json marshal

func NewByTypeScheme() ByTypeScheme {
	return orderedmap.New()
}

// ByResourceScheme maps Resource to the default scheme
type ByResourceScheme = *orderedmap.OrderedMap // Must be exported for json marshal

func NewByResourceScheme() ByResourceScheme {
	return orderedmap.New()
}

// BySeverityScheme maps Resource to the default scheme
type BySeverityScheme = *orderedmap.OrderedMap // Must be exported for json marshal

func NewBySeverityScheme() BySeverityScheme {
	return orderedmap.New()
}

func policiesSortBySeverityLess(i, j *orderedmap.Pair) bool {
	iSev := i.Value().(OutputData).PolicyInfo.Severity
	jSev := j.Value().(OutputData).PolicyInfo.Severity
	if iSev != jSev {
		return severity.Less(iSev, jSev)
	}

	iKey := i.Key()
	jKey := j.Key()
	return iKey < jKey
}

func policiesSortByNamespaceLess(i, j *orderedmap.Pair) bool {
	namespaceOrder := map[namespace.Namespace]int{
		namespace.Organization: 0,
		namespace.Actions:      1,
		namespace.Member:       2,
		namespace.Repository:   3,
		namespace.RunnerGroup:  4,
	}

	iNamespace := i.Value().(OutputData).PolicyInfo.Namespace
	jNamespace := j.Value().(OutputData).PolicyInfo.Namespace

	if iNamespace != jNamespace {
		return namespaceOrder[iNamespace] < namespaceOrder[jNamespace]
	}

	return policiesSortBySeverityLess(i, j)
}

type ViolationFilter func(violation Violation) bool

func FilterPoliciesByViolations(output FlattenedScheme, filter ViolationFilter) FlattenedScheme {
	filteredScheme := NewFlattenedScheme()
	for _, policyName := range output.Keys() {
		outputData := output.GetPolicyData(policyName)
		filteredViolations := []Violation{}
		any := false
		for _, violation := range outputData.Violations {
			if filter(violation) {
				filteredViolations = append(filteredViolations, violation)
				any = true
			}
		}
		if any {
			outputData.Violations = filteredViolations
			filteredScheme.Set(policyName, outputData)
		}
	}

	return filteredScheme
}

func FilterViolationsByStatus(output FlattenedScheme, status analyzers.PolicyStatus) FlattenedScheme {
	filter := func(violation Violation) bool {
		return violation.Status == status
	}
	return FilterPoliciesByViolations(output, filter)
}

func OnlyFailedViolations(output FlattenedScheme) FlattenedScheme {
	return FilterViolationsByStatus(output, analyzers.PolicyFailed)
}

func sortOutputData(outputData OutputData) OutputData {
	less := func(i, j int) bool {
		iLink := outputData.Violations[i].CanonicalLink
		jLink := outputData.Violations[j].CanonicalLink
		return iLink < jLink
	}

	sort.SliceStable(outputData.Violations, less)
	return outputData
}

func SortScheme(output FlattenedScheme, inplace bool, lessFunc func(a *orderedmap.Pair, b *orderedmap.Pair) bool) FlattenedScheme {
	if !inplace {
		output = output.Clone()
	}
	output.Sort(lessFunc)

	for _, policyName := range output.Keys() {
		outputData := output.GetPolicyData(policyName)
		outputData = sortOutputData(outputData)
		output.Set(policyName, outputData)
	}

	return output
}

func SortSchemeBySeverity(output FlattenedScheme, inplace bool) FlattenedScheme {
	return SortScheme(output, inplace, policiesSortBySeverityLess)
}

func SortSchemeByNamespace(output FlattenedScheme, inplace bool) FlattenedScheme {
	return SortScheme(output, inplace, policiesSortByNamespaceLess)
}
