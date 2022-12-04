package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
)

const defaultSeparator = ","

func main() {
	if os.Getenv("INPUT_TOKEN") == "" {
		fmt.Printf("::error ::Input token is required.")
		return
	}
	repo, owner, err := splitRepo(os.Getenv("GITHUB_REPOSITORY"))
	if err != nil {
		fmt.Printf("::error ::Unable to split repo: %s", err.Error())
		return
	}
	gitTags, err := getGitTags(repo, owner, os.Getenv("INPUT_TOKEN"))
	if err != nil {
		fmt.Printf("::error ::Unable to pull tags: %s", err.Error())
		return
	}
	latestVersion, err := getLatestVersion(gitTags)
	if err != nil {
		fmt.Printf("::error ::No latest version")
		return
	}
	output := getOutput(
		os.Getenv("GITHUB_REPOSITORY"),
		os.Getenv("INPUT_REPOSITORY"),
		os.Getenv("GITHUB_REF"),
		os.Getenv("GITHUB_SHA"),
		os.Getenv("INPUT_REGISTRIES"),
		os.Getenv("INPUT_SEPARATOR"),
		os.Getenv("INPUT_FULLNAME"),
		latestVersion,
	)
	err = AppendToOutputFile(output)
	if err != nil {
		fmt.Printf("::error ::No latest version")
		return
	}
}

func getOutput(gitRepo, inputRepo, gitRef, gitSHA, inputRegistries, separator, fullName string, latestVersion *semver.Version) map[string]string {
	imageName := getImageName(gitRepo, inputRepo)
	registries := parseRegistriesInput(inputRegistries)
	version := refToVersion(gitRef, gitSHA)
	versions := refToVersions(gitRef, latestVersion)
	tags := getTags(imageName, registries, versions, getFullName(fullName))
	separator = getSeparator(separator)
	return map[string]string{
		"tags":    strings.Join(tags, separator),
		"version": version,
	}
}

func AppendToOutputFile(output map[string]string) error {
	outFile := os.Getenv("GITHUB_OUTPUT")
	f, err := os.OpenFile(outFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	for key, value := range output {
		_, _ = f.WriteString(key)
		_, _ = f.WriteString("=")
		_, _ = f.WriteString(value)
		_, _ = f.WriteString("\n")
	}
	return nil
}
