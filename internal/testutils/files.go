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

func LoadFile(t *testing.T, directory string, filename string) string {
	dat, err := os.ReadFile(filepath.Join(directory, filename))
	require.NoError(t, err)

	return string(dat)
}

func WriteFile(t *testing.T, directory string, filename string, content string) {
	err := os.WriteFile(filepath.Join(directory, filename), []byte(content), 0644)
	require.NoError(t, err)
}
