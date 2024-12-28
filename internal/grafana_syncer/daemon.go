package grafana_syncer

import (
	"grafana-dashboards-downloader/internal/config"
	"log/slog"
	"time"
)

func BackgroundSyncingDaemon(logger *slog.Logger, dashboards []config.GrafanaDashboardsConfig, dashboardsDirectory string) {
	downloader := NewDashboardsDownloader(logger, dashboards)
	currentDashboards := NewCurrentDashboards()
	directory := NewDashboardsDirectory(dashboardsDirectory)
	dashboardsSyncer := DashboardsSyncer{
		logger:            logger,
		downloader:        downloader,
		currentDashboards: currentDashboards,
		directory:         directory,
	}

	for {
		logger.Info("starting syncing dashboards")
		downloadedFiles := dashboardsSyncer.downloadDashboards()
		dashboardsSyncer.cleanUpDashboards(downloadedFiles)

		time.Sleep(30 * time.Second)
	}
}
