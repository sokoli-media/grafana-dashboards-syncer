# How to setup grafana-dashboard-syncer

This file specified how to install and setup this project. It's still a work in progress.

## Requirements

 * `/etc/grafana` is mounted to a persistent storage

## Switch path for provisioning Grafana stuff from file

In `/etc/grafana/grafana.ini` find line:

```
provisioning = conf/provisioning
```

Change it to:

```
provisioning = /etc/grafana/provisioning
```

## Create a datasource in Grafana configuration

Since we're going to use a predefined dashboard, we need to make sure our data source has proper `uid` matching
whatever is specified in the dashboard.

Create a file `/etc/grafana/provisioning/datasources/prometheus.yaml` with content:

```
apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    uid: prometheus
    access: proxy
    url: http://X.X.X.X:9090
    isDefault: false
    editable: false
    version: 1
    jsonData:
      httpMethod: GET
```

Change `http://X.X.X.X:9090` to the url of your Prometheus instance.

Restart Grafana to load the datasource, check if it works.

## Install this project on your UnRaid instance

TODO: add details how to do it

Mount `/dashboards` in your container to `/mnt/user/appdata/grafana/dashboards` on your UnRaid server.

Add post arguments: `--dashboard "unique-dashboard-name.json=http://where-to-get.your/dashboard.json"`

## Mount dashboards to Grafana

Mount `/mnt/user/appdata/grafana/dashboards` on your UnRaid server as `/grafana-dashboards-syncer/dashboards/` in your
Grafana container (read only mode).

## Load dashboards automatically

Create a file `/etc/grafana/provisioning/dashboards/grafana-dashboards-syncer.yaml` with content:

```
apiVersion: 1

providers:
  - name: grafana-dashboards-syncer
    type: file
    disableDeletion: false
    allowUiUpdates: false
    updateIntervalSeconds: 30
    options:
      path: /grafana-dashboards-syncer/dashboards/
```

## Restart Grafana

## It works!

Hopefully.
