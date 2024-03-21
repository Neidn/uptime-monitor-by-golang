package workflows

import (
	"bufio"
	"context"
	"errors"
	"github.com/Neidn/uptime-monitor-by-golang/cmd/monitor/lib"
	"github.com/Neidn/uptime-monitor-by-golang/config"
	"log"
	"os"
	"strings"
)

var ctx = context.Background()

func GenerateSummary(
	githubClient *lib.GithubClient,
	owner string,
	repo string,
	check bool,
	uptimeConfig config.UptimeConfig,
) (err error) {
	if !check {
		return errors.New("health check failed")
	}

	// Create Dir HistoryYamlDir
	if _, err := os.Stat(config.HistoryYamlDir); os.IsNotExist(err) {
		_ = os.Mkdir(config.HistoryYamlDir, 0755)
	}

	// Read Readme.md And convert bytes to string
	readmeFile, err := os.Open(config.ReadmeFile)
	if err != nil {
		log.Println("Error reading Readme.md", err)
		return
	}
	defer readmeFile.Close()

	scanner := bufio.NewScanner(readmeFile)

	// Get startText before status pages
	// if there is uptimeConfig.SummaryStartHtmlComment, use it
	// else use default
	var startText string
	var endText string

	cutPoints := map[string]string{
		"start": config.DefaultStartStatusPageText,
		"end":   config.DefaultEndStatusPageText,
	}

	if uptimeConfig.SummaryStartHtmlComment != "" {
		cutPoints["start"] = uptimeConfig.SummaryStartHtmlComment
	}

	if uptimeConfig.SummaryEndHtmlComment != "" {
		cutPoints["end"] = uptimeConfig.SummaryEndHtmlComment
	}

	cutFlag := map[string]bool{
		"start": true,
		"end":   false,
	}

	for scanner.Scan() {
		line := scanner.Text()
		if cutFlag["start"] {
			if strings.Contains(line, cutPoints["start"]) {
				cutFlag["start"] = false
				continue
			}
			startText += line + "\n"
		} else {
			if strings.Contains(line, cutPoints["end"]) {
				cutFlag["end"] = true
				continue
			}
			if cutFlag["end"] {
				endText += line + "\n"
			}
		}
	}

	for _, site := range uptimeConfig.Sites {
		slugName := lib.GetSlug(site)
		log.Println("Checking : ", slugName)

		issues, err := githubClient.GetAllIssuesForSite(owner, repo, slugName)
		if err != nil {
			log.Println("Error getting all issues for site", err)
			continue
		}

		GetUptimePercentForSite(slugName, issues)
	}

	return
}
