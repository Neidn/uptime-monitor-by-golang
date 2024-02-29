package main

import (
	"github.com/Neidn/uptime-monitor-by-golang/cmd/monitor/workflows"
	"github.com/Neidn/uptime-monitor-by-golang/config"
	"github.com/actions-go/toolkit/core"
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

	// GetUptimeMonitorVersion() is a function from lib/version.go
	version, err := GetUptimeMonitorVersion()
	if err != nil {
		log.Println("Error getting version", err)
		return
	}

	log.Printf(`
🔼 Uptime Monitor @%s
GitHub-powered uptime monitor and status page by Neidn.

* Source: %s`,
		version,
		config.Generator,
	)

	command, _ := core.GetInput(config.GithubActionsInputKey)

	switch command {
	case config.CommandSummary:
		core.Debug("Starting summary")
		return

	case config.CommandReadme:
		core.Debug("Starting readme")
		return

	case config.CommandSite:
		core.Debug("Starting site")
		return

	case config.CommandGraph:
		core.Debug("Starting graph")
		return

	case config.CommandResponseTime:
		core.Debug("Starting response time")
		return

	case config.CommandUpdateDependencies:
		core.Debug("Starting update dependencies")
		return

	case config.CommandUpdateTemplate:
		core.Debug("Starting update template")
		return

	default:
		core.Debug("Starting update template")
		workflows.Update(false)
		return
	}
}
