package http_downloader

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"io"
	"net/http"
	"strconv"
)

var downloaderStatusCode = promauto.NewCounterVec(
	prometheus.CounterOpts{Name: "unraid_monitoring_operator_downloader_status_code"},
	[]string{"url", "status_code"})

func Download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	downloaderStatusCode.With(prometheus.Labels{"url": url, "status_code": strconv.Itoa(resp.StatusCode)}).Inc()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return content, fmt.Errorf("failed to fetch the file, status_code: %s", resp.Status)
	}

	return content, nil
}
