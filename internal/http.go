package internal

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"os"
	"unraid-monitoring-operator/internal/config"
	"unraid-monitoring-operator/internal/grafana_syncer"
	"unraid-monitoring-operator/internal/prometheus_syncer"
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

func BuildAndRun(logger *slog.Logger, config config.Config) {
	dashboardsDirectory := getEnv("GRAFANA_DASHBOARDS_DIRECTORY", "./dashboards/")
	go grafana_syncer.BackgroundSyncingDaemon(logger, config.Grafana.Dashboards, dashboardsDirectory)

	prometheusSyncer := prometheus_syncer.NewPrometheusSyncer(logger, config.Prometheus)
	go prometheus_syncer.RunBackgroundSyncingDaemon(logger, prometheusSyncer)

	runHTTPServer(logger)
}
