package screen

import (
	"io"
	"log"
	"os"

	"github.com/Legit-Labs/legitify/cmd/tty"
)

type screen struct {
	log   *log.Logger
	isTty bool
}

var singletone screen

func init() {
	Init(os.Stderr, tty.IsStderrTty())
}

func Init(writer io.Writer, isTty bool) {
	singletone.log = log.New(writer, "", 0)
	singletone.isTty = isTty
}

func Printf(fmt string, args ...interface{}) {
	singletone.log.Printf(fmt, args...)
}

func Writer() io.Writer {
	return singletone.log.Writer()
}

func IsTty() bool {
	return singletone.isTty
}
