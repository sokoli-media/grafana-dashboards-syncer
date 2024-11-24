package internal

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"grafana-dashboards-downloader/internal/syncer"
	"log/slog"
	"net/http"
	"os"
)

func getEnv(variableName string, defaultValue string) string {
	value := os.Getenv(variableName)
	if value == "" {
		return defaultValue
	}
	return value
}

func runHTTPServer(logger *slog.Logger) {
	http.HandleFunc("/dashboard.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/dashboards/dashboard.json")
	})
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		logger.Error("failed to run http server", "error", err)
		return
	}
}

func BuildAndRun(logger *slog.Logger, dashboards map[string]string) {
	dashboardsDirectory := getEnv("GRAFANA_DASHBOARDS_DIRECTORY", "./dashboards/")
	go syncer.BackgroundSyncingDaemon(logger, dashboards, dashboardsDirectory)

	runHTTPServer(logger)
}
