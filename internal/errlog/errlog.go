package errlog

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type forwarder struct{}

func (f *forwarder) Write(p []byte) (n int, err error) {
	Printf("%s", string(p))
	return len(p), nil
}

type errlog struct {
	log         *log.Logger
	everWritten bool
	permLogs    strings.Builder
	prereqLogs  strings.Builder
}

var singletone errlog

func init() {
	singletone.log = log.New(os.Stderr, "", log.LstdFlags)
	log.SetOutput(&forwarder{})

	PermIssueF("Missing permissions errors:\n")
	PrereqIssueF("Unmet Prerequisites errors:\n")
}

func SetOutput(writer io.Writer) {
	singletone.log.SetOutput(writer)
}

func Printf(format string, args ...interface{}) {
	singletone.everWritten = true
	singletone.log.Printf(format, args...)
}

func PermIssueF(format string, args ...interface{}) {
	singletone.permLogs.WriteString(fmt.Sprintf(format, args...))
}

func PrereqIssueF(format string, args ...interface{}) {
	singletone.prereqLogs.WriteString(fmt.Sprintf(format, args...))
}

func FlushAll() {
	singletone.log.Printf(singletone.permLogs.String())
	singletone.permLogs.Reset()

	singletone.log.Printf(singletone.prereqLogs.String())
	singletone.prereqLogs.Reset()
}

func HadErrors() bool {
	return singletone.everWritten
}
