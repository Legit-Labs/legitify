package formatter

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Legit-Labs/legitify/internal/common/severity"
	"github.com/Legit-Labs/legitify/internal/common/slice_utils"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
)

type markdownFormatter struct {
	colorizer markdownColorizer
}

func newMarkdownFormatter() OutputFormatter {
	return &markdownFormatter{
		colorizer: markdownColorizer{},
	}
}

func (m *markdownFormatter) Format(output scheme.Scheme, failedOnly bool) ([]byte, error) {
	var summary, failedPolicies []byte
	var typedOutput *scheme.Flattened

	typedOutput, ok := output.(*scheme.Flattened)
	if !ok {
		return nil, UnsupportedScheme{output}
	}

	if !failedOnly {
		summary = m.formatSummaryTable(typedOutput)
	}

	failedPolicies = m.formatFailedPolicies(typedOutput)

	return append(summary, failedPolicies...), nil
}

func (m *markdownFormatter) IsSchemeSupported(schemeType string) bool {
	return schemeType == scheme.TypeFlattened
}

func (m *markdownFormatter) formatSummaryTable(output *scheme.Flattened) []byte {
	tf := newMarkdownTableFormatter()
	tw := newTableContent(tf, m.colorizer)
	return tw.FormatSummary(output)
}

func (m *markdownFormatter) formatFailedPolicies(output *scheme.Flattened) []byte {
	failedPolicies := output.OnlyFailedViolations()
	pf := newMarkdownPolicyFormatter()
	pc := newPoliciesContent(pf, m.colorizer)
	return pc.FormatFailedPolicies(failedPolicies)
}

type markdownTableFormatter struct {
	buf    bytes.Buffer
	title  string
	header []string
}

func newMarkdownTableFormatter() *markdownTableFormatter {
	return &markdownTableFormatter{}
}

func (tf *markdownTableFormatter) SetTitle(title string) {
	tf.title = title
}

func (tf *markdownTableFormatter) SetHeaders(headers []string) {
	tf.header = headers
}

func (tf *markdownTableFormatter) WriteRow(row []string) {
	tf.buf.WriteString(asMarkdownRow(row))
}

func (tf *markdownTableFormatter) Render() []byte {
	beginning := []byte(asMarkdownTitle(tf.title) + "\n" + asMarkdownHeader(tf.header))
	lines := tf.buf.Bytes()

	return append(beginning, lines...)
}

func asMarkdownHeader(header []string) string {
	const headerRowSep = "--"
	sep := make([]string, 0, len(header))
	for range header {
		sep = append(sep, headerRowSep)
	}
	return asMarkdownRow(header) + asMarkdownRow(sep)
}

func asMarkdownRow(row []string) string {
	const colSep = "|"
	return fmt.Sprintf("%s%s%s\n", colSep, strings.Join(row, colSep), colSep)
}

func asMarkdownTitle(title string) string {
	return fmt.Sprintf("# %s", title)
}

func asMarkdownSubtitle(title string) string {
	return fmt.Sprintf("## %s", title)
}

const (
	markdownDangerEmoji  = ":no_entry:"
	markdownWarningEmoji = ":warning:"
	markdownRocketEmoji  = ":rocket:"
	markdownEyesEmoji    = ":eyes:"
)

type markdownColorizer struct {
}

func (mc markdownColorizer) withEmoji(emoji string, text interface{}) string {
	return fmt.Sprintf("%v %s", text, emoji)
}

func (mc markdownColorizer) asBold(text interface{}) string {
	asString := strings.TrimSpace(fmt.Sprintf("%v", text))
	bolded := slice_utils.Map(strings.Split(asString, "\n"), func(line string) string {
		return fmt.Sprintf("**%s**", line)
	})
	return strings.Join(bolded, "\n")
}

func (mc markdownColorizer) colorize(tColor themeColor, text interface{}) string {
	switch tColor {
	case themeColorFailure:
		return mc.withEmoji(markdownDangerEmoji, text)
	case themeColorAlert:
		return mc.withEmoji(markdownWarningEmoji, text)
	case themeColorSuccess:
		return mc.withEmoji(markdownRocketEmoji, text)
	case themeColorInteresting:
		return mc.withEmoji(markdownEyesEmoji, text)
	case themeColorWarning:
		return mc.withEmoji(markdownEyesEmoji, text)
	case themeColorBold:
		return mc.asBold(text)

	case themeColorNeutral:
		fallthrough
	default:
		return fmt.Sprintf("%v", text)
	}
}

// policy formatting
type markdownPolicyFormatter struct {
	colorizer markdownColorizer
}

func newMarkdownPolicyFormatter() markdownPolicyFormatter {
	return markdownPolicyFormatter{colorizer: markdownColorizer{}}
}

func (mp markdownPolicyFormatter) FormatTitle(title string, severity severity.Severity) string {
	color := severityToThemeColor(severity)
	title = mp.colorizer.colorize(color, title)

	return asMarkdownTitle(title)
}

func (mp markdownPolicyFormatter) FormatSubtitle(title string) string {
	return asMarkdownSubtitle(title)
}

func (mp markdownPolicyFormatter) FormatText(depth int, format string, args ...interface{}) string {
	return indentMultilineSpecial(depth, fmt.Sprintf(format, args...), mp.Indent(1), mp.Linebreak())
}

func (mp markdownPolicyFormatter) FormatList(depth int, title string, list []string, ordered bool) string {
	if len(list) == 0 {
		return ""
	}

	var sb strings.Builder
	bullet := "-"
	sb.WriteString(mp.FormatText(depth, "%s\n", title))
	for i, step := range list {
		if ordered {
			bullet = fmt.Sprintf("%d.", i+1)
		}
		sb.WriteString(mp.FormatText(depth, "%s %s\n", bullet, step))
	}

	return sb.String()
}

func (mp markdownPolicyFormatter) Linebreak() string {
	return "  \n"
}

func (mp markdownPolicyFormatter) Separator() string {
	return "---"
}

func (mp markdownPolicyFormatter) Indent(depth int) string {
	return strings.Repeat("> ", depth)
}
