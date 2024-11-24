package syncer

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCurrentDashboards__NoDashboards(t *testing.T) {
	currentDashboards := NewCurrentDashboards()
	newDashboard := Dashboard{"filename", "content"}

	require.Equal(t, true, currentDashboards.dashboardHasBeenUpdated(newDashboard))
}

func TestCurrentDashboards__SameDashboard(t *testing.T) {
	currentDashboards := NewCurrentDashboards()
	newDashboard := Dashboard{"filename", "content"}
	currentDashboards.saveDashboard(newDashboard)

	require.Equal(t, false, currentDashboards.dashboardHasBeenUpdated(newDashboard))
}

func TestCurrentDashboards__UpdatedDashboard(t *testing.T) {
	currentDashboards := NewCurrentDashboards()
	oldDashboard := Dashboard{"filename", "content"}
	newDashboard := Dashboard{"filename", "content-123"}
	currentDashboards.saveDashboard(oldDashboard)

	require.Equal(t, true, currentDashboards.dashboardHasBeenUpdated(newDashboard))
}
