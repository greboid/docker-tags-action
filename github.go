package main

import (
	"context"

	"github.com/google/go-github/v38/github"
	"golang.org/x/oauth2"
)

func getGitTags(owner, repo string, token string) ([]string, error) {
	client := github.NewClient(oauth2.NewClient(context.Background(),
		oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)))
	opt := &github.ListOptions{
		PerPage: 10,
	}
	var allTags []*github.RepositoryTag
	for {
		tags, resp, err := client.Repositories.ListTags(context.Background(), owner, repo, opt)
		if err != nil {
			return []string{}, err
		}
		allTags = append(allTags, tags...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	var allTagsString []string
	for index := range allTags {
		allTagsString = append(allTagsString, *allTags[index].Name)
	}
	return allTagsString, nil
}
