package gitlab

import (
	"github.com/Legit-Labs/legitify/internal/clients/gitlab"
	"github.com/Legit-Labs/legitify/internal/collected/gitlab_collected"
	"github.com/Legit-Labs/legitify/internal/collectors"
	"github.com/Legit-Labs/legitify/internal/common/group_waiter"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	"github.com/Legit-Labs/legitify/internal/common/types"
	"github.com/Legit-Labs/legitify/internal/context_utils"

	gitlab2 "github.com/xanzy/go-gitlab"
	"log"

	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"golang.org/x/net/context"
)

type repositoryCollector struct {
	collectors.BaseCollector
	Client  *gitlab.Client
	Context context.Context
}

func NewRepositoryCollector(ctx context.Context, client *gitlab.Client) collectors.Collector {
	c := &repositoryCollector{
		Client:  client,
		Context: ctx,
	}
	collectors.InitBaseCollector(&c.BaseCollector, c)
	return c
}

func (rc *repositoryCollector) Namespace() namespace.Namespace {
	return namespace.Repository
}

func (rc *repositoryCollector) CollectMetadata() collectors.Metadata {
	res := collectors.Metadata{}

	repositories, exist := context_utils.GetRepositories(rc.Context)

	if exist {
		for _, repository := range repositories {
			_, _, err := rc.Client.Client().Projects.GetProject(repository, &gitlab2.GetProjectOptions{})
			if err != nil {
				log.Printf("failed to collect metadata for repository %s", err)
			} else {
				res.TotalEntities++
			}
		}
	} else {
		organizations, err := rc.Client.Organizations()
		if err != nil {
			log.Printf("failed to collect list of orgniazations to get repositories metadata %s", err)
			return res
		}
		for _, org := range organizations {
			_, resp, err := rc.Client.Client().Groups.ListGroupProjects(org.Name, &gitlab2.ListGroupProjectsOptions{})

			if err != nil {
				log.Printf("failed to collect metadata for repositories %s", err)
			} else {
				res.TotalEntities = resp.TotalItems
			}
		}
	}
	return res
}

func (rc *repositoryCollector) Collect() collectors.SubCollectorChannels {
	repositories, exist := context_utils.GetRepositories(rc.Context)

	if exist {
		return rc.collectSpecific(repositories)
	}

	return rc.collectAll()
}

func (rc *repositoryCollector) collectSpecific(repositories []types.RepositoryWithOwner) collectors.SubCollectorChannels {
	return rc.WrappedCollection(func() {
		gw := group_waiter.New()
		for _, r := range repositories {
			repo := r
			gw.Do(func() {
				project, _, err := rc.Client.Client().Projects.GetProject(repo.Owner+"%2F"+repo.Name, &gitlab2.GetProjectOptions{})
				if err != nil {
					log.Println(err.Error())
					return
				}

				rc.extendedCollect(project)
			})

		}
		gw.Wait()
	})
}

func (rc *repositoryCollector) extendProjectWithProtectedBranches(project gitlab_collected.Repository) gitlab_collected.Repository {
	var completeProtectedBranches []*gitlab2.ProtectedBranch
	options := gitlab2.ListProtectedBranchesOptions{}

	err := gitlab.PaginateResults(func(opts *gitlab2.ListOptions) (*gitlab2.Response, error) {
		projectProtectedBranches, resp, err := rc.Client.Client().ProtectedBranches.ListProtectedBranches(project.ID, &options)
		if err != nil {
			return nil, err
		}
		for _, protectedBranch := range projectProtectedBranches {
			completeProtectedBranches = append(completeProtectedBranches, protectedBranch)
		}

		return resp, nil
	}, (*gitlab2.ListOptions)(&options))
	if err != nil {
		log.Printf("failed to list projects %s", err)
	}

	extendedProject := project
	extendedProject.ProtectedBranches = completeProtectedBranches
	return extendedProject
}

func (rc *repositoryCollector) extendProjectWithMembers(project gitlab_collected.Repository) gitlab_collected.Repository {
	var completeMembersList []*gitlab2.ProjectMember
	options := &gitlab2.ListProjectMembersOptions{}

	err := gitlab.PaginateResults(func(opts *gitlab2.ListOptions) (*gitlab2.Response, error) {
		projectMembers, resp, err := rc.Client.Client().ProjectMembers.ListAllProjectMembers(project.ID, options)
		if err != nil {
			return nil, err
		}
		for _, projectMember := range projectMembers {
			completeMembersList = append(completeMembersList, projectMember)
		}

		return resp, nil
	}, &options.ListOptions)

	if err != nil {
		log.Printf("failed to list projects %s", err)
	}

	extendedProject := project
	extendedProject.Members = completeMembersList
	return extendedProject
}

func (rc *repositoryCollector) collectAll() collectors.SubCollectorChannels {
	return rc.WrappedCollection(func() {

		var completeProjectsList []*gitlab2.Project
		maintainerPermissions := gitlab2.MaintainerPermissions
		options := gitlab2.ListProjectsOptions{MinAccessLevel: &maintainerPermissions}

		organizations, err := rc.Client.Organizations()
		if err != nil {
			log.Printf("failed to collect list of orgniazations to get repositories  %s", err)
			return
		}

		for _, org := range organizations {
			err := gitlab.PaginateResults(func(opts *gitlab2.ListOptions) (*gitlab2.Response, error) {
				repos, resp, err := rc.Client.Client().Groups.ListGroupProjects(org.Name, &gitlab2.ListGroupProjectsOptions{})
				if err != nil {
					return nil, err
				}
				for _, r := range repos {
					completeProjectsList = append(completeProjectsList, r)
				}
				return resp, nil
			}, &options.ListOptions)
			if err != nil {
				log.Printf("failed to list projects %s", err)
			}
		}

		for _, completeProject := range completeProjectsList {
			rc.extendedCollect(completeProject)

		}
	})
}

func (rc *repositoryCollector) extendedCollect(completeProjectsList *gitlab2.Project) {
	proj := gitlab_collected.Repository{
		Project: completeProjectsList,
	}
	extendedProject := rc.extendProjectWithMembers(proj)

	extendedProject = rc.extendProjectWithProtectedBranches(extendedProject)

	newContext := newCollectionContext(nil, []permissions.OrganizationRole{permissions.OrgRoleOwner})
	rc.CollectDataWithContext(extendedProject, extendedProject.Links.Self, &newContext)
}
