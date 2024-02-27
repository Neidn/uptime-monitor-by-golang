package main

import (
	"context"
	"github.com/google/go-github/v59/github"
	"golang.org/x/oauth2"
)

func GithubClient(token string) *github.Client {
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tokenClient := oauth2.NewClient(context.Background(), tokenService)

	client := github.NewClient(tokenClient)

	return client
}
