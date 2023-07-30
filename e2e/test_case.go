package test

type testCase struct {
	path          string
	failedEntity  string
	passedEntity  string
	skippedEntity string
}

type cliTestCase struct {
	legitifyCommand string
	field           string
	value           string
}
