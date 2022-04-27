package services

import (
	"caco/settings"
	"context"

	"github.com/google/go-github/v43/github"
	"golang.org/x/oauth2"
)

const GITHUB_PER_PAGE = 100

type GithubService struct {
	client *github.Client
}

func NewGithubService() *GithubService {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: settings.Config.GithubPersonalToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	git := github.NewClient(tc)

	return &GithubService{
		client: git,
	}
}

func (svc *GithubService) GetOpenPullRequests(project string) ([]*github.PullRequest, error) {
	options := github.PullRequestListOptions{
		ListOptions: github.ListOptions{
			PerPage: GITHUB_PER_PAGE,
		},
	}
	prs, _, err := svc.client.PullRequests.List(context.Background(), settings.Config.GithubOrg, project, &options)
	if err != nil {
		println(err)
	}

	return prs, err
}

func (svc *GithubService) GetReviewsFromPR(project string, prid int) {
	options := &github.ListOptions{
		PerPage: GITHUB_PER_PAGE,
	}
	svc.client.PullRequests.ListReviews(context.Background(), settings.Config.GithubOrg, project, prid, options)
}
