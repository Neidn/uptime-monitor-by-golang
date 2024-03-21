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

type GithubClient struct {
	Client *github.Client
}

var (
	IssueStatusOpen   = "open"
	IssueStatusClosed = "closed"
	IssueStatusAll    = "all"
)

var ctx = context.Background()

func NewGithubClient() (githubClient *GithubClient, err error) {
	token := config.GetToken()
	if token == "" {
		return nil, errors.New("token not found")
	}
	log.Println("token", token)

	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tokenClient := oauth2.NewClient(context.Background(), tokenService)

	githubClient = &GithubClient{
		Client: github.NewClient(tokenClient),
	}

	return githubClient, nil
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

func SendPush() {
	_ = exec.Command("git", "push").Run()
}

func (g *GithubClient) GetRepoReleases(
	owner string,
	repo string,
) (releases []*github.RepositoryRelease, err error) {
	opts := &github.ListOptions{
		PerPage: 1,
		Page:    1,
	}

	releases, _, err = g.Client.Repositories.ListReleases(
		ctx,
		owner,
		repo,
		opts,
	)
	return
}

func (g *GithubClient) GetIssues(
	owner string,
	repo string,
	slugName string,
) ([]*github.Issue, error) {
	issues, _, err := g.Client.Issues.ListByRepo(ctx, owner, repo, &github.IssueListByRepoOptions{
		State:     IssueStatusOpen,
		Sort:      "created",
		Direction: "desc",
		Labels: []string{
			slugName,
		},
	})
	return issues, err
}

func (g *GithubClient) CheckAndCloseMaintenanceEvents(
	owner string,
	repo string,
) (ongoingEvents []OnGoingMaintenanceEvent, err error) {
	EventOpt := &github.IssueListByRepoOptions{
		State:     IssueStatusAll,
		Sort:      "created",
		Direction: "desc",
		Labels: []string{
			"maintenance",
		},
	}

	events, resp, err := g.Client.Issues.ListByRepo(
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

func (g *GithubClient) CreateNewIssue(
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

	newIssue, _, err = g.Client.Issues.Create(ctx, owner, repo, issue)
	if err != nil {
		return nil, err
	}

	return newIssue, nil
}

func (g *GithubClient) GetAllIssuesForSite(
	owner string,
	repo string,
	slugName string,
) (issues []*github.Issue, err error) {
	issues, _, err = g.Client.Issues.ListByRepo(ctx, owner, repo, &github.IssueListByRepoOptions{
		State:     IssueStatusAll,
		Sort:      "created",
		Direction: "desc",
		Labels: []string{
			"status",
			slugName,
		},
		// Filter All
		ListOptions: github.ListOptions{PerPage: 100},
	})
	return
}

func (g *GithubClient) AddAssignees(
	owner string,
	repo string,
	issueNumber int,
	assignees []string,
) (err error) {
	_, _, err = g.Client.Issues.AddAssignees(ctx, owner, repo, issueNumber, assignees)
	return
}

func (g *GithubClient) LockIssue(
	owner string,
	repo string,
	issueNumber int,
) (err error) {
	_, err = g.Client.Issues.Lock(ctx, owner, repo, issueNumber, &github.LockIssueOptions{})
	return
}

func (g *GithubClient) UnlockIssue(
	owner string,
	repo string,
	issueNumber int,
) (err error) {
	_, err = g.Client.Issues.Unlock(ctx, owner, repo, issueNumber)
	return
}

func (g *GithubClient) CreateComment(
	owner string,
	repo string,
	issueNumber int,
	body string,
) (err error) {
	_, _, err = g.Client.Issues.CreateComment(ctx, owner, repo, issueNumber, &github.IssueComment{
		Body: &body,
	})
	return
}

func (g *GithubClient) CloseIssue(
	owner string,
	repo string,
	issueNumber int,
) (err error) {
	_, _, err = g.Client.Issues.Edit(ctx, owner, repo, issueNumber, &github.IssueRequest{
		State: &IssueStatusClosed,
	})
	return
}
