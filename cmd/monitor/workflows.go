package main

import (
	"context"
	"github.com/google/go-github/v59/github"
	"log"
)

var release string

func GetUptimeMonitorVersion(token string) (string, error) {
	if release != "" {
		return release, nil
	}

	client := GithubClient(token)

	opts := &github.ListOptions{
		PerPage: 1,
		Page:    1,
	}

	repositoryReleases, _, err := client.Repositories.ListReleases(context.Background(), "Neidn", "uptime-monitor-by-golang", opts)

	if err != nil {
		log.Println(err)
		return "", err
	}

	return repositoryReleases[0].GetTagName(), nil
}
