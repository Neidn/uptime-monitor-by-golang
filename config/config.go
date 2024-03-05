package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Site struct {
	Check                              string   `yaml:"check"`
	Method                             string   `yaml:"method"`
	Name                               string   `yaml:"name"`
	Url                                string   `yaml:"url"`
	Port                               int      `yaml:"port"`
	ExpectedStatusCode                 int      `yaml:"expectedStatusCode"`
	Assignees                          []string `yaml:"assignees"`
	Headers                            []string `yaml:"headers"`
	Tags                               []string `yaml:"tags"`
	Slug                               string   `yaml:"slug"`
	Body                               string   `yaml:"body"`
	Icon                               string   `yaml:"icon"`
	MaxResponseTime                    int      `yaml:"maxResponseTime"`
	MaxRedirects                       int      `yaml:"maxRedirects"`
	Verbose                            bool     `yaml:"verbose"`
	DangerousInsecure                  bool     `yaml:"dangerousInsecure"`
	DangerousDisableVerifyPeer         bool     `yaml:"dangerousDisableVerifyPeer"`
	DangerousDisableVerifyHost         bool     `yaml:"dangerousDisableVerifyHost"`
	DangerousBodyDown                  string   `yaml:"dangerousBodyDown"`
	DangerousBodyDownIfTextMissing     string   `yaml:"dangerousBodyDownIfTextMissing"`
	DangerousBodyDegraded              string   `yaml:"dangerousBodyDegraded"`
	DangerousBodyDegradedIfTextMissing string   `yaml:"dangerousBodyDegradedIfTextMissing"`
}

type navbar struct {
	Title string `yaml:"title"`
	Url   string `yaml:"url"`
}

type statusWebsite struct {
	Cname        string   `yaml:"cname"`
	LogoUrl      string   `yaml:"logoUrl"`
	Name         string   `yaml:"name"`
	IntroTitle   string   `yaml:"introTitle"`
	IntroMessage string   `yaml:"introMessage"`
	Navbar       []navbar `yaml:"navbar"`
	Publish      bool     `yaml:"publish"`
	SingleCommit bool     `yaml:"singleCommit"`
}

type CommitMessage struct {
	Content      string `yaml:"content"`
	Summary      string `yaml:"summary"`
	StatusChange string `yaml:"statusChange"`
	GraphUpdate  string `yaml:"graphUpdate"`
	AuthorName   string `yaml:"authorName"`
	AuthorEmail  string `yaml:"authorEmail"`
}

type i18n struct {
	PrefixStatus          PrefixStatus `yaml:"prefixStatus"`
	Url                   string       `yaml:"url"`
	Status                string       `yaml:"status"`
	History               string       `yaml:"history"`
	Ms                    string       `yaml:"ms"`
	ResponseTime          timeStruct   `yaml:"responseTime"`
	Uptime                timeStruct   `yaml:"uptime"`
	ResponseTimeGraphAlt  string       `yaml:"responseTimeGraphAlt"`
	LiveStatus            string       `yaml:"liveStatus"`
	AllSystemsOperational string       `yaml:"allSystemsOperational"`
	DegradedPerformance   string       `yaml:"degradedPerformance"`
	CompleteOutage        string       `yaml:"completeOutage"`
	PartialOutage         string       `yaml:"partialOutage"`
}

type timeStruct struct {
	Time  string `yaml:"time"`
	Day   string `yaml:"day"`
	Week  string `yaml:"week"`
	Month string `yaml:"month"`
	Year  string `yaml:"year"`
}

type PrefixStatus struct {
	Up       string `yaml:"up"`
	Down     string `yaml:"down"`
	Degraded string `yaml:"degraded"`
}

type UptimeConfig struct {
	Owner                   string        `yaml:"owner"`
	Repo                    string        `yaml:"repo"`
	UserAgent               string        `yaml:"user-agent"`
	Sites                   []Site        `yaml:"sites"`
	Assignees               []string      `yaml:"assignees"`
	Delay                   int           `yaml:"delay"`
	PAT                     string        `yaml:"PAT"`
	StatusWebSite           statusWebsite `yaml:"status-website"`
	SkipDescriptionUpdate   bool          `yaml:"skipDescriptionUpdate"`
	SkipTopicUpdate         bool          `yaml:"skipTopicUpdate"`
	SkipHomepageUpdate      bool          `yaml:"skipHomepageUpdate"`
	SkipDeletedIssues       bool          `yaml:"skipDeletedIssues"`
	SkipPoweredByReadme     bool          `yaml:"skipPoweredByReadme"`
	CommitMessages          CommitMessage `yaml:"commit-messages"`
	SummaryStartHtmlComment string        `yaml:"summaryStartHtmlComment"`
	SummaryEndHtmlComment   string        `yaml:"summaryEndHtmlComment"`
	LiveStatusHtmlComment   string        `yaml:"liveStatusHtmlComment"`
	CommitPrefixStatus      PrefixStatus  `yaml:"commit-prefix-status"`
	i18n                    i18n          `yaml:"i18n"`
}

func (c *UptimeConfig) GetConfig() *UptimeConfig {
	buf, err := os.ReadFile(".upptimerc.yml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(buf, &c)
	if err != nil {
		panic(err)
	}

	return c
}
