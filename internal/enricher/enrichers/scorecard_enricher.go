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

func NewScorecardEnricher(ctx context.Context) Enricher {
	return &scorecardEnricher{
		enabled: context_utils.GetScorecardVerbose(ctx),
	}
}

type scorecardEnricher struct {
	enabled bool
}

func (e *scorecardEnricher) Enrich(data analyzers.AnalyzedData) (Enrichment, bool) {
	if !e.enabled {
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

func (e *scorecardEnricher) ShouldEnrich(requestedEnricher string) bool {
	return requestedEnricher == e.Name()
}

func (e *scorecardEnricher) Name() string {
	return Scorecard
}

func createEnrichment(sc *sc.Result) (*ScorecardEnrichment, error) {
	var result ScorecardEnrichment
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

		result.Checks = append(result.Checks, ScorecardCheck{
			Reason:  checkResult.Reason,
			DocsUrl: doc.GetDocumentationURL(sc.Result.Scorecard.CommitSHA),
			Details: details,
		})
	}

	return &result, nil
}

type ScorecardCheck struct {
	Reason  string
	DocsUrl string
	Details []string
}

type ScorecardEnrichment struct {
	Checks []ScorecardCheck
}

func (se *ScorecardEnrichment) HumanReadable(prepend string) string {
	sb := utils.NewPrependedStringBuilder(prepend)

	for j, checkResult := range se.Checks {
		sb.WriteString(fmt.Sprintf("%d. %s:\n", j+1, checkResult.Reason))
		sb.WriteString(fmt.Sprintf("docs: %s\n", checkResult.DocsUrl))

		if len(checkResult.Details) > 0 {
			sb.WriteString(fmt.Sprintln("details: "))
			for i, detail := range checkResult.Details {
				clean := strings.Replace(detail, "\t", "", -1)
				clean = strings.Replace(clean, "\n", " ", -1)
				sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, clean))
			}
		}
	}

	return sb.String()
}
