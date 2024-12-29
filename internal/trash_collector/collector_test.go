package trash_collector

import (
	"errors"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func fileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func TestCollector_Integration(t *testing.T) {
	temporaryDirectory, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	defer os.RemoveAll(temporaryDirectory)

	collector := NewTrashCollector(temporaryDirectory)

	collector.AddKnownFile("file1.json")
	err = os.WriteFile(filepath.Join(temporaryDirectory, "file1.json"), []byte{}, 0644)
	require.NoError(t, err)

	collector.AddKnownFile("file2.json")
	err = os.WriteFile(filepath.Join(temporaryDirectory, "file2.json"), []byte{}, 0644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(temporaryDirectory, "file3.yml"), []byte{}, 0644)
	require.NoError(t, err)

	err = collector.PickUpTrash()
	require.NoError(t, err)

	require.True(t, fileExists(filepath.Join(temporaryDirectory, "file1.json")))
	require.True(t, fileExists(filepath.Join(temporaryDirectory, "file2.json")))
	require.False(t, fileExists(filepath.Join(temporaryDirectory, "file3.json")))
}
