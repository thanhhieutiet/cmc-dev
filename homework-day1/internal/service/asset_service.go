package service

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"

	"homework-day1/internal/model"
)

// assetService hiện thực model.AssetService
type assetService struct {
	repo model.AssetRepository
}

// NewAssetService tạo một instance mới của AssetService
func NewAssetService(repo model.AssetRepository) model.AssetService {
	return &assetService{repo: repo}
}

// CreateAsset tạo một asset mới
func (uc *assetService) CreateAsset(ctx context.Context, req *model.CreateAssetRequest) (*model.Asset, error) {
	now := time.Now()
	asset := &model.Asset{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Type:      req.Type,
		Status:    req.Status,
		Tags:      req.Tags,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Gán status mặc định nếu không truyền
	if asset.Status == "" {
		asset.Status = "active"
	}

	// Validate
	if err := asset.Validate(); err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, asset); err != nil {
		return nil, err
	}

	return asset, nil
}

// GetAssetByID lấy asset theo ID
func (uc *assetService) GetAssetByID(ctx context.Context, id string) (*model.Asset, error) {
	return uc.repo.GetByID(ctx, id)
}

// BatchCreateAssets tạo nhiều assets cùng lúc (Bài 2)
// Nguyên tắc: All or Nothing - validate toàn bộ trước khi insert
func (uc *assetService) BatchCreateAssets(ctx context.Context, req *model.BatchCreateRequest) (*model.BatchCreateResponse, error) {
	// Kiểm tra giới hạn 100 assets/request
	if len(req.Assets) > 100 {
		return nil, model.ErrBatchLimitExceeded
	}

	if len(req.Assets) == 0 {
		return &model.BatchCreateResponse{Created: 0, IDs: []string{}}, nil
	}

	// Bước 1: Validate TOÀN BỘ trước (All or Nothing)
	assets := make([]*model.Asset, 0, len(req.Assets))
	for _, item := range req.Assets {
		now := time.Now()
		asset := &model.Asset{
			ID:        uuid.New().String(),
			Name:      item.Name,
			Type:      item.Type,
			Status:    item.Status,
			Tags:      item.Tags,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if asset.Status == "" {
			asset.Status = "active"
		}

		if err := asset.Validate(); err != nil {
			// Một asset lỗi -> dừng tất cả, không tạo asset nào
			return nil, err
		}

		assets = append(assets, asset)
	}

	// Bước 2: Insert toàn bộ (đã qua validate)
	if err := uc.repo.BatchCreate(ctx, assets); err != nil {
		return nil, err
	}

	ids := make([]string, len(assets))
	for i, a := range assets {
		ids[i] = a.ID
	}

	return &model.BatchCreateResponse{
		Created: len(assets),
		IDs:     ids,
	}, nil
}

// BatchDeleteAssets xóa nhiều assets cùng lúc (Bài 3)
func (uc *assetService) BatchDeleteAssets(ctx context.Context, ids []string) (*model.BatchDeleteResponse, error) {
	deleted, notFound, err := uc.repo.BatchDelete(ctx, ids)
	if err != nil {
		return nil, err
	}

	return &model.BatchDeleteResponse{
		Deleted:  deleted,
		NotFound: notFound,
	}, nil
}

func (s *assetService) ToggleAutoScan(ctx context.Context, id string, autoScan bool) error {
	return s.repo.UpdateAutoScan(ctx, id, autoScan)
}

// GetStats lấy thống kê tổng quan (Bài 1.1)
func (uc *assetService) GetStats(ctx context.Context) (*model.StatsResponse, error) {
	return uc.repo.GetStats(ctx)
}

// CountAssets đếm assets theo bộ lọc (Bài 1.2)
func (uc *assetService) CountAssets(ctx context.Context, assetType, status string) (*model.CountResponse, error) {
	count, err := uc.repo.Count(ctx, assetType, status)
	if err != nil {
		return nil, err
	}

	filters := make(map[string]string)
	if assetType != "" {
		filters["type"] = assetType
	}
	if status != "" {
		filters["status"] = status
	}

	return &model.CountResponse{
		Count:   count,
		Filters: filters,
	}, nil
}

// ListAssets lấy danh sách assets có phân trang và lọc (Bài 6 - Bonus)
func (uc *assetService) ListAssets(ctx context.Context, page, limit int, assetType, status string) (*model.ListResponse, error) {
	// Giá trị mặc định
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	assets, total, err := uc.repo.List(ctx, page, limit, assetType, status)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &model.ListResponse{
		Data: assets,
		Pagination: model.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// SearchAssets tìm kiếm assets theo tên (Bài 7 - Bonus)
func (uc *assetService) SearchAssets(ctx context.Context, query string) ([]*model.Asset, error) {
	return uc.repo.Search(ctx, query, 100)
}

// GetAssetCount trả về tổng số assets (dùng cho Health Check)
func (uc *assetService) GetAssetCount() int {
	return uc.repo.TotalCount()
}

// GetStorageType trả về loại storage đang dùng (dùng cho Health Check)
func (uc *assetService) GetStorageType() string {
	if st, ok := uc.repo.(interface{ GetStorageType() string }); ok {
		return st.GetStorageType()
	}
	return "in-memory"
}
