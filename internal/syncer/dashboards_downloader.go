package syncer

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

func NewDashboardsDownloader(logger *slog.Logger, dashboards map[string]string) *DashboardsDownloader {
	return &DashboardsDownloader{
		logger:     logger,
		dashboards: dashboards,
	}
}

type DashboardsDownloader struct {
	logger     *slog.Logger
	dashboards map[string]string
}

func (d DashboardsDownloader) downloadAllDashboards() []Dashboard {
	var downloaded []Dashboard
	for filename, url := range d.dashboards {
		dashboardBody, err := d.downloadDashboard(url)
		if err != nil {
			d.logger.Error("couldn't download dashboard", "url", url, "error", err)
			continue
		}
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
