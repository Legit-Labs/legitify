package test

import (
	"flag"
	"strings"
	"testing"

	"github.com/thedevsaddam/gojsonq/v2"
)

var reportPath = flag.String("report_path", "/tmp/out.json", "legitify report output path")
var executionArgs = flag.String("execution_args", "", "arguments used to run legitify")

const pathToEntityName = "aux->entityName"

func TestGitHub(t *testing.T) {
	tests := [][]testCase{
		testCasesGitHubOrganization,
		testCasesGitHubActions,
		testCasesGitHubRunnerGroup,
		testCasesGitHubRepository,
	}

	for _, testCases := range tests {
		assertionLoop(t, testCases)
	}
}

func assertTestStatus(t *testing.T, jq *gojsonq.JSONQ, testPath, entityName, expectedStatus string) {
	jq.Reset()
	testFormattedPath := "content->" + testPath + "->violations"
	violation := jq.From(testFormattedPath)
	entity := violation.Where(pathToEntityName, "=", entityName)
	count := entity.Where("status", "=", expectedStatus).Count()
	if count == 0 {
		t.Logf("Failed on test %s Entity %s did not pass expected %s count %d", testPath, entityName, expectedStatus, count)
		t.Fail()
	}
}

type testPair struct {
	Got, Want string
}

func assertionLoop(t *testing.T, tests []testCase) {
	jq := gojsonq.New(gojsonq.SetSeparator("->")).File(*reportPath)
	for _, test := range tests {
		t.Logf("Testing: %s", test.path)

		pairs := []testPair{
			{
				Got:  test.passedEntity,
				Want: "PASSED",
			},
			{
				Got:  test.failedEntity,
				Want: "FAILED",
			},
			{
				Got:  test.skippedEntity,
				Want: "SKIPPED",
			},
		}

		for _, pair := range pairs {
			if pair.Got == "" {
				continue
			}
			assertTestStatus(t, jq, test.path, pair.Got, pair.Want)
		}
	}
}

func TestGitLab(t *testing.T) {
	tests := testCasesGitLab
	assertionLoop(t, tests)
}

func TestCLI(t *testing.T) {
	tests := [][]cliTestCase{
		analyzeFlagTests,
	}

	for _, cliTestCases := range tests {
		cliTestLoop(t, cliTestCases)
	}
}

func cliTestLoop(t *testing.T, cliTests []cliTestCase) {
	jq := gojsonq.New(gojsonq.SetSeparator("->")).File(*reportPath)
	for _, cliTest := range cliTests {
		if !strings.Contains(*executionArgs, cliTest.legitifyCommand) {
			continue
		}
		t.Logf("Testing: %s", cliTest.legitifyCommand)
		jq.Reset()
		content := jq.From("content")
		count := content.Where(cliTest.field, cliTest.op, cliTest.value).Count()
		if count != 0 {
			t.Logf("Failed on test %s", cliTest.legitifyCommand)
			t.Fail()
		}
	}
}
