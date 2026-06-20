package scheduler

import (
	"context"
	"log"
	"time"

	"homework-day1/internal/model"
)

type Scheduler struct {
	assetRepo   model.AssetRepository
	scanService model.ScanService
	interval    time.Duration
	stopCh      chan struct{}
}

func NewScheduler(assetRepo model.AssetRepository, scanService model.ScanService, interval time.Duration) *Scheduler {
	return &Scheduler{
		assetRepo:   assetRepo,
		scanService: scanService,
		interval:    interval,
		stopCh:      make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	ticker := time.NewTicker(s.interval)
	log.Printf("[Scheduler] Started auto-scan scheduler. Interval: %v", s.interval)

	go func() {
		for {
			select {
			case <-ticker.C:
				s.runScheduledScans()
			case <-s.stopCh:
				ticker.Stop()
				log.Println("[Scheduler] Stopped auto-scan scheduler.")
				return
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	close(s.stopCh)
}

func (s *Scheduler) runScheduledScans() {
	ctx := context.Background()
	// NOTE: We assume GetAutoScanAssets is implemented in AssetRepository
	// We might need to cast or use a specific interface if not in model.AssetRepository
	
	// Since GetAutoScanAssets is added to PostgresAssetRepo, we should ideally define it in AssetRepository interface.
	// We'll call it directly if it's there.
	type autoScanRepo interface {
		GetAutoScanAssets(ctx context.Context) ([]*model.Asset, error)
	}

	repo, ok := s.assetRepo.(autoScanRepo)
	if !ok {
		log.Println("[Scheduler] Warning: AssetRepository does not support GetAutoScanAssets")
		return
	}

	assets, err := repo.GetAutoScanAssets(ctx)
	if err != nil {
		log.Printf("[Scheduler] Error fetching auto-scan assets: %v", err)
		return
	}

	if len(assets) == 0 {
		return
	}

	log.Printf("[Scheduler] Triggering scheduled scans for %d assets...", len(assets))

	for _, asset := range assets {
		_, err := s.scanService.StartScan(ctx, asset.ID, model.ScanTypeAll)
		if err != nil {
			log.Printf("[Scheduler] Failed to trigger scan for asset %s: %v", asset.ID, err)
		} else {
			log.Printf("[Scheduler] Auto-scan triggered for asset %s", asset.Name)
		}
	}
}
