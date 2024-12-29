package testutils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func SetupFakeServer(t *testing.T, expectedPath string, content string) *FakeServer {
	f := &FakeServer{
		t:            t,
		ExpectedPath: expectedPath,
		Response:     content,
	}
	f.Init()
	return f
}

type FakeServer struct {
	t            *testing.T
	Server       *httptest.Server
	URL          string
	ExpectedPath string
	Response     string
}

func (f *FakeServer) Init() {
	f.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != f.ExpectedPath {
			f.t.Errorf("Expected to request '%s', got: %s", f.ExpectedPath, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(f.Response))
	}))
	f.URL = f.Server.URL
}

func (f *FakeServer) Close() {
	f.Server.Close()
}
