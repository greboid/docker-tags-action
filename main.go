package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
)

const defaultSeparator = ","

func main() {
	fmt.Printf(
		getOutput(
			os.Getenv("GITHUB_REPOSITORY"),
			os.Getenv("INPUT_REPOSITORY"),
			os.Getenv("GITHUB_REF"),
			os.Getenv("GITHUB_SHA"),
			os.Getenv("INPUT_REGISTRIES"),
			os.Getenv("INPUT_SEPARATOR"),
			os.Getenv("INPUT_FULLNAME"),
		))
}

func getSeparator(input string) string {
	if input == "" {
		return defaultSeparator
	}
	return input
}

func getFullName(input string) bool {
	return input == "" || input == "true" || input == "1"
}

func getOutput(gitRepo, inputRepo, gitRef, gitSHA, inputRegistries, separator string, fullName string) string {
	imageName := getImageName(gitRepo, inputRepo)
	registries := parseRegistriesInput(inputRegistries)
	version := refToVersion(gitRef, gitSHA)
	versions := refToVersions(gitRef)
	tags := getTags(imageName, registries, versions, getFullName(fullName))
	separator = getSeparator(separator)
	return fmt.Sprintf("::set-output name=tags::%s\n::set-output name=version::%s", strings.Join(tags, separator), version)
}

func getTags(imageName string, registries []string, versions []string, fullname bool) (tags []string) {
	for _, registry := range registries {
		for _, version := range versions {
			if fullname {
				tags = append(tags, fmt.Sprintf("%s/%s:%s", registry, imageName, version))
			} else {
				tags = append(tags, fmt.Sprintf("%s", version))
			}
		}
	}
	return
}

func refToVersion(ref string, sha string) string {
	if ref == "refs/heads/master" || ref == "refs/heads/main" {
		return sha
	}
	ref = strings.TrimPrefix(ref, "refs/tags/")
	ref = strings.TrimPrefix(ref, "v")
	version, err := semver.New(ref)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return "unknown"
	}
	return fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patch)
}

func refToVersions(ref string) (versions []string) {
	if ref == "refs/heads/master" || ref == "refs/heads/main" {
		versions = append(versions, "latest")
		return
	}
	ref = strings.TrimPrefix(ref, "refs/tags/")
	ref = strings.TrimPrefix(ref, "v")
	version, err := semver.New(ref)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	versions = append(versions, fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patch))
	versions = append(versions, fmt.Sprintf("%d.%d", version.Major, version.Minor))
	versions = append(versions, fmt.Sprintf("%d", version.Major))
	return
}

func parseRegistriesInput(input string) []string {
	input = strings.TrimSpace(input)
	if input == "" {
		return []string{"docker.io"}
	}
	if !strings.Contains(input, ",") {
		return []string{input}
	}
	var output []string
	inputSplit := strings.Split(input, ",")
	for _, reg := range inputSplit {
		reg = strings.TrimSpace(reg)
		if reg != "" {
			output = append(output, reg)
		}
	}
	return output
}

func getImageName(gitRepo string, inputRepo string) string {
	if strings.TrimSpace(inputRepo) != "" {
		return inputRepo
	}
	return gitRepo
}
