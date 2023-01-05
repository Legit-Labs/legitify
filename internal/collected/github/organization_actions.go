package githubcollected

import (
	"fmt"
	"github.com/Legit-Labs/legitify/internal/clients/github/types"

	"github.com/google/go-github/v49/github"
)

type OrganizationActions struct {
	Organization       ExtendedOrg                `json:"organization"`
	ActionsPermissions *github.ActionsPermissions `json:"actions_permissions"`
	TokenPermissions   *types.TokenPermissions    `json:"token_permissions"`
}

func (o OrganizationActions) ViolationEntityType() string {
	return "organization actions"
}

func (o OrganizationActions) CanonicalLink() string {
	const linkTemplate = "https://github.com/organizations/%s/settings/actions"
	return fmt.Sprintf(linkTemplate, *o.Organization.Login)
}

func (o OrganizationActions) Name() string {
	return *o.Organization.Login
}

func (o OrganizationActions) ID() int64 {
	return *o.Organization.ID
}
