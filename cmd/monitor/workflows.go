package main

var release string

func GetUptimeMonitorVersion(token string) string {
	if release != "" {
		return release
	}

	_, _ = GithubClient(token)

	return "v1.0.0"
}
