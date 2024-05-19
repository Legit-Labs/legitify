package policies

import (
	"embed"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func countFiles(fs embed.FS, p string) (int, error) {
	f, e := fs.ReadDir(p)
	totalCount := 0

	if e != nil {
		return 0, e
	}

	for _, de := range f {
		if de.IsDir() {
			count, e := countFiles(fs, path.Join(p, de.Name()))

			if e != nil {
				return 0, e
			}

			totalCount += count
		} else {
			totalCount += 1
		}
	}

	return totalCount, nil
}

func countBundles() (int, error) {
	count, e := countFiles(GitHubBundle, path.Dir(""))

	if e != nil {
		return 0, e
	}

	return count, nil
}

func TestPoliciesBundle(t *testing.T) {
	count, err := countBundles()

	require.Nilf(t, err, "counting files: %v", err)
	require.Equal(t, count, 9, "Expecting 9 files in bundle")
}
