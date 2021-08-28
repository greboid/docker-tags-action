package main

import (
	"fmt"
	"strings"
)

func getOutput(gitRepo, inputRepo, gitRef, gitSHA, inputRegistries, separator string, fullName string) string {
	imageName := getImageName(gitRepo, inputRepo)
	registries := parseRegistriesInput(inputRegistries)
	version := refToVersion(gitRef, gitSHA)
	versions := refToVersions(gitRef)
	tags := getTags(imageName, registries, versions, getFullName(fullName))
	separator = getSeparator(separator)
	return fmt.Sprintf("::set-output name=tags::%s\n::set-output name=version::%s", strings.Join(tags, separator), version)
}
