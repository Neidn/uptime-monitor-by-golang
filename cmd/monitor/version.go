package main

import (
	"context"
	"github.com/Neidn/uptime-monitor-by-golang/cmd/monitor/lib"
	"github.com/Neidn/uptime-monitor-by-golang/config"
	"github.com/google/go-github/v59/github"
	"log"
)

var release string

func GetUptimeMonitorVersion() (string, error) {
	if release != "" {
		return release, nil
	}

	client, err := lib.GithubClient()
	if err != nil {
		log.Println(err)
		return "", err
	}

	opts := &github.ListOptions{
		PerPage: 1,
		Page:    1,
	}

	repositoryReleases, _, err := client.Repositories.ListReleases(
		context.Background(),
		config.OwnerName,
		config.MonitorRepositoryName,
		opts,
	)

	if err != nil {
		log.Println(err)
		return "", err
	}

	return repositoryReleases[0].GetTagName(), nil
}
