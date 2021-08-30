package main

import (
	"errors"
	"strings"
)

func getSeparator(input string) string {
	if input == "" {
		return defaultSeparator
	}
	return input
}

func getFullName(input string) bool {
	return input == "" || input == "true" || input == "1"
}

func getImageName(gitRepo string, inputRepo string) string {
	if strings.TrimSpace(inputRepo) != "" {
		return inputRepo
	}
	return gitRepo
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

func splitRepo(repo string) (string, string, error) {
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return "", "", errors.New("invalid repo format")
	}
	return parts[0], parts[1], nil
}
