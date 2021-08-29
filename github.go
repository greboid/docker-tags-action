package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func getOutput(gitRepo, inputRepo, gitRef, gitSHA, inputRegistries, separator, fullName string, latestVersion *semver.Version) string {
	imageName := getImageName(gitRepo, inputRepo)
	registries := parseRegistriesInput(inputRegistries)
	version := refToVersion(gitRef, gitSHA)
	versions := refToVersions(gitRef, latestVersion)
	tags := getTags(imageName, registries, versions, getFullName(fullName))
	separator = getSeparator(separator)
	return fmt.Sprintf("::set-output name=tags::%s\n::set-output name=version::%s", strings.Join(tags, separator), version)
}

func getGitTags(directory string) ([]string, error) {
	repo, err := git.PlainOpen(directory)
	if err != nil {
		return nil, err
	}
	err = repo.Fetch(&git.FetchOptions{
		Tags: git.AllTags,
		Force: true,
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return nil, err
	}
	iter, err := repo.Tags()
	if err != nil {
		return nil, err
	}
	tags := make([]string, 0)
	err = iter.ForEach(func(reference *plumbing.Reference) error {
		tags = append(tags, reference.Name().Short())
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tags, nil
}