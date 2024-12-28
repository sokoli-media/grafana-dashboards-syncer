package config

import (
	"gopkg.in/yaml.v3"
)

type HTTPSourceConfig struct {
	Url string `yaml:"url"`
}

type GrafanaDashboardsConfig struct {
	SourceType string           `yaml:"source_type"`
	HTTPSource HTTPSourceConfig `yaml:"http_source"` // source_type: http
}

type GrafanaConfig struct {
	Dashboards []GrafanaDashboardsConfig `yaml:"dashboards"`
}

type PrometheusRuleConfig struct {
	SourceType string           `yaml:"source_type"`
	HTTPSource HTTPSourceConfig `yaml:"http_source"` // source_type: http
}

type PrometheusConfig struct {
	PrometheusRules []PrometheusRuleConfig `yaml:"prometheus_rules"`
}

type Config struct {
	Grafana    GrafanaConfig    `yaml:"grafana"`
	Prometheus PrometheusConfig `yaml:"prometheus"`
}

func LoadYamlConfig(configFileContent []byte) (Config, error) {
	var config Config
	err := yaml.Unmarshal(configFileContent, &config)
	return config, err
}
