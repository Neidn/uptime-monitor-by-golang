package lib

import (
	"fmt"
	"github.com/Neidn/uptime-monitor-by-golang/config"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"net/http"
	"reflect"
)

func HealthCheck() bool {
	owner, repo := config.GetOwnerRepo()
	if owner == config.OwnerName && repo == config.RepositoryName {
		return true
	}
	var remoteConfig config.UptimeConfig
	var defaultConfig config.UptimeConfig
	defaultConfig.GetConfig()

	resp, err := http.Get(
		fmt.Sprintf(
			"https://raw.githubusercontent.com/%s/%s/HEAD/%s",
			config.OwnerName,
			config.RepositoryName,
			config.UptimeRcYaml,
		),
	)

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if err != nil {
		return false
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	err = yaml.Unmarshal(data, &remoteConfig)
	if err != nil {
		log.Printf("Error unmarshalling remote config: %s", err)
		return false
	}

	if reflect.DeepEqual(remoteConfig, defaultConfig) {
		log.Printf(`

[warn] > UPPTIME WARNING
[warn] > You should change your Upptime configuration (.upptimerc.yml)
[warn] > Upptime workflows will NOT work until you've added custom configuration

`)
		return false
	}

	return true
}
