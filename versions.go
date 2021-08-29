package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/blang/semver/v4"
)

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

func refToVersions(ref string, latestVersion *semver.Version) (versions []string) {
	if ref == "refs/heads/master" || ref == "refs/heads/main" {
		versions = append(versions, "dev")
		return
	}
	ref = strings.TrimPrefix(ref, "refs/tags/")
	ref = strings.TrimPrefix(ref, "v")
	version, err := semver.New(ref)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	if latestVersion.Major == version.Major {
		versions = append(versions, "latest")
	}
	versions = append(versions, fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patch))
	versions = append(versions, fmt.Sprintf("%d.%d", version.Major, version.Minor))
	versions = append(versions, fmt.Sprintf("%d", version.Major))
	return
}

func getLatestVersion(versions []string) (*semver.Version, error) {
	hasVersion := false
	latestVersion := semver.MustParse("0.0.0")
	for index := range versions {
		semVersion, err := semver.New(versions[index])
		if err != nil {
			log.Printf("Error parsing semver: %s", err)
			continue
		}
		if len(semVersion.Pre) > 0 {
			log.Printf("Ignoring pre-release")
			continue
		}
		if semVersion.GT(latestVersion) {
			hasVersion = true
			latestVersion = *semVersion
		}
	}
	if hasVersion {
		return &latestVersion, nil
	}
	return &latestVersion, errors.New("no versions")
}
