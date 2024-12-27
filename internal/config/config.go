package config

import (
	"gopkg.in/yaml.v3"
)

type HTTPSourceConfig struct {
	Url string `yaml:"url"`
}

type GrafanaDashboardsConfig struct {
	HTTPSource HTTPSourceConfig `yaml:"http_source"`
}

type GrafanaConfig struct {
	Dashboards []GrafanaDashboardsConfig `yaml:"dashboards"`
}

type Config struct {
	Grafana GrafanaConfig `yaml:"grafana"`
}

func LoadYamlConfig(configFileContent []byte) (Config, error) {
	var config Config
	err := yaml.Unmarshal(configFileContent, &config)
	return config, err
}
