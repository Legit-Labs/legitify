package policies

import (
	"embed"
)

//go:embed github/*
var GitHubBundle embed.FS

//go:embed gitlab/*
var GitlabBundle embed.FS
