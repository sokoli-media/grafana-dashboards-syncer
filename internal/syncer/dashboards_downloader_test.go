package syncer

import (
	"fmt"
	"github.com/stretchr/testify/require"
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
	dashboards := map[string]string{
		expectedDashboard.filename: fmt.Sprintf("%s/dashboard.json", server.URL),
	}
	downloader := NewDashboardsDownloader(logger, dashboards)

	downloadedDashboards := downloader.downloadAllDashboards()

	require.NotEmpty(t, downloadedDashboards)
	require.Len(t, downloadedDashboards, 1)
	require.Equal(t, downloadedDashboards[0], expectedDashboard)
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
	dashboards := map[string]string{
		"d1.json": fmt.Sprintf("%s/dashboard1.json", server.URL),
		"d2.json": fmt.Sprintf("%s/dashboard2.json", server.URL),
	}
	downloader := NewDashboardsDownloader(logger, dashboards)

	downloadedDashboards := downloader.downloadAllDashboards()

	require.NotEmpty(t, downloadedDashboards)
	require.Len(t, downloadedDashboards, 2)

	require.Contains(t, downloadedDashboards, Dashboard{
		filename:  "d1.json",
		dashboard: expectedContent["/dashboard1.json"],
	})
	require.Contains(t, downloadedDashboards, Dashboard{
		filename:  "d2.json",
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
	dashboards := map[string]string{
		"d1.json": "http://not-existing-url/dashboard1.json",
		"d2.json": fmt.Sprintf("%s/dashboard2.json", server.URL),
	}
	downloader := NewDashboardsDownloader(logger, dashboards)

	downloadedDashboards := downloader.downloadAllDashboards()

	require.NotEmpty(t, downloadedDashboards)
	require.Len(t, downloadedDashboards, 1)
	require.Equal(t, downloadedDashboards[0], Dashboard{
		filename:  "d2.json",
		dashboard: expectedContent["/dashboard2.json"],
	})
}
