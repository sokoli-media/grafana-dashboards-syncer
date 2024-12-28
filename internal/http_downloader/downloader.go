package http_downloader

import (
	"fmt"
	"io"
	"net/http"
)

func Download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return content, fmt.Errorf("failed to fetch the file, status_code: %s", resp.Status)
	}

	return content, nil
}
