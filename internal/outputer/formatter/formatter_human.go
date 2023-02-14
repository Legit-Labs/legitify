package formatter

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Legit-Labs/legitify/internal/common/severity"
	tw "github.com/olekukonko/tablewriter"

	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
	"github.com/fatih/color"
)

type HumanFormatter struct {
	colorizer humanColorizer
}

func newHumanFormatter() OutputFormatter {
	return &HumanFormatter{
		colorizer: humanColorizer{},
	}
}

func (f *HumanFormatter) formatSummaryTable(output *scheme.Flattened) []byte {
	tf := newHumanTableWriter()
	tw := newTableContent(tf, f.colorizer)
	return tw.FormatSummary(output)
}

func (f *HumanFormatter) formatFailedPolicies(output *scheme.Flattened) []byte {
	failedPolicies := output.OnlyFailedViolations()
	pf := newHumanPolicyFormatter()
	pc := newPoliciesContent(pf, f.colorizer)
	return pc.FormatFailedPolicies(failedPolicies)
}

func (f *HumanFormatter) Format(output scheme.Scheme, failedOnly bool) ([]byte, error) {
	var summary, failedPolicies []byte

	typedOutput, ok := output.(*scheme.Flattened)
	if !ok {
		return nil, UnsupportedScheme{output}
	}

	if !failedOnly {
		summary = f.formatSummaryTable(typedOutput)
	}

	failedPolicies = f.formatFailedPolicies(typedOutput)

	return append(failedPolicies, summary...), nil
}

func (f *HumanFormatter) IsSchemeSupported(schemeType string) bool {
	return schemeType == scheme.TypeFlattened
}

// table formatting

type textTableFormatter struct {
	title string
	buf   bytes.Buffer
	tw    *tw.Table
}

func newHumanTableWriter() *textTableFormatter {
	var tf textTableFormatter
	tw := tw.NewWriter(&tf.buf)
	tf.tw = tw
	return &tf
}

func (tf *textTableFormatter) SetTitle(title string) {
	tf.title = title
}

func (tf *textTableFormatter) SetHeaders(headers []string) {
	tf.tw.SetHeader(headers)
	tf.tw.SetAutoFormatHeaders(false)
	tf.tw.SetRowLine(true)
}

func (tf *textTableFormatter) WriteRow(row []string) {
	tf.tw.Append(row)
}

func (tf *textTableFormatter) Render() []byte {
	tf.tw.Render()
	title := []byte(fmt.Sprintf("\n%s:\n", tf.title))
	return append(title, tf.buf.Bytes()...)
}

type humanColorizer struct {
}

func (hc humanColorizer) mapThemeColor(tColor themeColor) color.Attribute {
	switch tColor {
	case themeColorBold:
		return color.Bold
	case themeColorSuccess:
		return color.FgHiGreen
	case themeColorFailure:
		return color.FgHiRed
	case themeColorInteresting:
		return color.FgHiBlue
	case themeColorNeutral:
		return color.FgHiBlue
	case themeColorAlert:
		return color.FgHiYellow
	case themeColorWarning:
		return color.FgYellow
	default:
		return color.Reset
	}
}

func (hc humanColorizer) colorize(tColor themeColor, text interface{}) string {
	return color.New(hc.mapThemeColor(tColor)).Sprintf("%v", text)
}

// policy formatting
type humanPolicyFormatter struct {
	colorizer humanColorizer
}

func newHumanPolicyFormatter() humanPolicyFormatter {
	return humanPolicyFormatter{colorizer: humanColorizer{}}
}

func (hp humanPolicyFormatter) FormatTitle(title string, severity severity.Severity) string {
	sep := strings.Repeat("-", len(title))

	boldText := hp.colorizer.colorize(themeColorBold, title)
	color := severityToThemeColor(severity)
	title = hp.colorizer.colorize(color, boldText)

	return fmt.Sprintf("%s\n%s", title, sep)
}

func (hp humanPolicyFormatter) FormatSubtitle(title string) string {
	return hp.colorizer.colorize(themeColorBold, title) + hp.Linebreak() + hp.Separator()
}

func (hp humanPolicyFormatter) FormatText(depth int, format string, args ...interface{}) string {
	return indentMultiline(depth, fmt.Sprintf(format, args...))
}

func (hp humanPolicyFormatter) FormatList(depth int, title string, list []string, ordered bool) string {
	if len(list) == 0 {
		return ""
	}

	var sb strings.Builder
	bullet := "-"
	sb.WriteString(hp.FormatText(depth, "%s\n", title))
	for i, step := range list {
		if ordered {
			bullet = fmt.Sprintf("%d.", i+1)
		}
		sb.WriteString(hp.FormatText(depth, "%s %s\n", bullet, step))
	}

	return sb.String()
}

func (hp humanPolicyFormatter) Linebreak() string {
	return "\n"
}

func (hp humanPolicyFormatter) Separator() string {
	return "------------"
}

func (hp humanPolicyFormatter) Indent(depth int) string {
	return amplifyIndent(depth)
}
