package grafana_syncer

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"os"
	"path/filepath"
)

var directoryListFromDrive = promauto.NewCounter(
	prometheus.CounterOpts{Name: "grafana_dashboards_syncer_dashboard_directory_list_from_drive"})
var directoryListFromCache = promauto.NewCounter(
	prometheus.CounterOpts{Name: "grafana_dashboards_syncer_dashboard_directory_list_from_cache"})
var directorySaveDashboard = promauto.NewCounter(
	prometheus.CounterOpts{Name: "grafana_dashboards_syncer_dashboard_directory_save"})
var directoryRemoveDashboard = promauto.NewCounter(
	prometheus.CounterOpts{Name: "grafana_dashboards_syncer_dashboard_directory_remove"})

func NewDashboardsDirectory(dashboardsDirectory string) *DashboardsDirectory {
	var existingFiles []string
	return &DashboardsDirectory{
		directoryPath: dashboardsDirectory,
		existingFiles: &existingFiles,
	}
}

type DashboardsDirectory struct {
	directoryPath string
	existingFiles *[]string
}

func (d DashboardsDirectory) saveDashboard(dashboard Dashboard) error {
	*d.existingFiles = nil
	directorySaveDashboard.Inc()

	fullPath := filepath.Join(d.directoryPath, dashboard.filename)
	return os.WriteFile(fullPath, []byte(dashboard.dashboard), 0644)
}

func (d DashboardsDirectory) listDashboards() ([]string, error) {
	if len(*d.existingFiles) == 0 {
		directoryListFromDrive.Inc()

		entries, err := os.ReadDir(d.directoryPath)
		if err != nil {
			return nil, err
		}

		var dashboards []string
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			dashboards = append(dashboards, entry.Name())
		}

		*d.existingFiles = dashboards
	} else {
		directoryListFromCache.Inc()
	}

	return *d.existingFiles, nil
}

func (d DashboardsDirectory) removeDashboard(filename string) error {
	*d.existingFiles = nil
	directoryRemoveDashboard.Inc()

	fullPath := filepath.Join(d.directoryPath, filename)
	return os.Remove(fullPath)
}
