package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_LoadYamlConfig_FullConfig(t *testing.T) {
	content := `
grafana:
  dashboards:
    - source_type: http
      http_source:
        url: http://link.for/dashboard.json
    - source_type: http
      http_source:
        url: http://link.for.another/dashboard.json
`

	expectedConfig := Config{
		Grafana: GrafanaConfig{Dashboards: []GrafanaDashboardsConfig{
			{
				SourceType: "http",
				HTTPSource: HTTPSourceConfig{
					Url: "http://link.for/dashboard.json",
				},
			},
			{
				SourceType: "http",
				HTTPSource: HTTPSourceConfig{
					Url: "http://link.for.another/dashboard.json",
				},
			},
		}},
	}

	actualConfig, err := LoadYamlConfig([]byte(content))
	require.NoError(t, err)
	require.Equal(t, expectedConfig, actualConfig)
}
