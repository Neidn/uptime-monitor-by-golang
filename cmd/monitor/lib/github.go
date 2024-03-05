package lib

import (
	"context"
	"errors"
	"github.com/Neidn/uptime-monitor-by-golang/config"
	"github.com/google/go-github/v59/github"
	"golang.org/x/oauth2"
	"log"
	"os/exec"
	"strings"
)

type OnGoingMaintenanceEvent struct {
	IssueNumber int
	Metadata    EventMetadata
}

type EventMetadata struct {
	Start            string
	End              string
	ExpectedDown     []string
	ExpectedDegraded []string
}

var ctx = context.Background()

var EventOpt = &github.IssueListByRepoOptions{
	State:     "all",
	Sort:      "created",
	Direction: "desc",
	Labels: []string{
		"maintenance",
	},
}

func GithubClient() (*github.Client, error) {
	token := config.GetToken()
	if token == "" {
		return nil, errors.New("token not found")
	}
	log.Println("token", token)

	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tokenClient := oauth2.NewClient(context.Background(), tokenService)

	client := github.NewClient(tokenClient)

	return client, nil
}

func SendCommit(message string, name string, email string) {
	//_ = exec.Command("git", "config", "user.name", name).Run()
	//_ = exec.Command("git", "config", "user.email", email).Run()
	_ = exec.Command("git", "add", ".").Run()
	_ = exec.Command("git", "commit", "-m", message).Run()
}

func LastCommit() string {
	out, _ := exec.Command("git", "log", "-1", "--format=%H").Output()
	return string(out)
}

func GetIssues(
	client *github.Client,
	owner string,
	repo string,
	slugName string,
) ([]*github.Issue, error) {
	issues, _, err := client.Issues.ListByRepo(ctx, owner, repo, &github.IssueListByRepoOptions{
		State:     "open",
		Sort:      "created",
		Direction: "desc",
		Labels: []string{
			slugName,
		},
	})
	return issues, err
}

func CheckAndCloseMaintenanceEvents(client *github.Client, owner string, repo string) (ongoingEvents []OnGoingMaintenanceEvent, err error) {
	events, resp, err := client.Issues.ListByRepo(
		context.Background(),
		owner,
		repo,
		EventOpt,
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

	for _, event := range events {
		metadata := map[string]string{}

		if event.Body != nil && strings.Contains(*event.Body, "<!--") {
			summary := strings.Split(*event.Body, "<!--")[1]
			summary = strings.Split(summary, "-->")[0]

			lines := strings.Split(summary, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.Contains(line, ":") {
					parts := strings.Split(line, ":")
					if len(parts) == 2 {
						metadata[parts[0]] = parts[1]
					}
				}
			}
		}

		log.Println("Metadata: ", metadata)
	}

	return
}

func CreateNewIssue(
	client *github.Client,
	owner string,
	repo string,
	title string,
	body string,
	labels []string,
) (newIssue *github.Issue, err error) {
	issue := &github.IssueRequest{
		Title:  &title,
		Body:   &body,
		Labels: &labels,
	}

	newIssue, _, err = client.Issues.Create(ctx, owner, repo, issue)
	if err != nil {
		return nil, err
	}

	return newIssue, nil
}

func AddAssignees(
	client *github.Client,
	owner string,
	repo string,
	issueNumber int,
	assignees []string,
) (err error) {
	_, _, err = client.Issues.AddAssignees(ctx, owner, repo, issueNumber, assignees)
	return
}

func LockIssue(
	client *github.Client,
	owner string,
	repo string,
	issueNumber int,
) (err error) {
	_, err = client.Issues.Lock(ctx, owner, repo, issueNumber, &github.LockIssueOptions{})
	return
}
