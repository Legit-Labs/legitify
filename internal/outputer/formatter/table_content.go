package formatter

import (
	"github.com/Legit-Labs/legitify/internal/analyzers"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
)

type tableFormatter interface {
	SetTitle(title string)
	SetHeaders(headers []string)
	WriteRow(line []string)
	Render() []byte
}

type tableContent struct {
	tf        tableFormatter
	colorizer Colorizer
}

func NewTableContent(tf tableFormatter, colorizer Colorizer) *tableContent {
	return &tableContent{
		tf:        tf,
		colorizer: colorizer,
	}
}

func (tc *tableContent) countColorize(count int, color themeColor) string {
	if count == 0 {
		return tc.colorizer.colorize(themeColorNone, count)
	} else {
		return tc.colorizer.colorize(color, count)
	}
}

func (tc *tableContent) FormatSummary(output scheme.FlattenedScheme) []byte {
	output = scheme.SortSchemeByNamespace(output, false)

	tc.tf.SetTitle(tc.colorizer.colorize(themeColorBold, "Findings Summary"))

	headers := []string{"#", "Namespace", "Policy", "Severity", "Passed", "Failed", "Skipped"}
	for i, h := range headers {
		headers[i] = tc.colorizer.colorize(themeColorBold, h)
	}
	tc.tf.SetHeaders(headers)

	for i, policyName := range output.Keys() {
		rowNum := tc.colorizer.colorize(themeColorBold, i+1)
		data := output.GetPolicyData(policyName)
		policyInfo := data.PolicyInfo
		title := policyInfo.Title
		severity := tc.colorizer.colorize(severityToThemeColor(policyInfo.Severity), policyInfo.Severity)
		namespace := policyInfo.Namespace

		var passed, failed, skipped int
		for _, violation := range data.Violations {
			switch violation.Status {
			case analyzers.PolicyPassed:
				passed++
			case analyzers.PolicyFailed:
				failed++
			case analyzers.PolicySkipped:
				skipped++
			}
		}

		passedStr := tc.countColorize(passed, themeColorSuccess)
		failedStr := tc.countColorize(failed, themeColorFailure)
		skippedStr := tc.countColorize(skipped, themeColorInteresting)

		tc.tf.WriteRow([]string{rowNum, namespace, title, severity, passedStr, failedStr, skippedStr})
	}

	return tc.tf.Render()
}
