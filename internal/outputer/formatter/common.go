package formatter

import (
	"strings"
	"unicode"
)

const DefaultOutputIndent = "  "

func amplifyIndent(depth int) string {
	return strings.Repeat(DefaultOutputIndent, depth)
}
func indentMultiline(depth int, str string) string {
	return indentMultilineSpecial(depth, str, DefaultOutputIndent)
}

func amplifyIndentSpecial(depth int, indent string) string {
	return strings.Repeat(indent, depth)
}
func indentMultilineSpecial(depth int, str string, indent string) string {
	indent = amplifyIndentSpecial(depth, indent)

	if !strings.Contains(str, "\n") {
		return indent + str
	}

	lines := strings.Split(str, "\n")
	var sb strings.Builder
	lastIndex := len(lines) - 1
	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			sb.WriteString(indent + line)
		}
		if i < lastIndex {
			sb.WriteString("\n")
		}
	}

	return sb.String()
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
