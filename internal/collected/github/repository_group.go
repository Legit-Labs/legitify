package githubcollected

import (
	"fmt"
	"github.com/google/go-github/v44/github"
)

type RunnerGroup struct {
	Organization ExtendedOrg         `json:"organization"`
	RunnerGroup  *github.RunnerGroup `json:"runner_group"`
}

func (o RunnerGroup) ViolationEntityType() string {
	return "runner group"
}

func (o RunnerGroup) CanonicalLink() string {
	const linkTemplate = "https://github.com/organizations/%s/settings/actions/runner-groups/%d"
	return fmt.Sprintf(linkTemplate, *o.Organization.Login, *o.RunnerGroup.ID)
}

func (o RunnerGroup) Name() string {
	return *o.RunnerGroup.Name
}

func (o RunnerGroup) ID() int64 {
	return *o.RunnerGroup.ID
}
