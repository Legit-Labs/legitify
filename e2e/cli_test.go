package test

import (
	"bytes"
	"flag"
	"log"
	"os/exec"
	"strings"
	"testing"
)

var legitifyCLIPath = flag.String("legitify_cli_path", "/tmp/legitify", "legitify cli tool path")

func TestCLI(t *testing.T) {
	for _, test := range []struct {
		Args         []string
		ExpectedPass bool
		Output       string
	}{
		{
			Args:         []string{"5"},
			ExpectedPass: false,
			Output:       "unknown command \"5\" for \"legitify\"",
		},
	} {
		t.Run("", func(t *testing.T) {
			cmd := exec.Command(*legitifyCLIPath, test.Args...)
			var outb, errb bytes.Buffer
			cmd.Stdout = &outb
			cmd.Stderr = &errb
			err := cmd.Run()
			if !test.ExpectedPass {
				if err == nil {
					log.Println("out:", outb.String(), "err:", errb.String())
					log.Fatal(err)
				}
				if !strings.Contains(errb.String(), test.Output) {
					log.Println("out:", outb.String(), "err:", errb.String())
					log.Fatal(err)
				}
			} else {

				if err != nil {
					log.Println("out:", outb.String(), "err:", errb.String())
					log.Fatal(err)
				}
			}
		})
	}
}
