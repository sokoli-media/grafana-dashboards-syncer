package prometheus_syncer

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log/slog"
	"os"
	"path/filepath"
	"unraid-monitoring-operator/internal/config"
	"unraid-monitoring-operator/internal/http_downloader"
	"unraid-monitoring-operator/internal/trash_collector"
)

var prometheusRuleSynced = promauto.NewCounterVec(
	prometheus.CounterOpts{Name: "unraid_monitoring_operator_prometheus_rule_synced"},
	[]string{"url"})
var prometheusRuleUpdated = promauto.NewCounterVec(
	prometheus.CounterOpts{Name: "unraid_monitoring_operator_prometheus_rule_updated"},
	[]string{"url"})
var prometheusRuleSyncErrored = promauto.NewCounterVec(
	prometheus.CounterOpts{Name: "unraid_monitoring_operator_prometheus_rule_sync_errored"},
	[]string{"url", "reason"})

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
	trashCollector := trash_collector.NewTrashCollector(p.config.PrometheusRulesPath)

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

	err := trashCollector.PickUpTrash()
	if err != nil {
		p.logger.Error("couldn't delete unknown files", "error", err)
	}
}

func (p *PrometheusSyncer) generateFilename(prometheusRule config.PrometheusRuleConfig) string {
	md5sum := md5.New()
	md5sum.Write([]byte(prometheusRule.HTTPSource.Url))
	filenameBase := hex.EncodeToString(md5sum.Sum(nil))
	return fmt.Sprintf("%s.yml", filenameBase)
}
