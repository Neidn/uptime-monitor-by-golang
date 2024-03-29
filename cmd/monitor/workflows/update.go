package workflows

import (
	"fmt"
	"github.com/Neidn/uptime-monitor-by-golang/cmd/monitor/lib"
	"github.com/Neidn/uptime-monitor-by-golang/config"
	"log"
	"slices"
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
	hasDelta := false
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

	//client, err := githubClient.GetGithubClient()
	githubClient, err := lib.NewGithubClient()
	if err != nil {
		log.Println("Error getting github client", err)
		return
	}

	ongoingEvents, err := githubClient.CheckAndCloseMaintenanceEvents(owner, repo)

	if err != nil {
		log.Println("Error checking and closing maintenance events", err)
		return
	}

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
			startTime, err = time.Parse(config.TimeFormat, siteHistory.StartTime)
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

		if !shouldCommit && currentStatus == testResult.Status {
			log.Println("Skipping commit, status is ", testResult.Status)
		} else {
			siteHistory.Status = testResult.Status
			siteHistory.Code = testResult.Result.HttpCode
			siteHistory.ResponseTime = testResult.ResponseTime
			siteHistory.LastUpdated = time.Now().Format(config.TimeFormat)
			siteHistory.StartTime = startTime.Format(config.TimeFormat)
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

			if currentStatus == testResult.Status {
				log.Println("Status is the same", currentStatus, testResult.Status)
			} else {
				log.Println("Status changed from", currentStatus, "to", testResult.Status)
				hasDelta = true
				issues, err := githubClient.GetIssues(owner, repo, slugName)
				if err != nil {
					log.Println("Error getting issues", err)
					continue
				}
				log.Println("Found ", len(issues), " issues")

				expected := false

				switch testResult.Status {
				case StatusDown:
					// Check if there is match an ongoing maintenance event's metadata expectedDown with slugName
					for _, event := range ongoingEvents {
						if slices.Contains(event.Metadata.ExpectedDown, slugName) {
							expected = true
							break
						}
					}

				case StatusDegraded:
					// Check if there is match an ongoing maintenance event's metadata expectedDegraded with slugName
					for _, event := range ongoingEvents {
						if slices.Contains(event.Metadata.ExpectedDegraded, slugName) {
							expected = true
							break
						}
					}
				}

				if testResult.Status != StatusUp && !expected {
					//if !expected {
					if len(issues) > 0 {
						log.Println("Issue already exists")
					} else {
						log.Println("Creating issue")
						title, body, labels := CreateIssueMessage(owner, repo, slugName, lastCommit, testResult, site)
						newIssue, err := githubClient.CreateNewIssue(owner, repo, title, body, labels)
						if err != nil {
							log.Println("Error creating issue", err)
							continue
						}

						log.Println("Issue created", newIssue)

						// Add assignees
						assignees := append(site.Assignees, defaultConfig.Assignees...)
						if len(assignees) > 0 {
							err = githubClient.AddAssignees(owner, repo, *newIssue.Number, assignees)
							if err != nil {
								log.Println("Error adding assignees", err)
							}
						}

						// Lock the issue
						err = githubClient.LockIssue(owner, repo, *newIssue.Number)
						if err != nil {
							log.Println("Error locking issue", err)
						}

						log.Println("Opened and locked issue")
					}
				} else if len(issues) > 0 {
					log.Println("UnLocking issue")
					err = githubClient.UnlockIssue(owner, repo, issues[0].GetNumber())
					if err != nil {
						log.Println("Error unlocking issue", err)
					}

					commentMsg := CreateCommentMessage(owner, repo, lastCommit, issues[0], site)

					err = githubClient.CreateComment(owner, repo, issues[0].GetNumber(), commentMsg)
					if err != nil {
						log.Println("Error creating comment", err)
					}
					log.Println("Comment created")

					// Close the issue
					_ = githubClient.CloseIssue(owner, repo, issues[0].GetNumber())

					// Lock the issue
					_ = githubClient.LockIssue(owner, repo, issues[0].GetNumber())
				} else {
					log.Println("No Relevant issue found")
				}
			}

		}
	}

	// Git Push
	if false {
		//if len(defaultConfig.Sites) > 0 {
		lib.SendPush()
	}

	log.Println("Has delta: ", hasDelta)

	if hasDelta {
		GenerateSummary(
			githubClient,
			owner,
			repo,
			check,
			defaultConfig,
		)
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
