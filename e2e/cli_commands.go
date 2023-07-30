package test

var analyzeFlagTests = []cliTestCase{
	{
		legitifyCommand: "--failed-only",
		field:           "status",
		value:           "FAILED",
		op:              "=",
	},
	{
		legitifyCommand: "--namespace repository",
		field:           "violationEntityType",
		value:           "repository",
		op:              "=",
	},
	{
		legitifyCommand: "--org Legitify-E2E",
		field:           "canonicalLink",
		value:           "Legitify-E2E",
		op:              "contains",
	},
	{
		legitifyCommand: "--repo Legitify-E2E/bad_branch_protection",
		field:           "canonicalLink",
		value:           "https://github.com/Legitify-E2E/bad_branch_protection",
		op:              "=",
	},
}
