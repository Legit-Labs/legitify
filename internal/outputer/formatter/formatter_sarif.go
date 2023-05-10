package formatter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/owenrumney/go-sarif/v2/sarif"

	"github.com/Legit-Labs/legitify/internal/common/severity"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
)

type sarifFormatter struct {
	colorizer sarifColorizer
}

func newSarifFormatter() OutputFormatter {
	return &sarifFormatter{
		colorizer: sarifColorizer{},
	}
}

func (f *sarifFormatter) URIFromLink(link string) (base, uri string) {
	const (
		httpsPrefix = "https://"
		httpPrefix  = "http://"
	)

	if strings.HasPrefix(link, httpsPrefix) {
		return httpsPrefix, link[len(httpsPrefix):]
	} else if strings.HasPrefix(link, httpPrefix) {
		return httpPrefix, link[len(httpPrefix):]
	} else {
		return "", link
	}
}

func (f *sarifFormatter) Format(s scheme.Scheme, failedOnly bool) ([]byte, error) {
	report, err := sarif.New(sarif.Version210)
	if err != nil {
		return nil, err
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
		pb.Add("precision", "high")
		pb.Add("problem.severity", sarifProblemSeverity(policyInfo.Severity))
		pb.Add("security-severity", sarifSecuritySeverity(policyInfo.Severity))

		run.AddRule(policyInfo.FullyQualifiedPolicyName).
			WithDescription(policyInfo.Description).
			WithShortDescription(sarif.NewMultiformatMessageString(policyInfo.Title)).
			WithProperties(pb.Properties).
			WithTextHelp(getPlaintextPolicySummary(typedOutput, policyName)).
			WithMarkdownHelp(getMarkdownPolicySummary(typedOutput, policyName))

		// Tools like legitify don't fit perfectly into the SARIF model, so we're going to follow the
		// lead of OpenSSF's scorecard output as a starting point.
		// https://github.com/ossf/scorecard/blob/273dccda33590b7b46e98e19a9154f9da5400521/pkg/testdata/check6.sarif

		for _, violation := range data.Violations {
			base, uri := f.URIFromLink(violation.CanonicalLink)
			run.AddDistinctArtifact(violation.ViolationEntityType)
			run.CreateResultForRule(policyInfo.FullyQualifiedPolicyName).
				WithLevel(sarifSeverity(policyInfo.Severity)).
				WithMessage(sarif.NewTextMessage(getViolationMessage(&violation, &policyInfo))).
				WithHostedViewerUri(violation.CanonicalLink).
				AddLocation(
					sarif.NewLocationWithPhysicalLocation(
						sarif.NewPhysicalLocation().
							WithArtifactLocation(
								sarif.NewArtifactLocation().
									WithUri(uri).WithUriBaseId(base),
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

func (f *sarifFormatter) IsSchemeSupported(schemeType string) bool {
	return true
}

func getViolationMessage(violation *scheme.Violation, policyInfo *scheme.PolicyInfo) string {
	pc := getSarifContent()
	return pc.pf.FormatText(0, policyInfo.Description) +
		pc.pf.Linebreak() +
		string(pc.FormatViolation(violation))
}

// See https://github.com/github/docs/issues/21221
func sarifSeverity(s severity.Severity) string {
	switch s {
	case severity.Critical:
		return "error"
	case severity.High:
		return "error"
	case severity.Medium:
		return "warning"
	case severity.Low:
		return "note"
	default:
		return "none"
	}
}

func sarifProblemSeverity(s severity.Severity) string {
	switch s {
	case severity.Critical:
		return "error"
	case severity.High:
		return "error"
	case severity.Medium:
		return "warning"
	case severity.Low:
		return "recommendation"
	default:
		return "recommendation"
	}
}

func sarifSecuritySeverity(s severity.Severity) string {
	switch s {
	case severity.Critical:
		return "9.0"
	case severity.High:
		return "7.0"
	case severity.Medium:
		return "4.0"
	case severity.Low:
		return "1.0"
	default:
		return "1.0"
	}
}

func getSarifContent() *policiesContent {
	sFormatter := newSarifFormatter()
	typedFormatter := sFormatter.(*sarifFormatter)
	pf := newSarifPolicyFormatter()
	pc := newPoliciesContent(pf, typedFormatter.colorizer)
	return pc
}

func getMarkdownContent() *policiesContent {
	sFormatter := newMarkdownFormatter()
	typedFormatter := sFormatter.(*markdownFormatter)
	pf := newMarkdownPolicyFormatter()
	pc := newPoliciesContent(pf, typedFormatter.colorizer)
	return pc
}

func getPlaintextPolicySummary(output *scheme.Flattened, policyName string) string {
	return string(getSarifContent().FormatPolicy(output, policyName))
}

func getMarkdownPolicySummary(output *scheme.Flattened, policyName string) string {
	pc := getMarkdownContent()
	return string(pc.FormatPolicy(output, policyName))
}

type sarifColorizer struct {
}

func (sc sarifColorizer) colorize(tColor themeColor, text interface{}) string {
	return fmt.Sprintf("%v", text)
}

// plaintext policy formatting
type sarifPolicyFormatter struct {
	colorizer sarifColorizer
}

func newSarifPolicyFormatter() sarifPolicyFormatter {
	return sarifPolicyFormatter{colorizer: sarifColorizer{}}
}

func (sp sarifPolicyFormatter) FormatTitle(title string, severity severity.Severity) string {
	color := severityToThemeColor(severity)
	title = sp.colorizer.colorize(color, title)

	return title
}

func (sp sarifPolicyFormatter) FormatSubtitle(title string) string {
	return title
}

func (sp sarifPolicyFormatter) FormatText(depth int, format string, args ...interface{}) string {
	return indentMultilineSpecial(depth, fmt.Sprintf(format, args...), sp.Indent(1), sp.Linebreak())
}

func (sp sarifPolicyFormatter) FormatList(depth int, title string, list []string, ordered bool) string {
	if len(list) == 0 {
		return ""
	}

	var sb strings.Builder
	bullet := "*"
	sb.WriteString(sp.FormatText(depth, "%s\n", title))
	for i, step := range list {
		if ordered {
			bullet = fmt.Sprintf("%d.", i+1)
		}
		sb.WriteString(sp.FormatText(depth, "%s %s\n", bullet, step))
	}

	return sb.String()
}

func (sp sarifPolicyFormatter) Linebreak() string {
	return "  \n"
}

func (sp sarifPolicyFormatter) Separator() string {
	return "---"
}

func (sp sarifPolicyFormatter) Indent(depth int) string {
	return strings.Repeat(" ", depth)
}
