package syncer

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/require"
	"grafana-dashboards-downloader/internal/config"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestDashboardsDownloader__SingleDashboard(t *testing.T) {
	expectedDashboard := Dashboard{
		filename:  "test.json",
		dashboard: "some json that is returned",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/dashboard.json" {
			t.Errorf("Expected to request '/dashboard.json', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(expectedDashboard.dashboard))
	}))
	defer server.Close()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	dashboards := []config.GrafanaDashboardsConfig{
		{
			HTTPSource: config.GrafanaHTTPSourceConfig{
				Url: fmt.Sprintf("%s/dashboard.json", server.URL),
			},
		},
	}
	downloader := NewDashboardsDownloader(logger, dashboards)

	downloadedDashboards := downloader.downloadAllDashboards()

	require.NotEmpty(t, downloadedDashboards)
	require.Len(t, downloadedDashboards, 1)
	require.Equal(t, downloadedDashboards[0].dashboard, expectedDashboard.dashboard)
}

func TestDashboardsDownloader__MultipleDashboards(t *testing.T) {
	expectedContent := map[string]string{
		"/dashboard1.json": "some json that is returned for dashboard1.json",
		"/dashboard2.json": "another json that is returned for dashboard2.json",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/dashboard1.json" && r.URL.Path != "/dashboard2.json" {
			t.Errorf("Expected to request '/dashboard1.json' or '/dashboard2.json', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(expectedContent[r.URL.Path]))
	}))
	defer server.Close()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	dashboards := []config.GrafanaDashboardsConfig{
		{
			HTTPSource: config.GrafanaHTTPSourceConfig{
				Url: fmt.Sprintf("%s/dashboard1.json", server.URL),
			},
		},
		{
			HTTPSource: config.GrafanaHTTPSourceConfig{
				Url: fmt.Sprintf("%s/dashboard2.json", server.URL),
			},
		},
	}
	downloader := NewDashboardsDownloader(logger, dashboards)

	downloadedDashboards := downloader.downloadAllDashboards()

	require.NotEmpty(t, downloadedDashboards)
	require.Len(t, downloadedDashboards, 2)

	d1md5sum := md5.New()
	d1md5sum.Write([]byte(dashboards[0].HTTPSource.Url))
	d1filenameBase := hex.EncodeToString(d1md5sum.Sum(nil))
	d1filename := fmt.Sprintf("%s.json", d1filenameBase)

	d2md5sum := md5.New()
	d2md5sum.Write([]byte(dashboards[1].HTTPSource.Url))
	d2filenameBase := hex.EncodeToString(d2md5sum.Sum(nil))
	d2filename := fmt.Sprintf("%s.json", d2filenameBase)

	require.Contains(t, downloadedDashboards, Dashboard{
		filename:  d1filename,
		dashboard: expectedContent["/dashboard1.json"],
	})
	require.Contains(t, downloadedDashboards, Dashboard{
		filename:  d2filename,
		dashboard: expectedContent["/dashboard2.json"],
	})
}

func TestDashboardsDownloader__FirstDashboardUrlNotWorking(t *testing.T) {
	expectedContent := map[string]string{
		"/dashboard1.json": "some json that is returned for dashboard1.json",
		"/dashboard2.json": "another json that is returned for dashboard2.json",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/dashboard1.json" && r.URL.Path != "/dashboard2.json" {
			t.Errorf("Expected to request '/dashboard1.json' or '/dashboard2.json', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(expectedContent[r.URL.Path]))
	}))
	defer server.Close()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	dashboards := []config.GrafanaDashboardsConfig{
		{
			HTTPSource: config.GrafanaHTTPSourceConfig{
				Url: "http://not-existing-url/dashboard1.json",
			},
		},
		{
			HTTPSource: config.GrafanaHTTPSourceConfig{
				Url: fmt.Sprintf("%s/dashboard2.json", server.URL),
			},
		},
	}
	downloader := NewDashboardsDownloader(logger, dashboards)

	downloadedDashboards := downloader.downloadAllDashboards()

	require.NotEmpty(t, downloadedDashboards)
	require.Len(t, downloadedDashboards, 1)
	require.Equal(t, downloadedDashboards[0].dashboard, expectedContent["/dashboard2.json"])
}
