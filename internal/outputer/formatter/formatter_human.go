package formatter

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"github.com/Legit-Labs/legitify/internal/analyzers"
	"github.com/Legit-Labs/legitify/internal/enricher/enrichers"
	"github.com/olekukonko/tablewriter"

	"github.com/Legit-Labs/legitify/internal/common/severity"

	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme/converter"
	"github.com/fatih/color"
)

var severityToColor = map[severity.Severity]color.Attribute{
	severity.Critical: color.FgRed,
	severity.High:     color.FgHiRed,
	severity.Medium:   color.FgHiYellow,
	severity.Low:      color.FgYellow,
	severity.Unknown:  color.FgWhite,
}

type HumanFormatter struct {
	indent string
	sb     strings.Builder
}

func NewHumanFormatter(indent string) OutputFormatter {
	return &HumanFormatter{indent: indent}
}

func (f *HumanFormatter) amplifyIndent(depth int) string {
	return strings.Repeat(f.indent, depth)
}

func colorize(data interface{}, att color.Attribute) string {
	return color.New(att).Sprintf("%v", data)
}

func bold(data interface{}) string {
	return colorize(data, color.Bold)
}

func (f *HumanFormatter) sprintfWithColor(depth int, colorAttribute color.Attribute, format string, args ...interface{}) string {
	formatted := f.sprintf(depth, format, args...)
	return colorize(formatted, colorAttribute)
}

func (f *HumanFormatter) sprintf(depth int, format string, args ...interface{}) string {
	return fmt.Sprintf("%s%s", f.amplifyIndent(depth), fmt.Sprintf(format, args...))
}

func (f *HumanFormatter) indentMultiline(depth int, str string) string {
	indent := f.amplifyIndent(depth)
	lines := strings.Split(str, "\n")
	return strings.Join(lines, "\n"+indent)
}

func camelCaseToTitle(camelCased string) string {
	var sb strings.Builder

	for i, c := range camelCased {
		if unicode.IsLower(c) {
			if i == 0 {
				c = unicode.ToUpper(c)
			}
			sb.WriteRune(c)
		} else {
			sb.WriteRune(' ')
			sb.WriteRune(c)
		}
	}

	return sb.String()
}

func (f *HumanFormatter) formatAux(m map[string]enrichers.Enrichment) {
	for k, v := range m {
		indent := f.amplifyIndent(4)
		v := v.HumanReadable(indent)

		format := "- %s: "
		if strings.Contains(v, "\n") {
			format = "- %s:\n"
		}

		v = strings.TrimSuffix(v, "\n")
		f.sb.WriteString(f.sprintf(3, format, camelCaseToTitle(k)))
		f.sb.WriteString(v)
		f.sb.WriteString("\n")
	}
}

func (f *HumanFormatter) colorByPolicy(policyInfo scheme.PolicyInfo) color.Attribute {
	return severityToColor[policyInfo.Severity]
}

func (f *HumanFormatter) formatPolicyInfo(policyName string, policyInfo scheme.PolicyInfo) {
	f.sb.WriteString(f.sprintfWithColor(0, f.colorByPolicy(policyInfo), "%s\n", policyInfo.Title))
	f.sb.WriteString(strings.Repeat("-", len(policyInfo.Title)) + "\n")
	f.sb.WriteString(f.sprintf(1, "%s\n", f.indentMultiline(1, policyInfo.Description)))

	f.sb.WriteString(f.sprintf(1, "Policy Name: %s\n", policyName))
	f.sb.WriteString(f.sprintf(1, "Namespace: %s\n", policyInfo.Namespace))
	f.sb.WriteString(f.sprintfWithColor(1, f.colorByPolicy(policyInfo), "Severity: %s\n", policyInfo.Severity))
	f.sb.WriteString(f.sprintf(1, "Remediation Steps:\n"))
	for i, step := range policyInfo.RemediationSteps {
		f.sb.WriteString(f.sprintf(2, "%d. %s\n", i+1, step))
	}
}

func (f *HumanFormatter) formatViolation(violation scheme.Violation) {
	f.sb.WriteString(f.sprintf(2, "%sLink to %s: %s\n", f.indent, violation.ViolationEntityType, violation.CanonicalLink))
	if len(violation.Aux) > 0 {
		f.sb.WriteString(f.sprintf(2, "%sAuxiliary Info:\n", f.indent))
		f.formatAux(violation.Aux)
	}
}

func (f *HumanFormatter) formatSummaryTable(output scheme.FlattenedScheme) []byte {
	var buf bytes.Buffer

	output = scheme.SortSchemeByNamespace(output, false)
	tw := tablewriter.NewWriter(&buf)

	headers := []string{"#", "Namespace", "Policy", "Severity", "Passed", "Failed", "Skipped"}
	for i, h := range headers {
		headers[i] = bold(h)
	}
	tw.SetHeader(headers)
	tw.SetAutoFormatHeaders(false)
	tw.SetRowLine(true)

	for i, policyName := range output.Keys() {
		rowNum := bold(i + 1)
		data := output.GetPolicyData(policyName)
		policyInfo := data.PolicyInfo
		colorAtt := f.colorByPolicy(policyInfo)
		title := policyInfo.Title
		severity := colorize(policyInfo.Severity, colorAtt)
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

		passedStr := colorize(passed, color.FgGreen)
		failedStr := colorize(failed, color.FgRed)
		skippedStr := colorize(skipped, color.FgHiBlue)

		tw.Append([]string{rowNum, namespace, title, severity, passedStr, failedStr, skippedStr})
	}

	tw.Render()

	separator := []byte(color.New(color.Bold).Sprintf("\nFindings summary:\n"))
	return append(separator, buf.Bytes()...)
}

func (f *HumanFormatter) formatFailedViolations(output scheme.FlattenedScheme) ([]byte, error) {
	f.sb.Reset()

	i := 0
	lastIndex := len(output.Keys()) - 1
	for _, policyName := range output.Keys() {
		data := output.GetPolicyData(policyName)
		f.formatPolicyInfo(policyName, data.PolicyInfo)
		f.sb.WriteString("\n")

		f.sb.WriteString(f.sprintf(1, "Violations:\n"))
		for i, violation := range data.Violations {
			f.formatViolation(violation)
			if i < len(data.Violations)-1 {
				f.sb.WriteString(f.sprintf(2, "---\n"))
			}
		}
		if i < lastIndex {
			f.sb.WriteString("\n")
		}
		i++
	}

	return []byte(f.sb.String()), nil
}

func (f *HumanFormatter) Format(output interface{}, failedOnly bool) ([]byte, error) {
	var summary, failedViolations []byte
	var typedOutput scheme.FlattenedScheme

	typedOutput, ok := output.(scheme.FlattenedScheme)
	if !ok {
		return nil, UnsupportedScheme{output}
	}

	if !failedOnly {
		summary = f.formatSummaryTable(typedOutput)
		typedOutput = scheme.OnlyFailedViolations(typedOutput)
	}

	failedViolations, err := f.formatFailedViolations(typedOutput)
	if err != nil {
		return nil, err
	}

	return append(failedViolations, summary...), err
}

func (f *HumanFormatter) IsSchemeSupported(schemeType string) bool {
	return schemeType == converter.Flattened
}
