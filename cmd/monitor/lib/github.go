package lib

import (
	"context"
	"github.com/Neidn/uptime-monitor-by-golang/config"
	"github.com/google/go-github/v59/github"
	"golang.org/x/oauth2"
	"os/exec"
	"strings"
)

var ctx = context.Background()

func GithubClient() *github.Client {
	token := config.GetToken()
	if token == "" {
		return nil
	}

	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tokenClient := oauth2.NewClient(context.Background(), tokenService)

	client := github.NewClient(tokenClient)

	return client
}

func SendCommit(message string, name string, email string) {
	//_ = exec.Command("git", "config", "user.name", name).Run()
	//_ = exec.Command("git", "config", "user.email", email).Run()
	_ = exec.Command("git", "add", ".").Run()
	_ = exec.Command("git", "commit", "-m", message).Run()
}

func LastCommit() string {
	out, _ := exec.Command("git", "log", "-1", "--pretty=%B").Output()
	return string(out)
}

func GetIssues(client *github.Client, owner string, repo string) ([]*github.Issue, error) {
	issues, _, err := client.Issues.ListByRepo(ctx, owner, repo, nil)
	return issues, err
}

func UpdateMaintenanceEvents(events *[]github.Issue) {
	metadata := map[string]string{}

	for _, event := range *events {
		if event.Body != nil && strings.Contains(*event.Body, "<!--") {
			summary := strings.Split(*event.Body, "<!--")[1]
			summary = strings.Split(summary, "-->")[0]

			lines := strings.Split(summary, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.Contains(line, ":") {

				}
			}

		}
	}
}
