package gitlab

import (
	"github.com/Legit-Labs/legitify/internal/common/permissions"
)

type serverCollectionContext struct {
}

func newServerCollectionContext() serverCollectionContext {
	return serverCollectionContext{}
}

func (c serverCollectionContext) Premium() bool {
	return true
}

func (c serverCollectionContext) Roles() []permissions.Role {
	return []permissions.Role{permissions.OrgRoleOwner}
}
