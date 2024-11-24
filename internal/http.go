package internal

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"grafana-dashboards-downloader/internal/syncer"
	"log/slog"
	"net/http"
)

func runHTTPServer(logger *slog.Logger) {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		logger.Error("failed to run http server", "error", err)
		return
	}
}

func Run(logger *slog.Logger, dashboards map[string]string) {
	go syncer.BackgroundSyncingDaemon(logger, dashboards)

	runHTTPServer(logger)
}
