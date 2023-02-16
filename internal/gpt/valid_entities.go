package gpt

import (
	_ "embed"
)

//go:embed valid_entities/gh_org.json
var validGitHubOrganization []byte

//go:embed valid_entities/gh_repo.json
var validGitHubRepository []byte
