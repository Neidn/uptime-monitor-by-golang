package config

const SecretContextKey = "SECRETS_CONTEXT"
const GithubRepoKey = "GITHUB_REPOSITORY"

const ConfigFile = ".github/uptime.yml"

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
