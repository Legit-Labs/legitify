package test

var analyzeFlagTests = []cliTestCase{
	{
		legitifyCommand: "--failed-only",
		field:           "status",
		value:           "FAILED",
	},
	{
		legitifyCommand: "--namespace repository",
		field:           "violationEntityType",
		value:           "repository",
	},
	{
		legitifyCommand: "--org Legitify-E2E",
		field:           "canonicalLink",
		value:           "https://github.com/Legitify-E2E/",
	},
	{
		legitifyCommand: "--repo Legitify-E2E/bad_branch_protection",
		field:           "canonicalLink",
		value:           "https://github.com/Legitify-E2E/bad_branch_protection",
	},
}
