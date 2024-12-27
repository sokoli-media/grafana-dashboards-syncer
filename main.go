package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"grafana-dashboards-downloader/internal"
	"log/slog"
	"os"
	"strings"
	"syscall"
)

func parseDashboards(dashboards []string) (map[string]string, error) {
	parsed := make(map[string]string)
	for _, dashboard := range dashboards {
		filename, url, found := strings.Cut(dashboard, "=")
		if !found {
			return nil, errors.New("wrong number of parameters for dashboard: " + dashboard)
		}

		parsed[filename] = url
	}
	return parsed, nil
}

func loadConfigFromFile(configPath string) map[string]string {
	configFileContent, err := os.ReadFile(configPath)
	if err != nil {
		return map[string]string{}
	}

	config, err := internal.LoadYamlConfig(configFileContent)
	if err != nil {
		return map[string]string{}
	}

	var oldStyleConfig map[string]string
	for _, dashboard := range config.Grafana.Dashboards {
		filenameBase := md5.Sum([]byte(dashboard.HTTPSource.Url))
		filename := fmt.Sprintf("%s.json", filenameBase)
		oldStyleConfig[filename] = dashboard.HTTPSource.Url
	}

	return oldStyleConfig
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	configPath := os.Getenv("OPERATOR_CONFIG_PATH")

	var dashboards []string

	var rootCmd = &cobra.Command{
		Use: "grafana-dashboards-syncer",
	}
	rootCmd.PersistentFlags().StringSliceVarP(
		&dashboards,
		"dashboard",
		"d",
		[]string{},
		"Specify dashboard(s), format filename=url",
	)

	if err := rootCmd.Execute(); err != nil {
		logger.Error("couldn't start due to configuration error", "error", err)
	}

	mappedDashboards := loadConfigFromFile(configPath)
	if len(mappedDashboards) == 0 {
		mappedDashboards, err := parseDashboards(dashboards)
		if err != nil {
			logger.Error("couldn't parse dashboards", "error", err)
			syscall.Exit(1)
		}
		if len(mappedDashboards) < 1 {
			logger.Error("you must provide at least 1 dashboard")
			syscall.Exit(1)
		}
	}

	internal.BuildAndRun(logger, mappedDashboards)
}
