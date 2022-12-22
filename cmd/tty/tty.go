package tty

import (
	"os"

	"github.com/mattn/go-isatty"
)

func IsStdoutTty() bool {
	return IsTty(os.Stdout)
}
func IsStderrTty() bool {
	return IsTty(os.Stderr)
}

func IsTty(file *os.File) bool {
	// Inspired by the color package:
	// color package decides whether or not to use colors based on stdout,
	// but it does it on import time, which is too early for us.
	return os.Getenv("TERM") != "dumb" &&
		(isatty.IsTerminal(file.Fd()) || isatty.IsCygwinTerminal(file.Fd()))
}
