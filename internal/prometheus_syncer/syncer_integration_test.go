package prometheus_syncer

import (
	"fmt"
	"github.com/stretchr/testify/require"
	config2 "grafana-dashboards-downloader/internal/config"
	"grafana-dashboards-downloader/internal/testutils"
	"os"
	"path/filepath"
	"testing"
)

func Test_AddSingleFileToEmptyDirectory(t *testing.T) {
	temporaryDirectory, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	defer os.RemoveAll(temporaryDirectory)

	server := testutils.SetupFakeServer(t, "/1.yml", "1.yml content")
	fileUrl := fmt.Sprintf("%s/1.yml", server.URL)
	defer server.Close()

	config := config2.PrometheusConfig{
		PrometheusRulesPath: temporaryDirectory,
		PrometheusRules: []config2.PrometheusRuleConfig{
			{
				SourceType: "http",
				HTTPSource: config2.HTTPSourceConfig{Url: fileUrl},
			},
		},
	}

	NewPrometheusSyncer(testutils.LoggerForTesting, config).Sync()

	savedContent, err := os.ReadFile(
		filepath.Join(temporaryDirectory, testutils.GetHashedFilename(fileUrl, "yml")))
	require.NoError(t, err)
	require.Equal(t, "1.yml content", string(savedContent))
}

func Test_AddMultipleFilesToEmptyDirectory(t *testing.T) {
	temporaryDirectory, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	defer os.RemoveAll(temporaryDirectory)

	server1 := testutils.SetupFakeServer(t, "/1.yml", "1.yml content")
	file1Url := fmt.Sprintf("%s/1.yml", server1.URL)
	defer server1.Close()

	server2 := testutils.SetupFakeServer(t, "/2.yml", "2.yml content")
	file2Url := fmt.Sprintf("%s/2.yml", server2.URL)
	defer server2.Close()

	config := config2.PrometheusConfig{
		PrometheusRulesPath: temporaryDirectory,
		PrometheusRules: []config2.PrometheusRuleConfig{
			{
				SourceType: "http",
				HTTPSource: config2.HTTPSourceConfig{Url: file1Url},
			},
			{
				SourceType: "http",
				HTTPSource: config2.HTTPSourceConfig{Url: file2Url},
			},
		},
	}

	NewPrometheusSyncer(testutils.LoggerForTesting, config).Sync()

	require.True(t, testutils.FileExists(temporaryDirectory, testutils.GetHashedFilename(file1Url, "yml")))
	require.True(t, testutils.FileExists(temporaryDirectory, testutils.GetHashedFilename(file2Url, "yml")))
}

func Test_AddMultipleFilesToExistingFiles(t *testing.T) {
	temporaryDirectory, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	defer os.RemoveAll(temporaryDirectory)

	server1 := testutils.SetupFakeServer(t, "/1.yml", "1.yml content")
	file1Url := fmt.Sprintf("%s/1.yml", server1.URL)
	defer server1.Close()

	server2 := testutils.SetupFakeServer(t, "/2.yml", "2.yml content")
	file2Url := fmt.Sprintf("%s/2.yml", server2.URL)
	defer server2.Close()

	server3 := testutils.SetupFakeServer(t, "/3.yml", "3.yml content")
	file3Url := fmt.Sprintf("%s/3.yml", server3.URL)
	defer server3.Close()

	server4 := testutils.SetupFakeServer(t, "/4.yml", "4.yml content")
	file4Url := fmt.Sprintf("%s/4.yml", server4.URL)
	defer server4.Close()

	config := config2.PrometheusConfig{
		PrometheusRulesPath: temporaryDirectory,
		PrometheusRules: []config2.PrometheusRuleConfig{
			{
				SourceType: "http",
				HTTPSource: config2.HTTPSourceConfig{Url: file1Url},
			},
			{
				SourceType: "http",
				HTTPSource: config2.HTTPSourceConfig{Url: file2Url},
			},
		},
	}
	NewPrometheusSyncer(testutils.LoggerForTesting, config).Sync()

	require.True(t, testutils.FileExists(temporaryDirectory, testutils.GetHashedFilename(file1Url, "yml")))
	require.True(t, testutils.FileExists(temporaryDirectory, testutils.GetHashedFilename(file2Url, "yml")))

	config = config2.PrometheusConfig{
		PrometheusRulesPath: temporaryDirectory,
		PrometheusRules: []config2.PrometheusRuleConfig{
			{
				SourceType: "http",
				HTTPSource: config2.HTTPSourceConfig{Url: file1Url},
			},
			{
				SourceType: "http",
				HTTPSource: config2.HTTPSourceConfig{Url: file2Url},
			},
			{
				SourceType: "http",
				HTTPSource: config2.HTTPSourceConfig{Url: file3Url},
			},
			{
				SourceType: "http",
				HTTPSource: config2.HTTPSourceConfig{Url: file4Url},
			},
		},
	}
	NewPrometheusSyncer(testutils.LoggerForTesting, config).Sync()

	require.True(t, testutils.FileExists(temporaryDirectory, testutils.GetHashedFilename(file1Url, "yml")))
	require.True(t, testutils.FileExists(temporaryDirectory, testutils.GetHashedFilename(file2Url, "yml")))
	require.True(t, testutils.FileExists(temporaryDirectory, testutils.GetHashedFilename(file3Url, "yml")))
	require.True(t, testutils.FileExists(temporaryDirectory, testutils.GetHashedFilename(file4Url, "yml")))
}

func Test_RemoveOldFile(t *testing.T) {
	temporaryDirectory, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	defer os.RemoveAll(temporaryDirectory)

	server1 := testutils.SetupFakeServer(t, "/1.yml", "1.yml content")
	file1Url := fmt.Sprintf("%s/1.yml", server1.URL)
	defer server1.Close()

	server2 := testutils.SetupFakeServer(t, "/2.yml", "2.yml content")
	file2Url := fmt.Sprintf("%s/2.yml", server2.URL)
	defer server2.Close()

	server3 := testutils.SetupFakeServer(t, "/3.yml", "3.yml content")
	file3Url := fmt.Sprintf("%s/3.yml", server3.URL)
	defer server3.Close()

	server4 := testutils.SetupFakeServer(t, "/4.yml", "4.yml content")
	file4Url := fmt.Sprintf("%s/4.yml", server4.URL)
	defer server4.Close()

	config := config2.PrometheusConfig{
		PrometheusRulesPath: temporaryDirectory,
		PrometheusRules: []config2.PrometheusRuleConfig{
			{
				SourceType: "http",
				HTTPSource: config2.HTTPSourceConfig{Url: file1Url},
			},
			{
				SourceType: "http",
				HTTPSource: config2.HTTPSourceConfig{Url: file2Url},
			},
			{
				SourceType: "http",
				HTTPSource: config2.HTTPSourceConfig{Url: file3Url},
			},
			{
				SourceType: "http",
				HTTPSource: config2.HTTPSourceConfig{Url: file4Url},
			},
		},
	}
	NewPrometheusSyncer(testutils.LoggerForTesting, config).Sync()

	require.True(t, testutils.FileExists(temporaryDirectory, testutils.GetHashedFilename(file1Url, "yml")))
	require.True(t, testutils.FileExists(temporaryDirectory, testutils.GetHashedFilename(file2Url, "yml")))
	require.True(t, testutils.FileExists(temporaryDirectory, testutils.GetHashedFilename(file3Url, "yml")))
	require.True(t, testutils.FileExists(temporaryDirectory, testutils.GetHashedFilename(file4Url, "yml")))

	config = config2.PrometheusConfig{
		PrometheusRulesPath: temporaryDirectory,
		PrometheusRules: []config2.PrometheusRuleConfig{
			{
				SourceType: "http",
				HTTPSource: config2.HTTPSourceConfig{Url: file1Url},
			},
			{
				SourceType: "http",
				HTTPSource: config2.HTTPSourceConfig{Url: file2Url},
			},
		},
	}
	NewPrometheusSyncer(testutils.LoggerForTesting, config).Sync()

	require.True(t, testutils.FileExists(temporaryDirectory, testutils.GetHashedFilename(file1Url, "yml")))
	require.True(t, testutils.FileExists(temporaryDirectory, testutils.GetHashedFilename(file2Url, "yml")))
	require.False(t, testutils.FileExists(temporaryDirectory, testutils.GetHashedFilename(file3Url, "yml")))
	require.False(t, testutils.FileExists(temporaryDirectory, testutils.GetHashedFilename(file4Url, "yml")))
}

func Test_UpdateExistingFile(t *testing.T) {
	temporaryDirectory, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	defer os.RemoveAll(temporaryDirectory)

	server1 := testutils.SetupFakeServer(t, "/1.yml", "1.yml content")
	file1Url := fmt.Sprintf("%s/1.yml", server1.URL)
	hashedFilename := testutils.GetHashedFilename(file1Url, "yml")
	defer server1.Close()

	config := config2.PrometheusConfig{
		PrometheusRulesPath: temporaryDirectory,
		PrometheusRules: []config2.PrometheusRuleConfig{
			{
				SourceType: "http",
				HTTPSource: config2.HTTPSourceConfig{Url: file1Url},
			},
		},
	}
	syncer := NewPrometheusSyncer(testutils.LoggerForTesting, config)

	syncer.Sync()

	modificationTime := testutils.GetFileModificationTime(t, temporaryDirectory, hashedFilename)

	server1.Response = "updated content"
	syncer.Sync()

	newModificationTime := testutils.GetFileModificationTime(t, temporaryDirectory, hashedFilename)
	require.NotEqual(t, modificationTime, newModificationTime)
}

func Test_DontUpdateIfContentIsTheSame(t *testing.T) {
	temporaryDirectory, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	defer os.RemoveAll(temporaryDirectory)

	server1 := testutils.SetupFakeServer(t, "/1.yml", "1.yml content")
	file1Url := fmt.Sprintf("%s/1.yml", server1.URL)
	hashedFilename := testutils.GetHashedFilename(file1Url, "yml")
	defer server1.Close()

	config := config2.PrometheusConfig{
		PrometheusRulesPath: temporaryDirectory,
		PrometheusRules: []config2.PrometheusRuleConfig{
			{
				SourceType: "http",
				HTTPSource: config2.HTTPSourceConfig{Url: file1Url},
			},
		},
	}
	syncer := NewPrometheusSyncer(testutils.LoggerForTesting, config)

	syncer.Sync()

	modificationTime := testutils.GetFileModificationTime(t, temporaryDirectory, hashedFilename)

	syncer.Sync()

	newModificationTime := testutils.GetFileModificationTime(t, temporaryDirectory, hashedFilename)
	require.Equal(t, modificationTime, newModificationTime)
}
