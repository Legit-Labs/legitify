package enrichers

import (
	"fmt"

	"github.com/Legit-Labs/legitify/internal/common/utils"
	"github.com/iancoleman/orderedmap"
)

type GenericListEnrichment []orderedmap.OrderedMap

func (se GenericListEnrichment) HumanReadable(prepend string, linebreak string) string {
	sb := utils.NewPrependedStringBuilder(prepend)

	for i, enrichment := range []orderedmap.OrderedMap(se) {
		first := true
		for _, k := range enrichment.Keys() {
			v := utils.UnsafeGet[string](&enrichment, k)
			if first {
				sb.WriteStringf("%d. %s: %s%s", i+1, k, v, linebreak)
				first = false
			} else {
				sb.WriteStringf("   %s: %s%s", k, v, linebreak)
			}
		}
	}

	return linebreak + sb.String()
}

func NewGenericListEnrichmentFromInterface(data interface{}) (Enrichment, error) {
	if val, ok := data.([]interface{}); !ok {
		return nil, fmt.Errorf("expecting []map[string]string, found %T", data)
	} else {
		casted := utils.CastSliceOfInterface[orderedmap.OrderedMap](val)
		return GenericListEnrichment(casted), nil
	}
}
