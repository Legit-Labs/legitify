package gitlab_collected

import (
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/xanzy/go-gitlab"
)

type Server struct {
	url string
	*gitlab.Settings
}

func NewServer(url string, settings *gitlab.Settings) *Server {
	return &Server{
		url:      url,
		Settings: settings,
	}
}

func (o Server) ViolationEntityType() string {
	return namespace.Server
}

func (o Server) CanonicalLink() string {
	return o.url
}

func (o Server) Name() string {
	return o.url
}

func (o Server) ID() int64 {
	return int64(o.Settings.ID)
}
