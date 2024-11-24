package syncer

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestDashboardsDirectory__EmptyDirectory(t *testing.T) {
	dirPath, err := os.MkdirTemp("", "testDir")
	defer os.RemoveAll(dirPath)
	require.NoError(t, err)

	directory := NewDashboardsDirectory(dirPath)
	dashboards, err := directory.listDashboards()
	require.NoError(t, err)
	require.Empty(t, dashboards)
}

func TestDashboardsDirectory__MultipleFiles(t *testing.T) {
	dirPath, err := os.MkdirTemp("", "testDir")
	defer os.RemoveAll(dirPath)
	require.NoError(t, err)

	directory := NewDashboardsDirectory(dirPath)

	filenames := []string{"batman", "city", "gotham"}
	for _, filename := range filenames {
		err = directory.saveDashboard(Dashboard{
			filename:  filename,
			dashboard: "",
		})
		require.NoError(t, err)
	}

	dashboards, err := directory.listDashboards()
	require.NoError(t, err)
	require.Equal(t, filenames, dashboards)
}

func TestDashboardsDirectory__MultipleFilesUsingCache(t *testing.T) {
	dirPath, err := os.MkdirTemp("", "testDir")
	defer os.RemoveAll(dirPath)
	require.NoError(t, err)

	directory := NewDashboardsDirectory(dirPath)

	filenames := []string{"batman", "city", "gotham"}
	for _, filename := range filenames {
		err = directory.saveDashboard(Dashboard{
			filename:  filename,
			dashboard: "",
		})
		require.NoError(t, err)
	}

	dashboards, err := directory.listDashboards()
	require.NoError(t, err)
	require.Equal(t, filenames, dashboards)

	_ = os.RemoveAll(dirPath)

	dashboards, err = directory.listDashboards()
	require.NoError(t, err)
	require.Equal(t, filenames, dashboards)
}
