package gitlab

import (
	"context"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	"github.com/Legit-Labs/legitify/internal/common/types"
	"github.com/xanzy/go-gitlab"
)

type Client struct {
	context context.Context
	client  *gitlab.Client
}

func (c *Client) Client() *gitlab.Client {
	return c.client
}

func NewClient(ctx context.Context, token string, endpoint string, fillCache bool) (*Client, error) {
	var config []gitlab.ClientOptionFunc
	if endpoint != "" {
		config = []gitlab.ClientOptionFunc{gitlab.WithBaseURL(endpoint)}
	}

	git, err := gitlab.NewClient(token, config...)
	if err != nil {
		return nil, err
	}

	return &Client{
		context: ctx,
		client:  git,
	}, nil
}

func (c *Client) IsAnalyzable(ctx context.Context, repo types.RepositoryWithOwner) (bool, error) {
	return false, nil
}

func (c *Client) Scopes() permissions.TokenScopes {
	return nil
}
