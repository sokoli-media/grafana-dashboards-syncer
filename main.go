package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"grafana-dashboards-downloader/internal"
	"log/slog"
	"os"
)

func loadConfigFromFile(logger *slog.Logger, configPath string) map[string]string {
	configFileContent, err := os.ReadFile(configPath)
	if err != nil {
		logger.Warn("couldn't load yaml config file", "error", err)
		return map[string]string{}
	}

	config, err := internal.LoadYamlConfig(configFileContent)
	if err != nil {
		logger.Warn("couldn't parse yaml config file", "error", err)
		return map[string]string{}
	}

	oldStyleConfig := map[string]string{}
	for _, dashboard := range config.Grafana.Dashboards {
		md5sum := md5.New()
		md5sum.Write([]byte(dashboard.HTTPSource.Url))
		filenameBase := hex.EncodeToString(md5sum.Sum(nil))

		filename := fmt.Sprintf("%s.json", filenameBase)
		oldStyleConfig[filename] = dashboard.HTTPSource.Url
	}

	return oldStyleConfig
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	configPath := os.Getenv("OPERATOR_CONFIG_PATH")

	mappedDashboards := loadConfigFromFile(logger, configPath)
	internal.BuildAndRun(logger, mappedDashboards)
}
