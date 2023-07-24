package gitlab

import (
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	"github.com/xanzy/go-gitlab"
)

type collectionContext struct {
	group     *gitlab.Group
	roles     []permissions.Role
	isPremium bool
}

func newCollectionContext(group *gitlab.Group, roles []permissions.Role, isPremium bool) collectionContext {
	return collectionContext{
		group:     group,
		roles:     roles,
		isPremium: isPremium,
	}
}

func (c collectionContext) Premium() bool {
	return c.isPremium
}

func (c collectionContext) Roles() []permissions.Role {
	return c.roles
}
