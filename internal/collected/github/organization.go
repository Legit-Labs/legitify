package githubcollected

import (
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/common/permissions"

	"github.com/google/go-github/v44/github"
)

type ExtendedOrg struct {
	github.Organization
	Role permissions.OrganizationRole
}

func NewExtendedOrg(org *github.Organization, role permissions.OrganizationRole) ExtendedOrg {
	return ExtendedOrg{*org, role}
}

func (e ExtendedOrg) IsEnterprise() bool {
	const orgPlanEnterprise = "enterprise"
	if e.Plan == nil {
		return false
	}
	return e.Plan.GetName() == orgPlanEnterprise
}

func (e ExtendedOrg) IsFree() bool {
	const orgPlanFree = "free"
	if e.Plan == nil {
		return false
	}
	return e.Plan.GetName() == orgPlanFree
}

type Organization struct {
	Organization *ExtendedOrg   `json:"organization"`
	SamlEnabled  *bool          `json:"saml_enabled,omitempty"`
	Hooks        []*github.Hook `json:"hooks"`
	UserRole     permissions.OrganizationRole
}

func (o Organization) ViolationEntityType() string {
	return namespace.Organization
}

func (o Organization) CanonicalLink() string {
	return *o.Organization.HTMLURL
}

func (o Organization) Name() string {
	return *o.Organization.Login
}

func (o Organization) ID() int64 {
	return *o.Organization.ID
}
