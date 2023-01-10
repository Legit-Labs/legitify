package formatter

import "github.com/Legit-Labs/legitify/internal/common/severity"

type themeColor = int

const (
	themeColorSuccess     themeColor = iota
	themeColorFailure     themeColor = iota
	themeColorNeutral     themeColor = iota
	themeColorInteresting themeColor = iota

	themeColorAlert   themeColor = iota
	themeColorWarning themeColor = iota
	themeColorBold    themeColor = iota

	themeColorNone themeColor = iota
)

type Colorizer interface {
	colorize(color themeColor, text interface{}) string
}

func severityToThemeColor(sev severity.Severity) themeColor {
	switch sev {
	case severity.Critical:
		fallthrough
	case severity.High:
		return themeColorFailure
	case severity.Medium:
		return themeColorAlert
	case severity.Low:
		return themeColorWarning
	default:
		return themeColorNone
	}
}
