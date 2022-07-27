package cmd

import (
	"fmt"
	"github.com/Legit-Labs/legitify/internal/common/types"
	"strings"
)

func validateRepositories(repositories []string) ([]types.RepositoryWithOwner, error) {
	var result []types.RepositoryWithOwner

	for _, repo := range repositories {
		splitted := strings.Split(repo, "/")
		if len(splitted) != 2 {
			return nil, fmt.Errorf("invalid repository format %s expected \"owner/name\"", repo)
		}

		result = append(result, types.RepositoryWithOwner{
			Owner: splitted[0],
			Name:  splitted[1],
		})
	}

	return result, nil
}
