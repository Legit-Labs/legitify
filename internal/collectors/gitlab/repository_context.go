package gitlab

import (
	"github.com/Legit-Labs/legitify/internal/common/permissions"
)

type repositoryContext struct {
	testParam    bool
	isEnterprise bool
	roles        []permissions.Role
}

func (rc *repositoryContext) Premium() bool {
	return rc.isEnterprise
}

func (rc *repositoryContext) Roles() []permissions.Role {
	return rc.roles
}
