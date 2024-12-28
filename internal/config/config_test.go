package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_LoadYamlConfig_Grafana(t *testing.T) {
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

func Test_LoadYamlConfig_Prometheus(t *testing.T) {
	content := `
prometheus:
  prometheus_rules:
    - source_type: http
      http_source:
        url: http://link.for/file.yaml
    - source_type: http
      http_source:
        url: http://link.for.another/file.yaml
`

	expectedConfig := Config{
		Prometheus: PrometheusConfig{PrometheusRules: []PrometheusRuleConfig{
			{
				SourceType: "http",
				HTTPSource: HTTPSourceConfig{
					Url: "http://link.for/file.yaml",
				},
			},
			{
				SourceType: "http",
				HTTPSource: HTTPSourceConfig{
					Url: "http://link.for.another/file.yaml",
				},
			},
		}},
	}

	actualConfig, err := LoadYamlConfig([]byte(content))
	require.NoError(t, err)
	require.Equal(t, expectedConfig, actualConfig)
}
