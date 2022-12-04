package gitlab_collected

import (
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/xanzy/go-gitlab"
)

type Organization struct {
	*gitlab.Group
	Hooks []*gitlab.GroupHook `json:"hooks"`
}

func (o Organization) ViolationEntityType() string {
	return namespace.Organization
}

func (o Organization) CanonicalLink() string {
	return o.WebURL
}

func (o Organization) Name() string {
	return o.FullName
}

func (o Organization) ID() int64 {
	return int64(o.Group.ID)
}
