package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
)

const defaultSeparator = ","

func main() {
	repo, owner, err := splitRepo(os.Getenv("GITHUB_REPOSITORY"))
	if err != nil {
		fmt.Printf("::error ::Unable to split repo: %s", err.Error())
		return
	}
	gitTags, err := getGitTags(repo, owner, os.Getenv("secrets.GITHUB_TOKEN"))
	if err != nil {
		fmt.Printf("::error ::Unable to pull tags: %s", err.Error())
		return
	}
	latestVersion, err := getLatestVersion(gitTags)
	if err != nil {
		fmt.Printf("::error ::No latest version")
		return
	}
	fmt.Printf(
		getOutput(
			os.Getenv("GITHUB_REPOSITORY"),
			os.Getenv("INPUT_REPOSITORY"),
			os.Getenv("GITHUB_REF"),
			os.Getenv("GITHUB_SHA"),
			os.Getenv("INPUT_REGISTRIES"),
			os.Getenv("INPUT_SEPARATOR"),
			os.Getenv("INPUT_FULLNAME"),
			latestVersion,
		),
	)
}

func getOutput(gitRepo, inputRepo, gitRef, gitSHA, inputRegistries, separator, fullName string, latestVersion *semver.Version) string {
	imageName := getImageName(gitRepo, inputRepo)
	registries := parseRegistriesInput(inputRegistries)
	version := refToVersion(gitRef, gitSHA)
	versions := refToVersions(gitRef, latestVersion)
	tags := getTags(imageName, registries, versions, getFullName(fullName))
	separator = getSeparator(separator)
	return fmt.Sprintf("::set-output name=tags::%s\n::set-output name=version::%s", strings.Join(tags, separator), version)
}
