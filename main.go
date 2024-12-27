package main

import (
	"grafana-dashboards-downloader/internal"
	"grafana-dashboards-downloader/internal/config"
	"log/slog"
	"os"
)

func loadConfigFromFile(logger *slog.Logger, configPath string) (*config.Config, error) {
	configFileContent, err := os.ReadFile(configPath)
	if err != nil {
		logger.Warn("couldn't load yaml config file", "error", err)
		return nil, err
	}

	config, err := config.LoadYamlConfig(configFileContent)
	if err != nil {
		logger.Warn("couldn't parse yaml config file", "error", err)
		return nil, err
	}

	return &config, err
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	configPath := os.Getenv("OPERATOR_CONFIG_PATH")

	config, err := loadConfigFromFile(logger, configPath)
	if err != nil {
		os.Exit(1)
	}
	internal.BuildAndRun(logger, *config)
}
