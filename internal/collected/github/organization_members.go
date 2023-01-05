package githubcollected

import (
	"fmt"

	"github.com/google/go-github/v49/github"
)

type OrganizationMember struct {
	User       *github.User `json:"user"`
	LastActive int          `json:"last_active"`
	IsAdmin    bool         `json:"is_admin"`
}

type OrganizationMembers struct {
	Organization  ExtendedOrg          `json:"organization"`
	Members       []OrganizationMember `json:"members"`
	HasLastActive bool                 `json:"has_last_active"`
}

func NewOrganizationMember(user *github.User, lastActive int, memberType string) OrganizationMember {
	return OrganizationMember{
		User:       user,
		LastActive: lastActive,
		IsAdmin:    memberType == "admin",
	}
}

func (o OrganizationMembers) ViolationEntityType() string {
	return "organization members"
}

func (o OrganizationMembers) CanonicalLink() string {
	return fmt.Sprintf("https://github.com/orgs/%s/people", o.Name())
}

func (o OrganizationMembers) Name() string {
	// Deliberately using the Org; see membersList enricher
	return *o.Organization.Login
}

func (o OrganizationMembers) ID() int64 {
	// Deliberately using the Org; see membersList enricher
	return *o.Organization.ID
}
