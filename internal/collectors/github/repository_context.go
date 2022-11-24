package github

import (
	"context"
	"github.com/Legit-Labs/legitify/internal/clients/github"
	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
)

type repositoryContext struct {
	roles                       []permissions.Role
	isEnterprise                bool
	isBranchProtectionSupported bool
}

func (rc *repositoryContext) IsEnterprise() bool {
	return rc.isEnterprise
}

func (rc *repositoryContext) Roles() []permissions.Role {
	return rc.roles
}

func (rc *repositoryContext) IsBranchProtectionSupported() bool {
	return rc.isBranchProtectionSupported
}

type repositoryContextFactory struct {
	ctx    context.Context
	client github.Client
}

func newRepositoryContextFactory(ctx context.Context, client github.Client) *repositoryContextFactory {
	return &repositoryContextFactory{
		ctx:    ctx,
		client: client,
	}
}

func (rcf *repositoryContextFactory) newRepositoryContextForOrganization(login string, viewerCanAdminister *bool, repository *githubcollected.GitHubQLRepository) (*repositoryContext, error) {
	org, _, err := rcf.client.Client().Organizations.Get(rcf.ctx, login)
	if err != nil {
		return nil, err
	}
	role := permissions.GetOrgRole(viewerCanAdminister)
	extendedOrg := githubcollected.NewExtendedOrg(org, role)
	return rcf.newRepositoryContextForExtendedOrg(&extendedOrg, repository), nil
}

func (rcf *repositoryContextFactory) newRepositoryContextForExtendedOrg(org *githubcollected.ExtendedOrg, repository *githubcollected.GitHubQLRepository) *repositoryContext {
	return &repositoryContext{
		roles:                       []permissions.Role{org.Role, repository.ViewerPermission},
		isEnterprise:                org.IsEnterprise(),
		isBranchProtectionSupported: org.IsEnterprise() || !repository.IsPrivate,
	}
}

func (rcf *repositoryContextFactory) newRepositoryContextForUser(login string, repository *githubcollected.GitHubQLRepository) (*repositoryContext, error) {
	user, _, err := rcf.client.Client().Users.Get(rcf.ctx, login)
	if err != nil {
		return nil, err
	}

	return &repositoryContext{
		roles:                       []permissions.Role{repository.ViewerPermission},
		isEnterprise:                false,
		isBranchProtectionSupported: !repository.IsPrivate || (user.Plan != nil && *user.Plan.Name != "free"),
	}, nil
}
