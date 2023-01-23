package scheme

import (
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/common/utils"
	"github.com/iancoleman/orderedmap"
)

// ByNamespace maps Namespace to the default scheme
type ByNamespace orderedmap.OrderedMap // Must be exported for json marshal

func NewByNamespace() *ByNamespace {
	return ToByNamespace(orderedmap.New())
}
func ToByNamespace(m *orderedmap.OrderedMap) *ByNamespace {
	return (*ByNamespace)(m)
}
func (s *ByNamespace) AsOrderedMap() *orderedmap.OrderedMap {
	return (*orderedmap.OrderedMap)(s)
}
func (s *ByNamespace) Keys() []string {
	return s.AsOrderedMap().Keys()
}
func (s *ByNamespace) UnsafeGet(namespace string) *Flattened {
	return utils.UnsafeGet[*Flattened](s.AsOrderedMap(), namespace)
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
