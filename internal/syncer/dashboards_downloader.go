package syncer

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"grafana-dashboards-downloader/internal/config"
	"io"
	"log/slog"
	"net/http"
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

		dashboardBody, err := d.downloadDashboard(dashboard.HTTPSource.Url)
		if err != nil {
			labels := prometheus.Labels{"dashboard": filename, "reason": fmt.Sprintf("%s", err)}
			downloadsFailure.With(labels).Inc()

			d.logger.Error("couldn't download dashboard", "url", dashboard.HTTPSource.Url, "error", err)
			continue
		}
		downloadsSuccess.With(prometheus.Labels{"dashboard": filename}).Inc()

		dashboard := Dashboard{
			filename:  filename,
			dashboard: dashboardBody,
		}
		downloaded = append(downloaded, dashboard)
	}
	return downloaded
}

func (d DashboardsDownloader) downloadDashboard(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch the file, status_code: %s", resp.Status)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
