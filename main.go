package main

import (
	"fmt"
	"os"
)

const defaultSeparator = ","

func main() {
	gitTags, err := getGitTags(os.Getenv("GITHUB_WORKSPACE"))
	if err != nil {
		fmt.Printf("::error ::Unable to pull tags")
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
