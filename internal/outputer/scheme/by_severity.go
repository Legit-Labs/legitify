package scheme

import (
	"github.com/Legit-Labs/legitify/internal/common/severity"
	"github.com/Legit-Labs/legitify/internal/common/utils"
	"github.com/iancoleman/orderedmap"
)

// BySeverity maps Resource to the default scheme
type BySeverity orderedmap.OrderedMap // Must be exported for json marshal

func NewBySeverity() *BySeverity {
	return ToBySeverity(orderedmap.New())
}
func ToBySeverity(m *orderedmap.OrderedMap) *BySeverity {
	return (*BySeverity)(m)
}
func (s *BySeverity) AsOrderedMap() *orderedmap.OrderedMap {
	return (*orderedmap.OrderedMap)(s)
}
func (s *BySeverity) Keys() []string {
	return s.AsOrderedMap().Keys()
}
func (s *BySeverity) UnsafeGet(severity string) *Flattened {
	return utils.UnsafeGet[*Flattened](s.AsOrderedMap(), severity)
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
