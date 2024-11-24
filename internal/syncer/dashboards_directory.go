package syncer

import (
	"os"
	"path/filepath"
)

func NewDashboardsDirectory(dashboardsDirectory string) *DashboardsDirectory {
	return &DashboardsDirectory{
		directoryPath: dashboardsDirectory,
	}
}

type DashboardsDirectory struct {
	directoryPath string
	existingFiles *[]string
}

func (d DashboardsDirectory) saveDashboard(dashboard Dashboard) error {
	d.existingFiles = nil

	fullPath := filepath.Join(d.directoryPath, dashboard.filename)
	return os.WriteFile(fullPath, []byte(dashboard.dashboard), 0644)
}

func (d DashboardsDirectory) listDashboards() ([]string, error) {
	if len(*d.existingFiles) == 0 {
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

		d.existingFiles = &dashboards
	}

	return *d.existingFiles, nil
}

func (d DashboardsDirectory) removeDashboard(filename string) error {
	d.existingFiles = nil

	fullPath := filepath.Join(d.directoryPath, filename)
	return os.Remove(fullPath)
}
