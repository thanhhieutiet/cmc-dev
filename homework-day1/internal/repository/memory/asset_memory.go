package memory

import (
	"context"
	"strings"
	"sync"

	"homework-day1/internal/model"
)

// InMemoryAssetRepo hiện thực AssetRepository sử dụng bộ nhớ trong (map)
// Sử dụng sync.RWMutex để đảm bảo concurrent-safe (Bài 4)
type InMemoryAssetRepo struct {
	mu     sync.RWMutex
	assets map[string]*model.Asset
}

// NewInMemoryAssetRepo tạo một instance mới của InMemoryAssetRepo
func NewInMemoryAssetRepo() *InMemoryAssetRepo {
	return &InMemoryAssetRepo{
		assets: make(map[string]*model.Asset),
	}
}

// Create thêm một asset vào bộ nhớ
func (r *InMemoryAssetRepo) Create(ctx context.Context, asset *model.Asset) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.assets[asset.ID] = asset
	return nil
}

// GetByID tìm asset theo ID
func (r *InMemoryAssetRepo) GetByID(ctx context.Context, id string) (*model.Asset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	asset, exists := r.assets[id]
	if !exists {
		return nil, model.ErrAssetNotFound
	}
	return asset, nil
}

// Delete xóa một asset theo ID
func (r *InMemoryAssetRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.assets[id]; !exists {
		return model.ErrAssetNotFound
	}
	delete(r.assets, id)
	return nil
}

// BatchCreate thêm nhiều assets cùng lúc (đã được validate trước bởi usecase)
func (r *InMemoryAssetRepo) BatchCreate(ctx context.Context, assets []*model.Asset) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, asset := range assets {
		r.assets[asset.ID] = asset
	}
	return nil
}

// BatchDelete xóa nhiều assets theo danh sách IDs
// Trả về số lượng đã xóa thành công và số lượng không tìm thấy
func (r *InMemoryAssetRepo) BatchDelete(ctx context.Context, ids []string) (int, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	deleted := 0
	notFound := 0

	for _, id := range ids {
		if _, exists := r.assets[id]; exists {
			delete(r.assets, id)
			deleted++
		} else {
			notFound++
		}
	}

	return deleted, notFound, nil
}

// GetStats trả về thống kê tổng quan về các assets (Bài 1.1)
func (r *InMemoryAssetRepo) GetStats(ctx context.Context) (*model.StatsResponse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stats := &model.StatsResponse{
		Total:    len(r.assets),
		ByType:   make(map[string]int),
		ByStatus: make(map[string]int),
	}

	for _, asset := range r.assets {
		stats.ByType[asset.Type]++
		stats.ByStatus[asset.Status]++
	}

	return stats, nil
}

// Count đếm số assets theo bộ lọc type và status (Bài 1.2)
func (r *InMemoryAssetRepo) Count(ctx context.Context, assetType, status string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, asset := range r.assets {
		if matchFilter(asset, assetType, status) {
			count++
		}
	}
	return count, nil
}

// List trả về danh sách assets có phân trang và lọc (Bài 6 - Bonus)
func (r *InMemoryAssetRepo) List(ctx context.Context, page, limit int, assetType, status string) ([]*model.Asset, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Lọc trước
	var filtered []*model.Asset
	for _, asset := range r.assets {
		if matchFilter(asset, assetType, status) {
			filtered = append(filtered, asset)
		}
	}

	total := len(filtered)

	// Phân trang
	start := (page - 1) * limit
	if start >= total {
		return []*model.Asset{}, total, nil
	}

	end := start + limit
	if end > total {
		end = total
	}

	return filtered[start:end], total, nil
}

// Search tìm kiếm assets theo tên (partial match, case-insensitive) (Bài 7 - Bonus)
func (r *InMemoryAssetRepo) Search(ctx context.Context, query string, maxResults int) ([]*model.Asset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	lowerQuery := strings.ToLower(query)
	var results []*model.Asset

	for _, asset := range r.assets {
		if strings.Contains(strings.ToLower(asset.Name), lowerQuery) {
			results = append(results, asset)
			if len(results) >= maxResults {
				break
			}
		}
	}

	return results, nil
}

// TotalCount trả về tổng số assets hiện có (dùng cho Health Check)
func (r *InMemoryAssetRepo) TotalCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.assets)
}

// GetStorageType trả về loại storage là in-memory (dùng cho Health Check)
func (r *InMemoryAssetRepo) GetStorageType() string {
	return "in-memory"
}

// matchFilter kiểm tra asset có khớp với bộ lọc type và status không
func matchFilter(a *model.Asset, assetType, status string) bool {
	if assetType != "" && a.Type != assetType {
		return false
	}
	if status != "" && a.Status != status {
		return false
	}
	return true
}
