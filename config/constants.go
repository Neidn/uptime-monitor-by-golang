package config

import "time"

const OwnerName = "Neidn"
const MonitorRepositoryName = "uptime-monitor-by-golang"
const RepositoryName = "uptime"

const GithubActionsInputKey = "command"

const SecretContextKey = "SECRETS_CONTEXT"
const GithubRepoKey = "GITHUB_REPOSITORY"

const ConfigYaml = ".uptimerc.yml"
const HistoryYamlDir = "history"

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

const TcpPingDefaultCount = 5
const MaxResponseTime = 60000

const SecondWaitTime = time.Second
const ThirdWaitTime = time.Second * 10
