package namespace

import "fmt"

type Namespace = string

const (
	Enterprise   Namespace = "enterprise"
	Organization Namespace = "organization"
	Repository   Namespace = "repository"
	Member       Namespace = "member"
	Actions      Namespace = "actions"
	RunnerGroup  Namespace = "runner_group"
)

var All = []Namespace{
	Organization,
	Enterprise,
	Repository,
	Member,
	Actions,
	RunnerGroup,
}

func ValidateNamespaces(namespace []Namespace) error {
	for _, ns := range namespace {
		found := false
		for _, e := range All {
			if e == ns {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("invalid namespace %s", ns)
		}
	}

	return nil
}
