package main

import (
	"github.com/Neidn/uptime-monitor-by-golang/cmd/monitor/lib"
	"github.com/Neidn/uptime-monitor-by-golang/config"
	"log"
)

var release string

func GetUptimeMonitorVersion() (string, error) {
	if release != "" {
		return release, nil
	}

	githubClient, err := lib.NewGithubClient()
	if err != nil {
		log.Println(err)
		return "", err
	}

	repositoryReleases, err := githubClient.GetRepoReleases(
		config.OwnerName,
		config.RepositoryName,
	)

	if err != nil {
		log.Println(err)
		return "", err
	}

	if len(repositoryReleases) == 0 {
		return "", nil
	}

	return repositoryReleases[0].GetTagName(), nil
}
