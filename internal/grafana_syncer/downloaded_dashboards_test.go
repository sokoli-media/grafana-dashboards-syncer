package grafana_syncer

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDownloadedDashboards__NoDownloads(t *testing.T) {
	downloadedDashboards := NewDownloadedDashboards()

	require.Equal(t, false, downloadedDashboards.hasBeenDownloaded("filename"))
}

func TestDownloadedDashboards__AlreadyDownloaded(t *testing.T) {
	downloadedDashboards := NewDownloadedDashboards()
	dashboard := Dashboard{
		filename:  "filename",
		dashboard: "content",
	}
	downloadedDashboards.markAsDownloaded(dashboard)

	require.Equal(t, true, downloadedDashboards.hasBeenDownloaded(dashboard.filename))
}
