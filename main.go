package main

import (
	"fmt"
	"github.com/blang/semver/v4"
	"os"
	"strings"
)

func main() {
	fmt.Printf("Getting environmental variables\n")
	imageName := getImageName(os.Getenv("GITHUB_REPOSITORY"), os.Getenv("INPUT_REPOSITORY"))
	ref := os.Getenv("GITHUB_REF")
	fmt.Printf("Parsing registries\n")
	registries := parseRegistriesInput(os.Getenv("INPUT_REGISTRIES"))
	fmt.Printf("Getting versions\n")
	versions := refToVersions(ref)
	fmt.Printf("Getting tags\n")
	tags := getTags(imageName, registries, versions)
	fmt.Printf("::set-output name=tags::%s", strings.Join(tags, ","))
}

func getTags(imageName string, registries []string, versions []string) (tags []string){
	for _, registry := range registries {
		for _, version := range versions {
			tags = append(tags, fmt.Sprintf("%s/%s:%s", registry, imageName, version))
		}
	}
	return
}

func refToVersions(ref string) (versions []string) {
	fmt.Printf("Version: %s\n", ref)
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

func getImageName(gitRepo string, inputRepo string) string{
	if strings.TrimSpace(inputRepo) != "" {
		return inputRepo
	}
	return gitRepo
}