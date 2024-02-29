package workflows

import (
	"context"
	"fmt"
	"github.com/Neidn/uptime-monitor-by-golang/cmd/monitor/lib"
	"github.com/Neidn/uptime-monitor-by-golang/config"
	"github.com/google/go-github/v59/github"
	"github.com/gorilla/websocket"
	probing "github.com/prometheus-community/pro-bing"
	"gopkg.in/yaml.v3"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SiteHistory struct {
	Url          string `yaml:"url"`
	Status       string `yaml:"status"`
	Code         int    `yaml:"code"`
	ResponseTime int    `yaml:"responseTime"`
	LastUpdated  string `yaml:"lastUpdated"`
	StartTime    string `yaml:"startTime"`
	Generator    string `yaml:"generator"`
}

type PerformanceTest struct {
	Result       PerformanceTestResult
	ResponseTime int
	Status       string // "up" or "down" or "degraded"
}

type PerformanceTestResult struct {
	HttpCode int
}

const (
	StatusUp         = "up"
	StatusDown       = "down"
	StatusDegraded   = "degraded"
	SiteCheckTcpPing = "tcp-ping"
	SiteCheckWS      = "ws"
)

func Update() {

	log.Println("Update workflow started")

	check := lib.HealthCheck()
	if !check {
		log.Println("Health check failed")
		return
	}

	log.Printf("Health check passed")

	owner, repo := config.GetOwnerRepo()
	log.Printf("Owner: %s, Repo: %s", owner, repo)

	var defaultConfig config.UptimeConfig
	defaultConfig.GetConfig()

	client := lib.GithubClient()
	if client == nil {
		log.Println("Error getting client")
		return
	}

	opt := &github.IssueListByRepoOptions{
		State:     "all",
		Sort:      "created",
		Direction: "desc",
		Labels: []string{
			"maintenance",
		},
	}

	events, resp, err := client.Issues.ListByRepo(
		context.Background(),
		owner,
		repo,
		opt,
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

	for _, site := range defaultConfig.Sites {
		log.Printf("Checking : %s", site.Name)

		// Delay for custom time
		if defaultConfig.Delay > 0 {
			log.Printf("Delaying for %d seconds", defaultConfig.Delay)
			lib.Delay(defaultConfig.Delay)
		}

		slugName := lib.GetSlug(site)
		currentStatus := "unknown"
		startTime := time.Now()

		siteHistoryFile, err := os.ReadFile(filepath.Join(config.HistoryYamlDir, slugName+".yml"))
		if err != nil {
			log.Println("Error reading history", err)
			continue
		}

		siteHistory := SiteHistory{}
		err = yaml.Unmarshal(siteHistoryFile, &siteHistory)
		if err != nil {
			log.Println("Error unmarshalling history", err)
			continue
		}

		// Get the status of the site
		if siteHistory.Status != "" {
			currentStatus = siteHistory.Status
		}

		if siteHistory.StartTime != "" {
			startTime, err = time.Parse(time.RFC3339, siteHistory.StartTime)
			if err != nil {
				log.Println("Error parsing start time", err)
				startTime = time.Now()
			}
		}

		log.Println(slugName, "Current status: ", currentStatus, "Start time: ", startTime)

		switch site.Check {
		case SiteCheckTcpPing:
			_, err := TcpPingCheck(site)
			if err != nil {
				log.Println("Error checking site", err)
				continue
			}
			return
		case SiteCheckWS:
			log.Println("Checking ws")
			WsCheck(site)
			return
		default:
			log.Println("Checking http")

			log.Println("Performance test")

		}

		break
	}

}

var FailedPerformanceTest = PerformanceTest{
	Status:       StatusDown,
	ResponseTime: 0,
	Result: PerformanceTestResult{
		HttpCode: 0,
	},
}

func TcpPingCheck(site config.Site) (PerformanceTest, error) {
	address := site.Url
	status := StatusUp
	var responseTime time.Duration
	ip := net.ParseIP(site.Url)

	if ip == nil {
		_url, err := url.Parse(site.Url)
		if err != nil {
			log.Println("Error parsing URL", err)
			return FailedPerformanceTest, err
		}

		_hostname := strings.TrimPrefix(_url.Hostname(), "www.")

		ipList, err := net.LookupIP(_hostname)
		if err != nil {
			log.Println("Error looking up IP", err)
			return PerformanceTest{}, err
		}

		address = ipList[0].String()
	}

	// with specific port
	if site.Port != 0 {
		address = fmt.Sprintf("%s:%d", address, site.Port)
	}

	tcpPing, err := probing.NewPinger(address)
	defer tcpPing.Stop()
	if err != nil {
		log.Println(err)
		return FailedPerformanceTest, err
	}
	tcpPing.Count = config.TcpPingDefaultCount

	tcpPing.OnFinish = func(stats *probing.Statistics) {
		responseTime = stats.AvgRtt

		if responseTime > time.Duration(config.MaxResponseTime)*time.Millisecond {
			status = StatusDegraded
		}
	}

	err = tcpPing.Run()
	if err != nil {
		return FailedPerformanceTest, err
	}

	if responseTime < 0 {
		return FailedPerformanceTest, nil
	}

	return PerformanceTest{
		Result: PerformanceTestResult{
			HttpCode: 200,
		},
		ResponseTime: int(responseTime.Milliseconds()),
		Status:       status,
	}, nil
}

func WsCheck(site config.Site) (PerformanceTest, error) {
	log.Println("Using WebSocket to check site")
	status := StatusUp

	c, _, err := websocket.DefaultDialer.Dial(site.Url, nil)
	if err != nil {
		log.Println("Error dialing", err)
		return FailedPerformanceTest, err
	}
	defer c.Close()

	_textMessage := []byte("")
	if site.Body != "" {
		_textMessage = []byte(site.Body)
	}

	err = c.WriteMessage(websocket.TextMessage, _textMessage)
	if err != nil {
		log.Println("Error writing message", err)
		return FailedPerformanceTest, err
	}

	return PerformanceTest{
		Status:       status,
		ResponseTime: 0,
	}, nil
}
