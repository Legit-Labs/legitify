package types

import (
	"fmt"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	"strings"
)

type RepositoryWithOwner struct {
	Name  string
	Owner string
	Role  permissions.RepositoryRole
}

func NewRepositoryWithOwner(repositoryWithOwner string, perms permissions.RepositoryRole) RepositoryWithOwner {
	owner, name, _ := strings.Cut(repositoryWithOwner, "/")

	return RepositoryWithOwner{
		Owner: owner,
		Name:  name,
		Role:  perms,
	}
}

func (r RepositoryWithOwner) String() string {
	return fmt.Sprintf("%s/%s", r.Owner, r.Name)
}

type Organization struct {
	Name string
	Role permissions.OrganizationRole
}
