package cmd

import (
	"fmt"

	"github.com/Legit-Labs/legitify/cmd/tty"
	"github.com/fatih/color"
)

const (
	colorAuto          = "auto"
	colorAlways        = "always"
	colorNone          = "none"
	DefaultColorOption = colorAuto
)

func ColorOptions() []string {
	return []string{colorAuto, colorAlways, colorNone}
}

func InitColorPackage(colorWhen string) error {
	switch colorWhen {
	case colorAlways:
		color.NoColor = false
	case colorNone:
		color.NoColor = true
	case colorAuto:
		color.NoColor = !tty.IsStdoutTty()
	default:
		return fmt.Errorf("invalid color option: %s", colorWhen)
	}

	return nil
}
