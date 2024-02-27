package main

import (
	"context"
	"github.com/google/go-github/v59/github"
	"golang.org/x/oauth2"
	"log"
)

func GithubClient(token string) (*github.Client, error) {
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tokenClient := oauth2.NewClient(context.Background(), tokenService)

	client := github.NewClient(tokenClient)

	_ = &github.ListOptions{
		PerPage: 10,
		Page:    1,
	}
	log.Printf("Token: %s", token)

	//orgs, _, err := client.Repositories.Get(context.Background(), "Neidn", "uptime-monitor-by-golang")
	orgs, _, err := client.Repositories.GetRelease(context.Background(), "Neidn", "uptime-monitor-by-golang", 1)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println("Repositories:")
	log.Println(orgs)

	return client, nil

}
