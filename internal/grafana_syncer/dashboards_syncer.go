package grafana_syncer

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log/slog"
)

var dashboardUpdated = promauto.NewCounterVec(
	prometheus.CounterOpts{Name: "grafana_dashboards_syncer_dashboard_updated"},
	[]string{"dashboard"})
var dashboardNotUpdated = promauto.NewCounterVec(
	prometheus.CounterOpts{Name: "grafana_dashboards_syncer_dashboard_not_updated"},
	[]string{"dashboard"})
var oldDashboardCleanedUp = promauto.NewCounterVec(
	prometheus.CounterOpts{Name: "grafana_dashboards_syncer_old_dashboard_cleaned_up"},
	[]string{"dashboard"})

type DashboardsSyncer struct {
	logger            *slog.Logger
	downloader        *DashboardsDownloader
	currentDashboards *CurrentDashboards
	directory         *DashboardsDirectory
}

func (d DashboardsSyncer) downloadDashboards() DownloadedDashboards {
	d.logger.Info("downloading dashboards")
	downloadedFiles := NewDownloadedDashboards()

	for _, dashboard := range d.downloader.downloadAllDashboards() {
		if d.currentDashboards.dashboardHasBeenUpdated(dashboard) {
			dashboardUpdated.WithLabelValues(dashboard.filename).Inc()
			err := d.directory.saveDashboard(dashboard)
			if err != nil {
				d.logger.Error("saving dashboard failed", "error", err)
				continue
			}

			d.currentDashboards.saveDashboard(dashboard)
			d.logger.Info("saved new or updated dashboard", "filename", dashboard.filename)
		} else {
			dashboardNotUpdated.WithLabelValues(dashboard.filename).Inc()
			d.logger.Info("dashboard is the same, ignoring", "filename", dashboard.filename)
		}

		downloadedFiles.markAsDownloaded(dashboard)
	}
	return downloadedFiles
}

func (d DashboardsSyncer) cleanUpDashboards(downloadedDashboards DownloadedDashboards) {
	d.logger.Info("checking if there are old dashboards to be removed")
	existingDashboards, err := d.directory.listDashboards()
	if err != nil {
		d.logger.Error("couldn't list dashboards", "error", err)
		return
	}

	for _, filename := range existingDashboards {
		if !downloadedDashboards.hasBeenDownloaded(filename) {
			oldDashboardCleanedUp.WithLabelValues(filename).Inc()

			d.logger.Info("removing old dashboard", "filename", filename)
			err = d.directory.removeDashboard(filename)
			if err != nil {
				d.logger.Error("couldn't remove old dashboard", "file", filename)
			}
		}
	}
}
