package enrichers

import (
	"fmt"

	"github.com/Legit-Labs/legitify/internal/analyzers"
	"github.com/Legit-Labs/legitify/internal/common/utils"
)

type Enrichment interface {
	HumanReadable(prepend string) string
	Name() string
}

type BasicEnrichment struct {
	val  string
	name string
}

func (s *BasicEnrichment) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, s.val)), nil
}

func (be *BasicEnrichment) HumanReadable(_ string) string {
	sb := utils.NewPrependedStringBuilder("")
	sb.WriteString(be.val)
	return sb.String()
}

func (be *BasicEnrichment) Name() string {
	return be.name
}

func NewBasicEnrichment(str string, name string) Enrichment {
	return &BasicEnrichment{
		val:  str,
		name: name,
	}
}

type Enricher interface {
	Enrich(data analyzers.AnalyzedData) (enrichment Enrichment, ok bool)
	Name() string
}
