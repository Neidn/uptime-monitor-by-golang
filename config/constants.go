package config

import (
	"flag"
	"time"
)

const TimeFormat = time.RFC3339

const OwnerName = "Neidn"
const MonitorRepositoryName = "uptime-monitor-by-golang"
const RepositoryName = "uptime"
const CommitBranch = "main"
const Generator = "https://github.com/Neidn/uptime"

const GithubActionsInputKey = "command"

const SecretContextKey = "SECRETS_CONTEXT"
const GithubRepoKey = "GITHUB_REPOSITORY"
const NotificationDownMessageKey = "NOTIFICATION_DOWN_MESSAGE"
const DefaultCommitMessage = "$PREFIX $SITE_NAME is $STATUS ($RESPONSE_CODE in $RESPONSE_TIME ms) [skip ci] [$REPOSITORY_NAME]"

const UptimeRcYaml = ".uptimerc.yml"
const ReadmeFile = "README.md"
const HistoryYamlDir = "history"
const DefaultStartStatusPageText = "<!-- start: status pages -->"
const DefaultEndStatusPageText = "<!-- end: status pages -->"

const GraphsCiSchedule = "0 0 * * *"
const ResponseTimeCiSchedule = "0 23 * * *"
const StaticSiteCiSchedule = "0 1 * * *"
const SummaryCiSchedule = "0 0 * * *"
const UpdateTemplateCiSchedule = "0 0 * * *"
const UpdatesCiSchedule = "0 3 * * *"
const UptimeCiSchedule = "*/5 * * * *"

const DefaultRunner = "ubuntu-latest"
const DynamicRandomNumber = "$DynamicRandomNumber"
const RandomMinDefault = "0"
const RandomMaxDefault = "1000000"
const DynamicAlphanumericString = "$DynamicAlphanumericString"
const DynamicStringLengthDefault = "10"

const TcpPingDefaultCount = 5    // seconds
const TcpPingDefaultTimeout = 5  // seconds
const TcpPingDefaultInterval = 1 // seconds
const MaxResponseTime = 60000

const SecondWaitTime = time.Second
const ThirdWaitTime = time.Second * 10

const DefaultUp = "🟩"
const DefaultDegraded = "🟨"
const DefaultDown = "🟥"

var (
	AuthorName  = flag.String("authorName", "Uptime Check Bot", "Author name for commit")
	AuthorEmail = flag.String("authorEmail", "dty3152@gmail.com", "Author email for commit")
)
