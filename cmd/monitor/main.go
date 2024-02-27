package main

import (
	"github.com/Neidn/uptime-monitor-by-golang/config"
	"log"
	"sync"
)

func main() {
	var wait sync.WaitGroup
	wait.Add(1)

	go run(&wait)

	wait.Wait()
}

func run(wait *sync.WaitGroup) {
	defer wait.Done()

	_TokenKey := []string{"GITHUB_TOKEN", "token", "GH_PAT"}
	var token string

	for _, key := range _TokenKey {
		token = config.GetSecret(key)
		if token != "" {
			log.Printf("Token found: %s, %s", key, token)
			break
		}
	}

	// GetUptimeMonitorVersion() is a function from workflows.go
	version, err := GetUptimeMonitorVersion(token)
	if err != nil {
		log.Println("Error getting version", err)
		return
	}

	log.Printf(`
🔼 Uptime Monitor @%s
GitHub-powered uptime monitor and status page by Neidn.

* Source: https://github.com/Neidn/uptime`, version)
}
