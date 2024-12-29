package prometheus_syncer

import (
	"log/slog"
	"time"
)

func RunBackgroundSyncingDaemon(logger *slog.Logger, syncer *PrometheusSyncer) {
	logger.Info("starting background syncing for prometheus")
	for {
		go syncer.Sync()

		time.Sleep(30 * time.Second)
	}
}
