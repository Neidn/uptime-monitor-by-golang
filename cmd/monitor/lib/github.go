package lib

import (
	"context"
	"github.com/Neidn/uptime-monitor-by-golang/config"
	"github.com/google/go-github/v59/github"
	"golang.org/x/oauth2"
)

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
