package main

var release string

func GetUptimeMonitorVersion() string {
	if release != "" {
		return release
	}

	return "v1.0.0"
}
