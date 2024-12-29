package trash_collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var addedFile = promauto.NewCounterVec(
	prometheus.CounterOpts{Name: "unraid_monitoring_operator_trash_collector_added_file"},
	[]string{"directory"})
var removedFile = promauto.NewCounterVec(
	prometheus.CounterOpts{Name: "unraid_monitoring_operator_trash_collector_removed_file"},
	[]string{"directory"})

func NewTrashCollector(directory string) *DiskCollector {
	return &DiskCollector{
		directory:  directory,
		knownFiles: []string{},
	}
}

type DiskCollector struct {
	directory  string
	knownFiles []string
}

func (c *DiskCollector) AddKnownFile(filename string) {
	if strings.Contains(filename, "/") {
		// better panic than sorry
		panic("filename must be a basename, not full path, this should never happen")
	}

	c.knownFiles = append(c.knownFiles, filename)
	addedFile.With(prometheus.Labels{"directory": c.directory}).Inc()
}

func (c *DiskCollector) PickUpTrash() error {
	files, err := os.ReadDir(c.directory)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !slices.Contains(c.knownFiles, file.Name()) {
			fullPath := filepath.Join(c.directory, file.Name())
			err := os.Remove(fullPath)
			removedFile.With(prometheus.Labels{"directory": c.directory}).Inc()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
