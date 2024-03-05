package workflows

import (
	"fmt"
	"github.com/Neidn/uptime-monitor-by-golang/config"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func ReadSiteHistory(slugName string) (*SiteHistory, error) {
	siteHistoryFile, err := os.ReadFile(filepath.Join(config.HistoryYamlDir, slugName+".yml"))
	if err != nil {
		return nil, err
	}

	siteHistory := SiteHistory{}
	err = yaml.Unmarshal(siteHistoryFile, &siteHistory)
	if err != nil {
		return nil, err
	}

	return &siteHistory, nil
}

func WriteSiteHistory(slugName string, siteHistory *SiteHistory) error {
	historyBody := fmt.Sprintf(`url: %s
status: %s
code: %d
responseTime: %d
lastUpdated: %s
startTime: %s
generator: %s
`, siteHistory.Url, siteHistory.Status, siteHistory.Code, siteHistory.ResponseTime, siteHistory.LastUpdated, siteHistory.StartTime, siteHistory.Generator)

	_ = os.WriteFile(filepath.Join(config.HistoryYamlDir, slugName+".yml"), []byte(historyBody), 0644)

	return nil
}
