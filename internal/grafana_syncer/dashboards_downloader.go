package grafana_syncer

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log/slog"
	"unraid-monitoring-operator/internal/config"
	"unraid-monitoring-operator/internal/http_downloader"
)

var downloadsFailure = promauto.NewCounterVec(
	prometheus.CounterOpts{Name: "grafana_dashboards_syncer_download_failure"},
	[]string{"dashboard", "reason"})

var downloadsSuccess = promauto.NewCounterVec(
	prometheus.CounterOpts{Name: "grafana_dashboards_syncer_download_success"},
	[]string{"dashboard"})

func NewDashboardsDownloader(logger *slog.Logger, dashboards []config.GrafanaDashboardsConfig) *DashboardsDownloader {
	return &DashboardsDownloader{
		logger:     logger,
		dashboards: dashboards,
	}
}

type DashboardsDownloader struct {
	logger     *slog.Logger
	dashboards []config.GrafanaDashboardsConfig
}

func (d DashboardsDownloader) downloadAllDashboards() []Dashboard {
	var downloaded []Dashboard
	for _, dashboard := range d.dashboards {
		md5sum := md5.New()
		md5sum.Write([]byte(dashboard.HTTPSource.Url))
		filenameBase := hex.EncodeToString(md5sum.Sum(nil))
		filename := fmt.Sprintf("%s.json", filenameBase)

		dashboardBody, err := http_downloader.Download(dashboard.HTTPSource.Url)
		if err != nil {
			labels := prometheus.Labels{"dashboard": filename, "reason": fmt.Sprintf("%s", err)}
			downloadsFailure.With(labels).Inc()

			d.logger.Error("couldn't download dashboard", "url", dashboard.HTTPSource.Url, "error", err)
			continue
		}
		downloadsSuccess.With(prometheus.Labels{"dashboard": filename}).Inc()

		dashboard := Dashboard{
			filename:  filename,
			dashboard: string(dashboardBody),
		}
		downloaded = append(downloaded, dashboard)
	}
	return downloaded
}
