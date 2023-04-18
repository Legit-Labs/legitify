package githubcollected

import (
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
)

type Enterprise struct {
	MembersCanChangeRepositoryVisibilitySetting string `json:"visibility_change_disabled"`
	EnterpriseName                              string `json:"name"`
	Url                                         string `json:"url"`
	Id                                          int64  `json:"id"`
	UserRole                                    string
}

func NewEnterprise(MembersCanChangeRepositoryVisibilitySetting string, Name string, Url string, Id int64, admin bool) Enterprise {
	UserRole := ""
	if admin {
		UserRole = permissions.EnterpriseAdmin
	}
	return Enterprise{MembersCanChangeRepositoryVisibilitySetting: MembersCanChangeRepositoryVisibilitySetting,
		EnterpriseName: Name,
		Url:            Url,
		Id:             Id,
		UserRole:       UserRole,
	}
}

func (o Enterprise) ViolationEntityType() string {
	return namespace.Enterprise
}

func (o Enterprise) CanonicalLink() string {
	return o.Url
}

func (o Enterprise) Name() string {
	return o.EnterpriseName
}

func (o Enterprise) ID() int64 {
	return o.Id
}
