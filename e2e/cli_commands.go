package test

var analyzeFlagTests = []cliTestCase{
	{
		legitifyCommand: "--failed-only",
		field:           "status",
		op:              "!=",
		value:           "FAILED",
	},
	{
		legitifyCommand: "--namespace repository",
		field:           "violationEntityType",
		op:              "!=",
		value:           "repository",
	},
	{
		legitifyCommand: "--org Legitify-E2E",
		field:           "canonicalLink",
		op:              "notIn",
		value:           "https://github.com/Legitify-E2E/",
	},
	{
		legitifyCommand: "--repo Legitify-E2E/bad_branch_protection",
		field:           "canonicalLink",
		op:              "!=",
		value:           "https://github.com/Legitify-E2E/bad_branch_protection",
	},
}
