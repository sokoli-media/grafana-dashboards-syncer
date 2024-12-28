package trash_collector

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
)

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
		panic("filename must be a basename, not full path, this should never happen")
	}

	c.knownFiles = append(c.knownFiles, filename)
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
			if err != nil {
				return err
			}
		}
	}

	return nil
}
