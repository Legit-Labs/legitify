package github

import "github.com/Legit-Labs/legitify/internal/common/permissions"

type enterpriseContext struct {
	roles        []permissions.Role
	isEnterprise bool
}

func (ec *enterpriseContext) Premium() bool {
	return ec.isEnterprise
}

func (ec *enterpriseContext) Roles() []permissions.Role {
	return ec.roles
}
func newEnterpriseContext(roles []permissions.Role) *enterpriseContext {
	return &enterpriseContext{
		roles:        roles,
		isEnterprise: true,
	}
}
