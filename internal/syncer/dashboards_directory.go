package syncer

import (
	"os"
	"path/filepath"
)

func getEnv(variableName string, defaultValue string) string {
	value := os.Getenv(variableName)
	if value == "" {
		return defaultValue
	}
	return value
}

func NewDashboardsDirectory() DashboardsDirectory {
	return DashboardsDirectory{
		directoryPath: getEnv("GRAFANA_DASHBOARDS_DIRECTORY", "./dashboards/"),
	}
}

type DashboardsDirectory struct {
	directoryPath string
}

func (d DashboardsDirectory) saveDashboard(dashboard Dashboard) error {
	fullPath := filepath.Join(d.directoryPath, dashboard.filename)
	return os.WriteFile(fullPath, []byte(dashboard.dashboard), 0644)
}

func (d DashboardsDirectory) listDashboards() ([]string, error) {
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

	return dashboards, nil
}

func (d DashboardsDirectory) removeDashboard(filename string) error {
	fullPath := filepath.Join(d.directoryPath, filename)
	return os.Remove(fullPath)
}
