package gitlab

import (
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	"github.com/xanzy/go-gitlab"
)

type collectionContext struct {
	group *gitlab.Group
	roles []permissions.Role
}

func newCollectionContext(group *gitlab.Group, roles []permissions.Role) collectionContext {
	return collectionContext{
		group: group,
		roles: roles,
	}
}

func (c collectionContext) Premium() bool {
	// TODO: currently we don't know about GitLab policies that required premium account, need to update this once we do.
	return true
}

func (c collectionContext) Roles() []permissions.Role {
	return c.roles
}
