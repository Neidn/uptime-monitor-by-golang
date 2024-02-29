package workflows

import (
	"context"
	"github.com/Neidn/uptime-monitor-by-golang/cmd/monitor/lib"
	"github.com/Neidn/uptime-monitor-by-golang/config"
	"github.com/google/go-github/v59/github"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"time"
)

type SiteHistory struct {
	Url          string `yaml:"url"`
	Status       string `yaml:"status"`
	Code         int    `yaml:"code"`
	ResponseTime int    `yaml:"responseTime"`
	LastUpdated  string `yaml:"lastUpdated"`
	StartTime    string `yaml:"startTime"`
	Generator    string `yaml:"generator"`
}

type PerformanceTestResult struct {
	Result       PerformanceTestCode
	ResponseTime int
	Status       string // "up" or "down" or "degraded"
}

type PerformanceTestCode struct {
	HttpCode int
}

const (
	StatusUp         = "up"
	StatusDown       = "down"
	StatusDegraded   = "degraded"
	SiteCheckTcpPing = "tcp-ping"
	SiteCheckWS      = "ws"
)

func Update(shouldCommit bool) {
	log.Println("Update workflow started")

	check := lib.HealthCheck()
	if !check {
		log.Println("Health check failed")
		return
	}

	log.Printf("Health check passed")

	owner, repo := config.GetOwnerRepo()
	log.Printf("Owner: %s, Repo: %s", owner, repo)

	var defaultConfig config.UptimeConfig
	defaultConfig.GetConfig()

	client := lib.GithubClient()
	if client == nil {
		log.Println("Error getting client")
		return
	}

	opt := &github.IssueListByRepoOptions{
		State:     "all",
		Sort:      "created",
		Direction: "desc",
		Labels: []string{
			"maintenance",
		},
	}

	events, resp, err := client.Issues.ListByRepo(
		context.Background(),
		owner,
		repo,
		opt,
	)

	if err != nil {
		log.Println("Error getting issues", err)
		return
	}

	if resp.StatusCode != 200 {
		log.Println("Error getting issues", resp.Status)
		return
	}

	log.Printf("Found ongoing maintenance events: %d", len(events))

	for _, site := range defaultConfig.Sites {
		log.Printf("Checking : %s", site.Name)
		var testResult PerformanceTestResult

		// Delay for custom time
		if defaultConfig.Delay > 0 {
			log.Printf("Delaying for %d seconds", defaultConfig.Delay)
			lib.Delay(defaultConfig.Delay)
		}

		slugName := lib.GetSlug(site)
		currentStatus := "unknown"
		startTime := time.Now()

		siteHistory, err := ReadSiteHistory(slugName)
		if err != nil {
			log.Println("Error reading site history", err)
			siteHistory = &SiteHistory{}
		}

		// Get the status of the site
		if siteHistory.Status != "" {
			currentStatus = siteHistory.Status
		}

		if siteHistory.StartTime != "" {
			startTime, err = time.Parse(time.RFC3339, siteHistory.StartTime)
			if err != nil {
				log.Println("Error parsing start time", err)
				startTime = time.Now()
			}
		}

		log.Println(slugName, "Current status: ", currentStatus, "Start time: ", startTime)

		testResult, err = ServerCheck(site)
		if err != nil {
			log.Println("Error checking site", err)
			continue
		}

		// if the status is not up, check again
		if testResult.Status != StatusUp {
			log.Println("Status is not up, checking again")
			testResult, err = ServerCheck(site)
			if err != nil {
				log.Println("Error checking site", err)
				continue
			}

			// if the status is still not up, check again
			if testResult.Status != StatusUp {
				log.Println("Status is still not up, checking again")
				testResult, err = ServerCheck(site)
				if err != nil {
					log.Println("Error checking site", err)
					continue
				}
			}
		}

		if shouldCommit {

		}

		break
	}
}

func ServerCheck(site config.Site) (PerformanceTestResult, error) {
	switch site.Check {
	case SiteCheckTcpPing:
		tcpPingCheckResult, err := TcpPingCheck(site)
		if err != nil {
			log.Println("Error checking site using tcp ping", err)
			return PerformanceTestResult{}, err
		}
		return tcpPingCheckResult, nil

	case SiteCheckWS:
		wsCheckResult, err := WsCheck(site)
		if err != nil {
			log.Println("Error checking site using websocket", err)
			return PerformanceTestResult{}, err
		}

		return wsCheckResult, nil

	default:
		httpResult, err := HttpCheck(site)
		if err != nil {
			log.Println("Error checking site using http", err)
			return PerformanceTestResult{}, err
		}
		return httpResult, nil
	}
}

func ReadSiteHistory(slugName string) (*SiteHistory, error) {
	siteHistoryFile, err := os.ReadFile(filepath.Join(config.HistoryYamlDir, slugName+".yml"))
	if err != nil {
		return nil, err
	}

	siteHistory := SiteHistory{}
	err = yaml.Unmarshal(siteHistoryFile, &siteHistory)
	if err != nil {
		return nil, err
	}

	return &siteHistory, nil
}
