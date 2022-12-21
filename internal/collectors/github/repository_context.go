package github

import (
	"github.com/Legit-Labs/legitify/internal/common/permissions"
)

type repositoryContext struct {
	roles                         []permissions.Role
	isEnterprise                  bool
	isBranchProtectionSupported   bool
	hasBranchProtectionPermission bool
}

func (rc *repositoryContext) Premium() bool {
	return rc.isEnterprise
}

func (rc *repositoryContext) Roles() []permissions.Role {
	return rc.roles
}

func (rc *repositoryContext) IsBranchProtectionSupported() bool {
	return rc.isBranchProtectionSupported
}

func (rc *repositoryContext) SetHasBranchProtectionPermission(value bool) {
	rc.hasBranchProtectionPermission = value
}

func (rc *repositoryContext) HasBranchProtectionPermission() bool {
	return rc.hasBranchProtectionPermission
}

func newRepositoryContext(roles []permissions.RepositoryRole, isBranchProtectionSupported bool, isEnterprise bool, hasBranchProtectionPermission bool) *repositoryContext {
	return &repositoryContext{
		roles:                         roles,
		isEnterprise:                  isEnterprise,
		isBranchProtectionSupported:   isBranchProtectionSupported,
		hasBranchProtectionPermission: hasBranchProtectionPermission,
	}
}
