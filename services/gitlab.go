package services

import (
	"caco/settings"
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/xanzy/go-gitlab"
)

const PER_PAGE = 999999999

type GitlabService struct {
	client *gitlab.Client
}

// NewGitlabService ...
func NewGitlabService() *GitlabService {
	git, err := gitlab.NewClient(settings.Config.GitlabToken, gitlab.WithBaseURL(settings.Config.GitlabURL))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return &GitlabService{
		client: git,
	}
}

func (s *GitlabService) GetGroupByName(team string) (*gitlab.Group, error) {
	page := 1
	groups, response, err := s.callListGroups(page)
	if err != nil {
		return nil, err
	}

	for hasPagination(response) {
		page++
		newGroups, newResponse, err := s.callListGroups(page)
		if err != nil {
			return nil, err
		}
		groups = append(groups, newGroups...)
		response = newResponse
	}

	// find if we find the group in the array
	group := findGroup(groups, team)

	return group, nil
}

func (s *GitlabService) GetGroupProjects(gid int) ([]*gitlab.Project, error) {
	page := 1
	projects, response, err := s.callListGroupProjects(gid, page)

	if err != nil {
		return nil, err
	}

	for hasPagination(response) {
		page++
		newProjects, newResponse, err := s.callListGroupProjects(gid, page)
		if err != nil {
			return nil, err
		}
		projects = append(projects, newProjects...)
		response = newResponse
	}

	// filter exclusion list
	filteredProjects := filterProjects(projects, func(i int) bool {
		_, ok := settings.Config.GitlabProjectExclusionsMap[i]
		return !ok //invertemos o bool pois só queremos o que não está na lista
	})

	return filteredProjects, nil
}

func (s *GitlabService) GetProjectOpenedMergeRequests(pid int) ([]*gitlab.MergeRequest, error) {
	// get group projects
	mrs, _, err := s.client.MergeRequests.ListProjectMergeRequests(
		pid,
		&gitlab.ListProjectMergeRequestsOptions{

			State: gitlab.String("opened"),
			Scope: gitlab.String("all"),
			ListOptions: gitlab.ListOptions{
				PerPage: PER_PAGE,
			},
		},
		gitlab.WithContext(context.Background()),
	)

	if err != nil {
		return nil, err
	}

	return mrs, nil
}

func findGroup(groups []*gitlab.Group, name string) *gitlab.Group {
	for _, g := range groups {
		if settings.Config.Debug {
			fmt.Println("Group: id: " + strconv.Itoa(g.ID) + ", name: " + g.Name)
		}
		if name == g.Name {
			return g
		}
	}
	return nil
}

func filterProjects(source []*gitlab.Project, predicate func(int) bool) []*gitlab.Project {
	target := make([]*gitlab.Project, 0)
	for _, p := range source {
		if predicate(p.ID) {
			target = append(target, p)
		}
	}
	return target
}

func hasPagination(response *gitlab.Response) bool {
	return response.TotalPages > response.CurrentPage
}

func (s *GitlabService) callListGroups(page int) ([]*gitlab.Group, *gitlab.Response, error) {
	return s.client.Groups.ListGroups(
		&gitlab.ListGroupsOptions{
			AllAvailable: gitlab.Bool(true),
			ListOptions: gitlab.ListOptions{
				PerPage: PER_PAGE,
				Page:    page,
			},
		},
		gitlab.WithContext(context.Background()),
	)
}

func (s *GitlabService) callListGroupProjects(gid int, page int) ([]*gitlab.Project, *gitlab.Response, error) {
	return s.client.Groups.ListGroupProjects(gid,
		&gitlab.ListGroupProjectsOptions{
			Archived: gitlab.Bool(false),
			Simple:   gitlab.Bool(true),
			ListOptions: gitlab.ListOptions{
				PerPage: PER_PAGE,
				Page:    page,
			},
		},
		gitlab.WithContext(context.Background()),
	)
}
