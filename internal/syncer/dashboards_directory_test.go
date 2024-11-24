package syncer

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
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

func TestDashboardsDirectory__ClearCacheWhenAddingDashboard(t *testing.T) {
	dirPath, err := os.MkdirTemp("", "testDir")
	defer os.RemoveAll(dirPath)
	require.NoError(t, err)

	directory := NewDashboardsDirectory(dirPath)

	// we want to create some dashboard to enable caching...
	err = directory.saveDashboard(Dashboard{
		filename:  "filename1",
		dashboard: "",
	})
	require.NoError(t, err)

	dashboards, err := directory.listDashboards()
	require.NoError(t, err)
	require.Len(t, dashboards, 1)

	// and then we add a new file without invalidating our cache
	err = os.WriteFile(filepath.Join(dirPath, "filename2"), []byte(""), 0644)
	require.NoError(t, err)

	// we want to double-check we're using cache here, so that...
	dashboards, err = directory.listDashboards()
	require.NoError(t, err)
	require.Len(t, dashboards, 1)

	// ...when we save a new dashboard...
	err = directory.saveDashboard(Dashboard{
		filename:  "filename3",
		dashboard: "",
	})
	require.NoError(t, err)

	// cache is now invalidated, our old file is now also returned!
	dashboards, err = directory.listDashboards()
	require.NoError(t, err)
	require.Len(t, dashboards, 3)
}

func TestDashboardsDirectory__ClearCacheWhenRemovingDashboard(t *testing.T) {
	dirPath, err := os.MkdirTemp("", "testDir")
	defer os.RemoveAll(dirPath)
	require.NoError(t, err)

	directory := NewDashboardsDirectory(dirPath)

	// we want to create some dashboard to enable caching...
	err = directory.saveDashboard(Dashboard{
		filename:  "filename1",
		dashboard: "",
	})
	require.NoError(t, err)

	dashboards, err := directory.listDashboards()
	require.NoError(t, err)
	require.Len(t, dashboards, 1)
	require.Equal(t, "filename1", dashboards[0])

	// and then we add a new file without invalidating our cache
	err = os.WriteFile(filepath.Join(dirPath, "filename2"), []byte(""), 0644)
	require.NoError(t, err)

	// we want to double-check we're using cache here, so that...
	dashboards, err = directory.listDashboards()
	require.NoError(t, err)
	require.Len(t, dashboards, 1)
	require.Equal(t, "filename1", dashboards[0])

	// ...when we remove a dashboard...
	err = directory.removeDashboard("filename1")
	require.NoError(t, err)

	// cache is now invalidated, our secret file is the only one returned!
	dashboards, err = directory.listDashboards()
	require.NoError(t, err)
	require.Len(t, dashboards, 1)
	require.Equal(t, "filename2", dashboards[0])
}
