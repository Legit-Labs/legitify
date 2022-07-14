package enrichers

import (
	"fmt"

	"github.com/Legit-Labs/legitify/internal/analyzers"
	"github.com/Legit-Labs/legitify/internal/common/utils"
)

type Enrichment interface {
	HumanReadable(prepend string) string
}

type BasicEnrichment struct {
	val string
}

func (s *BasicEnrichment) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, s.val)), nil
}

func (be *BasicEnrichment) HumanReadable(prepend string) string {
	sb := utils.NewPrependedStringBuilder(prepend)
	sb.WriteString(be.val)
	return sb.String()
}

func NewBasicEnrichment(str string) Enrichment {
	return &BasicEnrichment{
		val: str,
	}
}

type Enricher interface {
	Enrich(data analyzers.AnalyzedData) (enrichment Enrichment, ok bool)
	ShouldEnrich(requestedEnricher string) bool
	Name() string
}
