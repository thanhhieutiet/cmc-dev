package memory

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"homework-day1/internal/domain"
)

// TestConcurrentCreate kiểm tra race condition khi tạo asset đồng thời (Bài 4)
func TestConcurrentCreate(t *testing.T) {
	repo := NewInMemoryAssetRepo()
	ctx := context.Background()

	const numGoroutines = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(n int) {
			defer wg.Done()
			asset := &domain.Asset{
				ID:        fmt.Sprintf("id-%d", n),
				Name:      fmt.Sprintf("concurrent-%d.com", n),
				Type:      "domain",
				Status:    "active",
				CreatedAt: time.Now(),
			}
			if err := repo.Create(ctx, asset); err != nil {
				t.Errorf("Create failed: %v", err)
			}
		}(i)
	}

	wg.Wait()

	if count := repo.TotalCount(); count != numGoroutines {
		t.Errorf("Expected %d assets, got %d", numGoroutines, count)
	}
}

// TestConcurrentBatchCreateAndStats kiểm tra batch create và stats đồng thời
func TestConcurrentBatchCreateAndStats(t *testing.T) {
	repo := NewInMemoryAssetRepo()
	ctx := context.Background()

	var wg sync.WaitGroup

	// 10 goroutines tạo batch, 10 goroutines đọc stats
	for i := 0; i < 10; i++ {
		wg.Add(2)

		go func(n int) {
			defer wg.Done()
			assets := []*domain.Asset{
				{ID: fmt.Sprintf("batch-%d-1", n), Name: fmt.Sprintf("batch-%d-1.com", n), Type: "domain", Status: "active", CreatedAt: time.Now()},
				{ID: fmt.Sprintf("batch-%d-2", n), Name: fmt.Sprintf("batch-%d-2.com", n), Type: "ip", Status: "inactive", CreatedAt: time.Now()},
			}
			if err := repo.BatchCreate(ctx, assets); err != nil {
				t.Errorf("BatchCreate failed: %v", err)
			}
		}(i)

		go func() {
			defer wg.Done()
			_, _ = repo.GetStats(ctx)
		}()
	}

	wg.Wait()

	if count := repo.TotalCount(); count != 20 {
		t.Errorf("Expected 20 assets, got %d", count)
	}
}
