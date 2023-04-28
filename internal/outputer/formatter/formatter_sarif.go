package formatter

import (
	"encoding/json"
	"strings"

	"github.com/owenrumney/go-sarif/v2/sarif"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
)

type SarifFormatter struct {
}

func newSarifFormatter() OutputFormatter {
	return &SarifFormatter{}
}

func (f *SarifFormatter) Format(s scheme.Scheme, failedOnly bool) ([]byte, error) {
	report, err := sarif.New(sarif.Version210)
	if err != nil {
		panic(err)
	}

	typedOutput, ok := s.(*scheme.Flattened)
	if !ok {
		return nil, UnsupportedScheme{s}
	}

	run := sarif.NewRunWithInformationURI("legitify", "https://legitify.dev/")

	for _, policyName := range s.AsOrderedMap().Keys() { 
		data := typedOutput.GetPolicyData(policyName)
		policyInfo := data.PolicyInfo

		pb := sarif.NewPropertyBag()
		pb.Add("impact", policyInfo.Threat)
		pb.Add("resolution", policyInfo.RemediationSteps)

		run.AddRule(policyInfo.FullyQualifiedPolicyName).
			WithDescription(policyInfo.Description).
			WithProperties(pb.Properties).
			WithMarkdownHelp(getMarkdownPolicySummary(s.(*scheme.Flattened), policyName))

		for _, violation := range data.Violations {
			run.AddDistinctArtifact(violation.ViolationEntityType)
			run.CreateResultForRule(policyInfo.FullyQualifiedPolicyName).
				WithLevel(strings.ToLower(policyInfo.Severity)).
				WithMessage(sarif.NewTextMessage(policyInfo.Description)).
				AddLocation(
					sarif.NewLocationWithPhysicalLocation(
						sarif.NewPhysicalLocation().WithArtifactLocation(
							sarif.NewSimpleArtifactLocation(violation.CanonicalLink),
						),
					),
				)
		}
	}

	report.AddRun(run)

	bytes, err := json.MarshalIndent(report, "", DefaultOutputIndent)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func getMarkdownPolicySummary(output *scheme.Flattened, policyName string) string {
	mdFormatter := newMarkdownFormatter()
	typedFormatter := mdFormatter.(*markdownFormatter)
	pf := newMarkdownPolicyFormatter()
	pc := newPoliciesContent(pf, typedFormatter.colorizer)
	return string(pc.FormatPolicy(output, policyName))
}

func (f *SarifFormatter) IsSchemeSupported(schemeType string) bool {
	return true
}
