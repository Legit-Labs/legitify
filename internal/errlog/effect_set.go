package errlog

import (
	"encoding/json"

	"github.com/Legit-Labs/legitify/internal/common/map_utils"
)

type effectSet map[string]bool

func (es effectSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(es.filtered())
}

func (es effectSet) filtered() []string {
	filteredEffects := []string{}
	keys := map_utils.ToKeySortedMap[bool]((map[string]bool)(es))
	for _, effect := range keys.Keys() {
		if effect == "" {
			continue
		}
		filteredEffects = append(filteredEffects, effect)
	}
	return filteredEffects
}
