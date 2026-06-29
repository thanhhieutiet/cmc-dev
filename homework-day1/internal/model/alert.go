package model

import (
	"context"
	"time"
)

// Alert đại diện cho một cảnh báo an ninh được sinh ra từ việc so sánh kết quả scan
type Alert struct {
	ID        string    `json:"id"`
	AssetID   string    `json:"asset_id"`
	Type      string    `json:"type"` // "port", "ssl", "subdomain"
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

// AlertRepository định nghĩa giao tiếp với storage layer cho bảng alerts
type AlertRepository interface {
	Create(ctx context.Context, alert *Alert) error
	ListUnread(ctx context.Context) ([]*Alert, error)
	ListAll(ctx context.Context) ([]*Alert, error)
	ListByAsset(ctx context.Context, assetID string) ([]*Alert, error)
	MarkAsRead(ctx context.Context, id string) error
	MarkAllAsRead(ctx context.Context) error
}
