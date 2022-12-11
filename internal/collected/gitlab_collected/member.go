package gitlab_collected

import (
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/xanzy/go-gitlab"
)

type Member struct {
	*gitlab.User
}

func (o Member) ViolationEntityType() string {
	return namespace.Member
}

func (o Member) CanonicalLink() string {
	return o.WebURL
}

func (o Member) Name() string {
	return o.Username
}

func (o Member) ID() int64 {
	return int64(o.User.ID)
}
