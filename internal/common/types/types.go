package types

import (
	"fmt"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
)

type RepositoryWithOwner struct {
	Name  string
	Owner string
}

func (r RepositoryWithOwner) String() string {
	return fmt.Sprintf("%s/%s", r.Owner, r.Name)
}

type Organization struct {
	Name string
	Role permissions.OrganizationRole
}
