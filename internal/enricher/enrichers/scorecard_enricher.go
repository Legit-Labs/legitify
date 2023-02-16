package enrichers

import (
	"context"
	"fmt"
	"strings"

	"github.com/Legit-Labs/legitify/internal/analyzers"
	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/utils"
	"github.com/Legit-Labs/legitify/internal/context_utils"
	sc "github.com/Legit-Labs/legitify/internal/scorecard"
	"github.com/ossf/scorecard/v4/checker"
	docs "github.com/ossf/scorecard/v4/docs/checks"
)

const Scorecard = "scorecard"
const maxScore = 10

func NewScorecardEnricher() Enricher {
	return &scorecardEnricher{}
}

type scorecardEnricher struct {
}

func (e *scorecardEnricher) Enrich(ctx context.Context, data analyzers.AnalyzedData) (Enrichment, bool) {
	if !context_utils.GetScorecardVerbose(ctx) {
		return nil, false
	}

	switch t := data.Entity.(type) {
	case githubcollected.Repository:
		if t.Scorecard == nil {
			return nil, false
		}
		result, err := createEnrichment(t.Scorecard)
		if err != nil {
			return nil, false
		}
		return result, true
	}
	return nil, false
}

func (e *scorecardEnricher) Parse(data interface{}) (Enrichment, error) {
	if val, ok := data.([]ScorecardCheck); !ok {
		return nil, fmt.Errorf("expecting []ScorecardCheck")
	} else {
		return ScorecardEnrichment(val), nil
	}
}

func createEnrichment(sc *sc.Result) (ScorecardEnrichment, error) {
	var checks []ScorecardCheck
	d, err := docs.Read()
	if err != nil {
		return nil, err
	}

	for _, checkResult := range sc.Result.Checks {
		if checkResult.Score == maxScore {
			continue
		}

		doc, err := d.GetCheck(checkResult.Name)
		if err != nil {
			return nil, err
		}

		var details []string
		for _, detail := range checkResult.Details {
			if detail.Type == checker.DetailWarn {
				details = append(details, detail.Msg.Text)
			}
		}

		checks = append(checks, ScorecardCheck{
			Reason:  checkResult.Reason,
			DocsUrl: doc.GetDocumentationURL(sc.Result.Scorecard.CommitSHA),
			Details: details,
		})
	}

	return checks, nil
}

type ScorecardCheck struct {
	Reason  string
	DocsUrl string
	Details []string
}

type ScorecardEnrichment []ScorecardCheck

func (se ScorecardEnrichment) HumanReadable(prepend string, linebreak string) string {
	sb := utils.NewPrependedStringBuilder(prepend)

	for j, checkResult := range []ScorecardCheck(se) {
		sb.WriteStringf("%d. %s:%s", j+1, checkResult.Reason, linebreak)
		sb.WriteStringf("docs: %s%s", checkResult.DocsUrl, linebreak)

		if len(checkResult.Details) > 0 {
			sb.WriteStringf("details: %s", linebreak)
			for i, detail := range checkResult.Details {
				clean := strings.Replace(detail, "\t", "", -1)
				clean = strings.Replace(clean, "\n", " ", -1)
				sb.WriteStringf("  %d. %s%s", i+1, clean, linebreak)
			}
		}
	}

	return linebreak + sb.String()
}
