package errlog

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type forwarder struct{}

func (f *forwarder) Write(p []byte) (n int, err error) {
	Printf("%s", string(p))
	return len(p), nil
}

type errlog struct {
	log         *log.Logger
	everWritten bool
	permIssues  bool
	skiplog     *SkipLog
	permLog     *PermLog
	permWriter  io.Writer
}

var singletone errlog

func init() {
	singletone.log = log.New(os.Stderr, "", log.LstdFlags)
	singletone.skiplog = NewSkipLog()
	singletone.permLog = NewPermLog()
	log.SetOutput(&forwarder{})
}

func SetOutput(writer io.Writer) {
	singletone.log.SetOutput(writer)
}

func SetPermissionsOutput(writer io.Writer) {
	singletone.permWriter = writer
}

func Printf(format string, args ...interface{}) {
	singletone.everWritten = true
	singletone.log.Printf(format, args...)
}

func AddPermIssue(issue PermIssue) {
	singletone.permLog.Add(issue)
}

func AddSkipIssue(policyName string, entityName string, skipReason SkipReason) {
	singletone.skiplog.Add(policyName, entityName, skipReason)
}

type PermissionsOutput struct {
	Permissions     interface{} `json:"missing_permissions"`
	SkippedPolicies interface{} `json:"skipped_policies"`
}

func FlushAll() {
	if singletone.permLog.Empty() && singletone.skiplog.Empty() {
		return
	}
	singletone.permIssues = true

	issuesOutput := PermissionsOutput{
		Permissions:     singletone.permLog,
		SkippedPolicies: singletone.skiplog,
	}

	permIssues, err := json.MarshalIndent(issuesOutput, "", "  ")
	if err != nil {
		singletone.log.Printf("Failed to marshal permission issues: %s", err)
	} else {
		if _, err := singletone.permWriter.Write(permIssues); err != nil {
			singletone.log.Printf("Failed to dump permission issues: %s", err)
		}
	}
}

func HadErrors() bool {
	return singletone.everWritten
}
func HadPermIssues() bool {
	return singletone.permIssues
}
