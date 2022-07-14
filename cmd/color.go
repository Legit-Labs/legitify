package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
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

func isTty() bool {
	// Inpsired by the color package:
	// color package decides whether or not to use colors based on stdout,
	// but it does it on import time, which is too early for us.
	return os.Getenv("TERM") != "dumb" &&
		(isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()))
}

func InitColorPackage(colorWhen string) error {
	switch colorWhen {
	case colorAlways:
		color.NoColor = false
	case colorNone:
		color.NoColor = true
	case colorAuto:
		color.NoColor = !isTty()
	default:
		return fmt.Errorf("invalid color option: %s", colorWhen)
	}

	return nil
}
