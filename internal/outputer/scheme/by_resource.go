package scheme

import (
	"github.com/Legit-Labs/legitify/internal/common/map_utils"
	"github.com/iancoleman/orderedmap"
)

// ByResource maps Resource to the default scheme
type ByResource orderedmap.OrderedMap // Must be exported for json marshal

func NewByResource() *ByResource {
	return ToByResource(orderedmap.New())
}
func ToByResource(m *orderedmap.OrderedMap) *ByResource {
	return (*ByResource)(m)
}
func (s *ByResource) AsOrderedMap() *orderedmap.OrderedMap {
	return (*orderedmap.OrderedMap)(s)
}
func (s *ByResource) Keys() []string {
	return s.AsOrderedMap().Keys()
}
func (s *ByResource) UnsafeGet(resourceType string) *Flattened {
	return map_utils.UnsafeGet[*Flattened](s.AsOrderedMap(), resourceType)
}
