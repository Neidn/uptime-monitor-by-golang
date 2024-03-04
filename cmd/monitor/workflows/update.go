package workflows

import (
	"context"
	"fmt"
	"github.com/Neidn/uptime-monitor-by-golang/cmd/monitor/lib"
	"github.com/Neidn/uptime-monitor-by-golang/config"
	"github.com/google/go-github/v59/github"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"strings"
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

		//if shouldCommit || currentStatus != testResult.Status {
		if true {
			siteHistory.Status = testResult.Status
			siteHistory.Code = testResult.Result.HttpCode
			siteHistory.ResponseTime = testResult.ResponseTime
			siteHistory.LastUpdated = time.Now().Format(time.RFC3339)
			siteHistory.StartTime = startTime.Format(time.RFC3339)
			siteHistory.Generator = fmt.Sprintf("%s %s", config.RepositoryName, config.Generator)

			// Write the history
			err = WriteSiteHistory(slugName, siteHistory)
			if err != nil {
				log.Println("Error writing site history", err)
				continue
			}

			// Commit the changes
			var message string

			if defaultConfig.CommitMessages.StatusChange != "" {
				message = defaultConfig.CommitMessages.StatusChange
			} else {
				message = config.DefaultCommitMessage
			}

			message = ReplaceCommitMessage(
				message,
				defaultConfig.CommitPrefixStatus,
				testResult,
				site,
				config.RepositoryName,
			)
			log.Println("Commit message", message)

			authorName := defaultConfig.CommitMessages.AuthorName
			if authorName == "" {
				authorName = *config.AuthorName
			}
			authorEmail := defaultConfig.CommitMessages.AuthorEmail
			if authorEmail == "" {
				authorEmail = *config.AuthorEmail
			}
			//lib.SendCommit(message, authorName, authorEmail)
			lastCommit := lib.LastCommit()
			log.Println("Last commit", lastCommit)

			if currentStatus != testResult.Status {
				log.Println("Status changed from", currentStatus, "to", testResult.Status)
				//hasDelta := false
				issues, err := lib.GetIssues(client, owner, repo)
				if err != nil {
					log.Println("Error getting issues", err)
					continue
				}
				log.Println("Found ", len(issues), " issues")

				expected := false
				if testResult.Status != StatusUp {

				}
			}

		}
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

func WriteSiteHistory(slugName string, siteHistory *SiteHistory) error {
	log.Println("siteHistory", siteHistory)

	historyBody := fmt.Sprintf(`url: %s
status: %s
code: %d
responseTime: %d
lastUpdated: %s
startTime: %s
generator: %s
`, siteHistory.Url, siteHistory.Status, siteHistory.Code, siteHistory.ResponseTime, siteHistory.LastUpdated, siteHistory.StartTime, siteHistory.Generator)

	_ = os.WriteFile(filepath.Join(config.HistoryYamlDir, slugName+".yml"), []byte(historyBody), 0644)

	return nil
}

func ReplaceCommitMessage(
	message string,
	prefixStatus config.PrefixStatus,
	performanceTestResult PerformanceTestResult,
	site config.Site,
	repositoryName string,
) string {
	var prefix string
	switch performanceTestResult.Status {
	case StatusUp:
		if prefixStatus.Up != "" {
			prefix = prefixStatus.Up
		} else {
			prefix = config.DefaultUp
		}
		break
	case StatusDegraded:
		if prefixStatus.Degraded != "" {
			prefix = prefixStatus.Degraded
		} else {
			prefix = config.DefaultDegraded
		}
		break
	default:
		if prefixStatus.Down != "" {
			prefix = prefixStatus.Down
		} else {
			prefix = config.DefaultDown
		}

	}
	message = strings.ReplaceAll(message, "$PREFIX", prefix)

	message = strings.ReplaceAll(message, "$SITE_NAME", site.Name)
	message = strings.ReplaceAll(message, "$SITE_URL", site.Url)
	message = strings.ReplaceAll(message, "$STATUS", performanceTestResult.Status)
	message = strings.ReplaceAll(message, "$RESPONSE_CODE", fmt.Sprintf("%d", performanceTestResult.Result.HttpCode))
	message = strings.ReplaceAll(message, "$RESPONSE_TIME", fmt.Sprintf("%d", performanceTestResult.ResponseTime))
	message = strings.ReplaceAll(message, "$REPOSITORY_NAME", repositoryName)

	return message
}
