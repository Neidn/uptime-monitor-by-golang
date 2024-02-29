package lib

import (
	"context"
	"github.com/Neidn/uptime-monitor-by-golang/config"
	"github.com/google/go-github/v59/github"
	"golang.org/x/oauth2"
	"os/exec"
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
