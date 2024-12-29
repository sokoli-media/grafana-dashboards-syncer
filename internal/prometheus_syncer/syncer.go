package prometheus_syncer

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"unraid-monitoring-operator/internal/config"
	"unraid-monitoring-operator/internal/http_downloader"
	"unraid-monitoring-operator/internal/trash_collector"
)

var prometheusSyncStarted = promauto.NewCounter(
	prometheus.CounterOpts{Name: "unraid_monitoring_operator_prometheus_sync_started"})
var prometheusRuleSynced = promauto.NewCounterVec(
	prometheus.CounterOpts{Name: "unraid_monitoring_operator_prometheus_rule_synced"},
	[]string{"url"})
var prometheusRuleUpdated = promauto.NewCounterVec(
	prometheus.CounterOpts{Name: "unraid_monitoring_operator_prometheus_rule_updated"},
	[]string{"url"})
var prometheusRuleSyncErrored = promauto.NewCounterVec(
	prometheus.CounterOpts{Name: "unraid_monitoring_operator_prometheus_rule_sync_errored"},
	[]string{"url", "reason"})
var prometheusSyncErrored = promauto.NewCounterVec(
	prometheus.CounterOpts{Name: "unraid_monitoring_operator_prometheus_sync_errored"},
	[]string{"reason"})

func NewPrometheusSyncer(logger *slog.Logger, config config.PrometheusConfig) *PrometheusSyncer {
	return &PrometheusSyncer{
		logger:               logger,
		config:               config,
		downloadedFilesCache: make(map[string]string),
	}
}

type PrometheusSyncer struct {
	logger               *slog.Logger
	config               config.PrometheusConfig
	downloadedFilesCache map[string]string
}

func (p *PrometheusSyncer) Sync() {
	prometheusSyncStarted.Inc()
	trashCollector := trash_collector.NewTrashCollector(p.config.PrometheusRulesPath)

	if err := p.configurationIsValid(); err != nil {
		p.logger.Info("ignoring prometheus sync", "reason", err)
		prometheusSyncErrored.With(
			prometheus.Labels{
				"reason": fmt.Sprintf("%s", err),
			},
		)
		return
	}

	updated := false
	for _, prometheusRule := range p.config.PrometheusRules {
		p.logger.Info("syncing prometheusRule", "url", prometheusRule.HTTPSource.Url)
		content, err := http_downloader.Download(prometheusRule.HTTPSource.Url)
		if err != nil {
			p.logger.Error("couldn't download prometheus rules", "error", err, "url", prometheusRule.HTTPSource.Url)
			prometheusRuleSyncErrored.With(
				prometheus.Labels{
					"url":    prometheusRule.HTTPSource.Url,
					"reason": fmt.Sprintf("%s", err),
				},
			)
			continue
		}

		filename := p.generateFilename(prometheusRule)

		cachedValue, exists := p.downloadedFilesCache[filename]
		if !exists || cachedValue != string(content) {
			fullPath := filepath.Join(p.config.PrometheusRulesPath, filename)
			p.logger.Debug("writing file to disk")
			err = os.WriteFile(fullPath, content, 0644)
			if err != nil {
				p.logger.Error(
					"couldn't save prometheus rule file to disk",
					"error", err,
					"url", prometheusRule.HTTPSource.Url,
					"path", fullPath,
				)
				prometheusRuleSyncErrored.With(
					prometheus.Labels{
						"url":    prometheusRule.HTTPSource.Url,
						"reason": fmt.Sprintf("%s", err),
					},
				)
				continue
			}

			updated = true
			p.downloadedFilesCache[filename] = string(content)
			prometheusRuleUpdated.With(
				prometheus.Labels{
					"url": prometheusRule.HTTPSource.Url,
				},
			)
		}

		p.logger.Debug("calling trash collector")
		trashCollector.AddKnownFile(filename)
		prometheusRuleSynced.With(prometheus.Labels{"url": prometheusRule.HTTPSource.Url})
	}

	if updated && p.config.ReloadConfigUrl != "" {
		resp, err := http.Post(p.config.ReloadConfigUrl, "", nil)
		if err != nil || resp.StatusCode != 200 {
			p.logger.Error("couldn't reload prometheus configuration", "error", err)
			prometheusSyncErrored.With(
				prometheus.Labels{
					"reason": fmt.Sprintf("%s", err),
				},
			)
		}
		defer resp.Body.Close()
	}

	err := trashCollector.PickUpTrash()
	if err != nil {
		p.logger.Error("couldn't delete unknown files", "error", err)
		prometheusSyncErrored.With(
			prometheus.Labels{
				"reason": fmt.Sprintf("%s", err),
			},
		)
	}
}

func (p *PrometheusSyncer) generateFilename(prometheusRule config.PrometheusRuleConfig) string {
	md5sum := md5.New()
	md5sum.Write([]byte(prometheusRule.HTTPSource.Url))
	filenameBase := hex.EncodeToString(md5sum.Sum(nil))
	return fmt.Sprintf("%s.yml", filenameBase)
}

func (p *PrometheusSyncer) configurationIsValid() error {
	if fileStat, err := os.Stat(p.config.PrometheusRulesPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("prometheus rule path '%s' doesn't exist", p.config.PrometheusRulesPath)
		} else {
			if !fileStat.IsDir() {
				return fmt.Errorf("prometheus rule path '%s' is not a directory", p.config.PrometheusRulesPath)
			}
		}
	}

	if len(p.config.PrometheusRules) == 0 {
		return errors.New("no prometheus rules to be updated, ignoring sync")
	}

	return nil
}
