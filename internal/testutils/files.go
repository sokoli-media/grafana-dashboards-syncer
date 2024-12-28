package testutils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func GetHashedFilename(url string, extension string) string {
	md5sum := md5.New()
	md5sum.Write([]byte(url))
	filenameBase := hex.EncodeToString(md5sum.Sum(nil))
	return fmt.Sprintf("%s.%s", filenameBase, extension)
}

func FileExists(directory string, filename string) bool {
	if _, err := os.Stat(filepath.Join(directory, filename)); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func GetFileModificationTime(t *testing.T, directory string, filename string) time.Time {
	fileInfo, err := os.Stat(filepath.Join(directory, filename))
	require.NoError(t, err)

	return fileInfo.ModTime()
}
