package domain

import (
	"context"
	"errors"
	"time"
)

// --- Entities ---

// Asset đại diện cho tài sản số - thực thể cốt lõi của bài toán
type Asset struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`       // domain, ip, service
	Status    string    `json:"status"`     // active, inactive
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// --- Domain Errors ---

var (
	ErrInvalidAssetType   = errors.New("invalid asset type: must be domain, ip, or service")
	ErrInvalidAssetStatus = errors.New("invalid asset status: must be active or inactive")
	ErrInvalidName        = errors.New("asset name cannot be empty")
	ErrAssetNotFound      = errors.New("asset not found")
	ErrBatchLimitExceeded = errors.New("batch limit exceeded: maximum 100 assets per request")
	ErrIDsRequired        = errors.New("ids parameter is required")
	ErrSearchQueryRequired = errors.New("search query parameter 'q' is required")
)

// ValidTypes chứa các loại asset hợp lệ
var ValidTypes = map[string]bool{
	"domain":  true,
	"ip":      true,
	"service": true,
}

// ValidStatuses chứa các trạng thái hợp lệ
var ValidStatuses = map[string]bool{
	"active":   true,
	"inactive": true,
}

// Validate kiểm tra tính hợp lệ của Asset
func (a *Asset) Validate() error {
	if a.Name == "" {
		return ErrInvalidName
	}
	if !ValidTypes[a.Type] {
		return ErrInvalidAssetType
	}
	if a.Status != "" && !ValidStatuses[a.Status] {
		return ErrInvalidAssetStatus
	}
	return nil
}

// --- Request / Response DTOs ---

// CreateAssetRequest là DTO cho việc tạo 1 asset
type CreateAssetRequest struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status,omitempty"`
}

// BatchCreateRequest là DTO cho Bài 2 - tạo nhiều assets
type BatchCreateRequest struct {
	Assets []CreateAssetRequest `json:"assets"`
}

// BatchCreateResponse trả về kết quả tạo hàng loạt
type BatchCreateResponse struct {
	Created int      `json:"created"`
	IDs     []string `json:"ids"`
}

// BatchDeleteResponse trả về kết quả xóa hàng loạt
type BatchDeleteResponse struct {
	Deleted  int `json:"deleted"`
	NotFound int `json:"not_found"`
}

// StatsResponse trả về thống kê tổng quan - Bài 1.1
type StatsResponse struct {
	Total    int            `json:"total"`
	ByType   map[string]int `json:"by_type"`
	ByStatus map[string]int `json:"by_status"`
}

// CountResponse trả về kết quả đếm - Bài 1.2
type CountResponse struct {
	Count   int            `json:"count"`
	Filters map[string]string `json:"filters"`
}

// PaginationMeta chứa metadata phân trang - Bài 6
type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// ListResponse trả về danh sách có phân trang
type ListResponse struct {
	Data       []*Asset       `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// HealthResponse trả về trạng thái sức khỏe - Bài 5
type HealthResponse struct {
	Status    string        `json:"status"`
	Storage   StorageHealth `json:"storage"`
	Uptime    float64       `json:"uptime_seconds"`
	Timestamp time.Time     `json:"timestamp"`
}

// StorageHealth chứa thông tin về storage
type StorageHealth struct {
	Type       string `json:"type"`
	AssetCount int    `json:"asset_count"`
}

// --- Interfaces (Hợp đồng) ---

// AssetRepository định nghĩa interface cho lớp truy cập dữ liệu
type AssetRepository interface {
	// CRUD cơ bản
	Create(ctx context.Context, asset *Asset) error
	GetByID(ctx context.Context, id string) (*Asset, error)
	Delete(ctx context.Context, id string) error

	// Batch operations
	BatchCreate(ctx context.Context, assets []*Asset) error
	BatchDelete(ctx context.Context, ids []string) (deleted int, notFound int, err error)

	// Statistics & Counting
	GetStats(ctx context.Context) (*StatsResponse, error)
	Count(ctx context.Context, assetType, status string) (int, error)

	// Listing & Searching
	List(ctx context.Context, page, limit int, assetType, status string) ([]*Asset, int, error)
	Search(ctx context.Context, query string, maxResults int) ([]*Asset, error)

	// Health
	TotalCount() int
}

// AssetUsecase định nghĩa interface cho lớp logic nghiệp vụ
type AssetUsecase interface {
	// CRUD cơ bản
	CreateAsset(ctx context.Context, req *CreateAssetRequest) (*Asset, error)
	GetAssetByID(ctx context.Context, id string) (*Asset, error)

	// Batch operations
	BatchCreateAssets(ctx context.Context, req *BatchCreateRequest) (*BatchCreateResponse, error)
	BatchDeleteAssets(ctx context.Context, ids []string) (*BatchDeleteResponse, error)

	// Statistics & Counting
	GetStats(ctx context.Context) (*StatsResponse, error)
	CountAssets(ctx context.Context, assetType, status string) (*CountResponse, error)

	// Listing & Searching
	ListAssets(ctx context.Context, page, limit int, assetType, status string) (*ListResponse, error)
	SearchAssets(ctx context.Context, query string) ([]*Asset, error)

	// Health
	GetAssetCount() int
}
