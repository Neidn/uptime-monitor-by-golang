package workflows

import (
	"fmt"
	"github.com/Neidn/uptime-monitor-by-golang/cmd/monitor/lib"
	"github.com/Neidn/uptime-monitor-by-golang/config"
	"github.com/google/go-github/v59/github"
	"strings"
)

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

func CreateIssueMessage(
	owner string,
	repo string,
	slugName string,
	lastCommit string,
	testResult PerformanceTestResult,
	site config.Site,
) (title string, body string, labels []string) {
	if testResult.Status == StatusDown {
		title = fmt.Sprintf(`🚨 %s is down`, site.Name)
	} else {
		title = fmt.Sprintf(`⚠️ %s has degraded performance`, site.Name)
	}

	var statusWord string
	if testResult.Status == StatusDown {
		statusWord = "was **down**"
	} else {
		statusWord = "had **degraded performance**"
	}

	body = fmt.Sprintf(`In [%s](https://github.com/%s/%s), %s (%s) %s
- HTTP code: %d
- Response time: %d ms
`, lastCommit[:7], owner, repo, site.Name, site.Url, statusWord, testResult.Result.HttpCode, testResult.ResponseTime)

	labels = []string{"status", slugName}
	labels = append(labels, site.Tags...)

	return title, body, labels
}

func CreateCommentMessage(
	owner string,
	repo string,
	lastCommit string,
	issue *github.Issue,
	site config.Site,
) (commentMsg string) {
	var issueState string
	if strings.Contains(*issue.Title, "degraded") {
		issueState = "performance has improved"
	} else {
		issueState = "Site is back up"
	}

	minutes := lib.ConvertDateToHumanReadableTimeDifference(issue.GetCreatedAt())

	commentMsg = fmt.Sprintf(
		`**Resolved:** %s %s In [%s](https://github.com/%s/%s/commit/%s) after %s`,
		site.Name, issueState, lastCommit[:7], owner, repo, lastCommit, minutes,
	)
	return
}
