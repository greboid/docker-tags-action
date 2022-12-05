package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/blang/semver/v4"
)

const defaultSeparator = ","

func main() {
	fmt.Printf("Repository: '%s'\n", os.Getenv("GITHUB_REPOSITORY"))
	fmt.Printf("Ref: '%s'\n", os.Getenv("GITHUB_REF"))
	fmt.Printf("SHA: '%s'\n", os.Getenv("GITHUB_SHA"))

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
	latestVersion := getLatestVersion(gitTags)
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
		log.Printf("Unable to save to file: %s", err)
		fmt.Printf("::error ::Unable to save to file")
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
	if outFile != "" {
		f, err := os.OpenFile(outFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Unable to open output file: %s", err)
			return err
		}
		content := strings.Builder{}
		defer f.Close()
		for key, value := range output {
			content.WriteString(key)
			content.WriteString("=")
			content.WriteString(value)
			content.WriteString("\n")
		}
		fmt.Print(content.String())
		_, err = f.WriteString(content.String())
		if err != nil {
			return err
		}
	} else {
		for key, value := range output {
			fmt.Println("::set-output name=" + key + "::" + value)
		}
	}
	return nil
}
