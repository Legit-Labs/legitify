package gitlab

import (
	"fmt"
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
		res.TotalEntities = len(repositories)
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
				res.TotalEntities += resp.TotalItems
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
				project, _, err := rc.Client.Client().Projects.GetProject(getRepositoryEncodedName(repo), &gitlab2.GetProjectOptions{})
				if err != nil {
					log.Println(err.Error())
					return
				}

				rc.extendedCollection(project)
			})

		}
		gw.Wait()
	})
}

func getRepositoryEncodedName(repo types.RepositoryWithOwner) string {
	return repo.Owner + "/" + repo.Name
}

func (rc *repositoryCollector) extendProjectWithProtectedBranches(project gitlab_collected.Repository) (gitlab_collected.Repository, error) {
	var completeProtectedBranches []*gitlab2.ProtectedBranch
	options := gitlab2.ListProtectedBranchesOptions{}

	err := gitlab.PaginateResults(func(opts *gitlab2.ListOptions) (*gitlab2.Response, error) {
		projectProtectedBranches, resp, err := rc.Client.Client().ProtectedBranches.ListProtectedBranches(int(project.ID()), &options)
		if err != nil {
			return nil, err
		}
		completeProtectedBranches = append(completeProtectedBranches, projectProtectedBranches...)

		return resp, nil
	}, (*gitlab2.ListOptions)(&options))
	if err != nil {
		log.Printf("failed to list projects %s", err)
		return project, err
	}

	extendedProject := project
	extendedProject.ProtectedBranches = completeProtectedBranches
	return extendedProject, nil
}

func (rc *repositoryCollector) extendProjectWithMembers(project gitlab_collected.Repository) (gitlab_collected.Repository, error) {
	var completeMembersList []*gitlab2.ProjectMember
	options := &gitlab2.ListProjectMembersOptions{}

	err := gitlab.PaginateResults(func(opts *gitlab2.ListOptions) (*gitlab2.Response, error) {
		projectMembers, resp, err := rc.Client.Client().ProjectMembers.ListAllProjectMembers(int(project.ID()), options)
		if err != nil {
			return nil, err
		}
		completeMembersList = append(completeMembersList, projectMembers...)
		return resp, nil
	}, &options.ListOptions)

	if err != nil {
		log.Printf("failed to list projects %s", err)
		return project, err
	}

	extendedProject := project
	extendedProject.Members = completeMembersList
	return extendedProject, nil
}

func (rc *repositoryCollector) extendProjectWithWebhooks(project gitlab_collected.Repository) (gitlab_collected.Repository, error) {
	var completeProjectWebhookList []*gitlab2.ProjectHook
	options := gitlab2.ListProjectHooksOptions{}

	err := gitlab.PaginateResults(func(opts *gitlab2.ListOptions) (*gitlab2.Response, error) {
		projectWebhooks, resp, err := rc.Client.Client().Projects.ListProjectHooks(int(project.ID()), &options)
		if err != nil {
			return nil, err
		}
		completeProjectWebhookList = append(completeProjectWebhookList, projectWebhooks...)
		return resp, nil
	}, (*gitlab2.ListOptions)(&options))
	if err != nil {
		log.Printf("failed to list project: %s webhook. error message: %s", project.Name(), err)
		return project, err
	}

	extendedProject := project
	extendedProject.Webhooks = completeProjectWebhookList
	return extendedProject, nil
}

func (rc *repositoryCollector) extendProjectWithPushRules(project gitlab_collected.Repository) (gitlab_collected.Repository, error) {
	rules, _, err := rc.Client.Client().Projects.GetProjectPushRules(int(project.ID()))
	if err != nil {
		log.Printf("failed to get project push rule %s", err)
		return project, err
	}
	extendedProject := project
	extendedProject.PushRules = rules
	return extendedProject, nil
}

func (rc *repositoryCollector) extendProjectWithMergeRequestApprovalRules(project gitlab_collected.Repository) (gitlab_collected.Repository, error) {
	rules, _, err := rc.Client.Client().Projects.GetProjectApprovalRules(int(project.ID()))
	if err != nil {
		log.Printf("failed to get project merge request approval rules %s", err)
		return project, err
	}
	extendedProject := project
	extendedProject.ApprovalRules = rules
	return extendedProject, nil
}

func (rc *repositoryCollector) extendProjectWithApprovalConfiguration(project gitlab_collected.Repository) (gitlab_collected.Repository, error) {
	config, _, err := rc.Client.Client().Projects.GetApprovalConfiguration(int(project.ID()))
	if err != nil {
		log.Printf("failed to get project approval configuration %s", err)
		return project, err
	}
	extendedProject := project
	extendedProject.ApprovalConfiguration = config
	return extendedProject, nil
}

func (rc *repositoryCollector) extendProjectWithMinimumRequiredApprovals(project gitlab_collected.Repository) (gitlab_collected.Repository, error) {
	minRequiredApprovals := 0

	for _, rule := range project.ApprovalRules {
		if rule.ApprovalsRequired > minRequiredApprovals {
			minRequiredApprovals = rule.ApprovalsRequired
		}
	}

	project.MinimumRequiredApprovals = minRequiredApprovals
	return project, nil
}

func (rc *repositoryCollector) collectAll() collectors.SubCollectorChannels {
	return rc.WrappedCollection(func() {
		options := gitlab2.ListProjectsOptions{}
		organizations, err := rc.Client.Organizations()
		if err != nil {
			log.Printf("failed to collect list of orgniazations to get repositories  %s", err)
			return
		}
		gw := group_waiter.New()
		for _, org := range organizations {
			err := gitlab.PaginateResults(func(opts *gitlab2.ListOptions) (*gitlab2.Response, error) {
				repos, resp, err := rc.Client.Client().Groups.ListGroupProjects(org.Name, &gitlab2.ListGroupProjectsOptions{})
				if err != nil {
					return nil, err
				}
				for _, completeProject := range repos {
					k := completeProject
					gw.Do(func() {
						rc.extendedCollection(k)
					})
				}
				return resp, nil
			}, &options.ListOptions)
			if err != nil {
				log.Printf("failed to list projects %s", err)
			}
		}
		gw.Wait()
	})
}

func (rc *repositoryCollector) extendedCollection(completeProjectsList *gitlab2.Project) {
	proj := gitlab_collected.Repository{
		Project: completeProjectsList,
	}
	extensionFunctions := []func(gitlab_collected.Repository) (gitlab_collected.Repository, error){
		rc.extendProjectWithMembers,
		rc.extendProjectWithProtectedBranches,
		rc.extendProjectWithWebhooks,
		rc.extendProjectWithPushRules,
		rc.extendProjectWithMergeRequestApprovalRules,
		rc.extendProjectWithApprovalConfiguration,
		rc.extendProjectWithMinimumRequiredApprovals,
	}
	var err error
	for _, f := range extensionFunctions {
		proj, err = f(proj)
		if err != nil {
			fmt.Printf("Project '%s' collection failed with error:%s", proj.Name(), err)
			break
		}
	}

	if err == nil {
		newContext := newCollectionContext(nil, []permissions.OrganizationRole{permissions.OrgRoleOwner})
		rc.CollectDataWithContext(proj, proj.Links.Self, &newContext)
	}
	rc.CollectionChangeByOne()
}
