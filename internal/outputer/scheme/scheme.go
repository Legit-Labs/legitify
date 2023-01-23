package scheme

import (
	"fmt"
	"sort"

	"github.com/iancoleman/orderedmap"
)

type Scheme interface {
	AsOrderedMap() *orderedmap.OrderedMap
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

type SchemeType = string

const (
	TypeFlattened        SchemeType = "flattened"
	TypeGroupByNamespace SchemeType = "group-by-namespace"
	TypeGroupByResource  SchemeType = "group-by-resource"
	TypeGroupBySeverity  SchemeType = "group-by-severity"

	DefaultScheme = TypeFlattened
)

func SchemeTypes() []SchemeType {
	return []SchemeType{
		TypeFlattened,
		TypeGroupByNamespace,
		TypeGroupByResource,
		TypeGroupBySeverity,
	}
}

func DetectSchemeType(s interface{}) (SchemeType, error) {
	switch t := s.(type) {
	case *Flattened:
		return TypeFlattened, nil
	case *ByNamespace:
		return TypeGroupByNamespace, nil
	case *ByResource:
		return TypeGroupByResource, nil
	case *BySeverity:
		return TypeGroupBySeverity, nil
	default:
		return DefaultScheme, fmt.Errorf("invalid scheme type: %T", t)
	}
}
