package policies

import (
	"embed"
)

//go:embed github/*
var Bundle embed.FS
