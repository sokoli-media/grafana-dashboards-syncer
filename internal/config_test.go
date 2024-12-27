package internal

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_LoadYamlConfig_FullConfig(t *testing.T) {
	content := `
grafana:
  dashboards:
    - http_source:
        url: http://link.for/dashboard.json
    - http_source:
        url: http://link.for.another/dashboard.json
`

	expectedConfig := Config{
		Grafana: GrafanaConfig{Dashboards: []GrafanaDashboardsConfig{
			{
				HTTPSource: GrafanaHTTPSourceConfig{
					Url: "http://link.for/dashboard.json",
				},
			},
			{
				HTTPSource: GrafanaHTTPSourceConfig{
					Url: "http://link.for.another/dashboard.json",
				},
			},
		}},
	}

	actualConfig, err := LoadYamlConfig([]byte(content))
	require.NoError(t, err)
	require.Equal(t, expectedConfig, actualConfig)
}
