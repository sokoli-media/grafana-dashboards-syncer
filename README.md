# unraid-monitoring-operator

Simple project to help deploying monitoring _stuff_ in UnRaid.

In the real world we usually have big-and-heavy k8s _stuff_ automatically deploying monitoring _stuff_.
The goal of this project is to make it possible also in the UnRaid world (or any other world that is using
docker containers instead of the whole k8s stack.)

### How to set it up?

Coming soon.

### Example config

```yaml
grafana:
  dashboards:
    - source_type: http
      http_source:
        url: http://192.168.1.10:1234/dashboard.json

prometheus:
  reload_config_url: http://192.168.1.1:9000/-/reload  # requires "--web.enable-lifecycle" to be added to the prometheus command
  prometheus_rules_path: /prometheus_rules/
  prometheus_rules:
    - source_type: http
      http_source:
        url: http://192.168.1.10:1234/prometheusRule.yml
```
