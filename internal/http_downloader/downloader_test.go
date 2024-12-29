package http_downloader

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownloader__ReturnFile(t *testing.T) {
	expectedPath := "/dashboard.json"
	expectedContent := "some json or other value that is returned"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != expectedPath {
			t.Errorf("Expected to request '%s', got: %s", expectedPath, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(expectedContent))
	}))
	defer server.Close()

	content, err := Download(fmt.Sprintf("%s/dashboard.json", server.URL))
	require.NoError(t, err)
	require.Equal(t, []byte(expectedContent), content)
}

func TestDownloader__FileNotReturned(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500 - Something bad happened!"))
	}))
	defer server.Close()

	content, err := Download(fmt.Sprintf("%s/dashboard.json", server.URL))
	require.Error(t, err)
	require.Equal(t, []byte("500 - Something bad happened!"), content)
}

func TestDownloader__ServerDoesntWork(t *testing.T) {
	content, err := Download("http://this-url-is-fake/dashboard.json")
	require.Error(t, err)
	require.Equal(t, []byte(nil), content)
}
