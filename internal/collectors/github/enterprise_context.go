package github

import "github.com/Legit-Labs/legitify/internal/common/permissions"

type enterpriseContext struct {
	roles []permissions.Role
}

func (ec *enterpriseContext) Premium() bool {
	return true
}

func (ec *enterpriseContext) Roles() []permissions.Role {
	return ec.roles
}
func newEnterpriseContext(roles []permissions.Role) *enterpriseContext {
	return &enterpriseContext{
		roles: roles,
	}
}
