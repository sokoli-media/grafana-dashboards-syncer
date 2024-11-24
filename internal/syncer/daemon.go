package syncer

import (
	"log/slog"
	"time"
)

func BackgroundSyncingDaemon(logger *slog.Logger, dashboards map[string]string, dashboardsDirectory string) {
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
