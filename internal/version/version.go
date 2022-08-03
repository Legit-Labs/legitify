package version

import "fmt"

var (
	/* Note: name, version and commit are injected during CI. */
	Name    = "legitify"
	Version = "na"
	Commit  = "na"
)

var ReadableVersion = fmt.Sprintf("%s version %s commit %s", Name, Version, Commit)
