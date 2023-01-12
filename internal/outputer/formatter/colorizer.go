package formatter

type colorizer interface {
	colorize(color themeColor, text interface{}) string
}
