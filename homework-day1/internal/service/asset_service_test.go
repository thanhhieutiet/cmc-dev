package service

import (
	"context"
	"errors"
	"testing"

	"homework-day1/internal/model"
)

type MockAssetRepository struct {
	CreateFn      func(ctx context.Context, asset *model.Asset) error
	GetByIDFn     func(ctx context.Context, id string) (*model.Asset, error)
	DeleteFn      func(ctx context.Context, id string) error
	BatchCreateFn func(ctx context.Context, assets []*model.Asset) error
	BatchDeleteFn func(ctx context.Context, ids []string) (int, int, error)
	UpdateAutoScanFn   func(ctx context.Context, id string, autoScan bool) error
	GetAutoScanAssetsFn func(ctx context.Context) ([]*model.Asset, error)
	GetStatsFn    func(ctx context.Context) (*model.StatsResponse, error)
	CountFn       func(ctx context.Context, assetType, status string) (int, error)
	ListFn        func(ctx context.Context, page, limit int, assetType, status string) ([]*model.Asset, int, error)
	SearchFn      func(ctx context.Context, query string, maxResults int) ([]*model.Asset, error)
	TotalCountFn  func() int
}

func (m *MockAssetRepository) Create(ctx context.Context, asset *model.Asset) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, asset)
	}
	return nil
}

func (m *MockAssetRepository) GetByID(ctx context.Context, id string) (*model.Asset, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockAssetRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func (m *MockAssetRepository) BatchCreate(ctx context.Context, assets []*model.Asset) error {
	if m.BatchCreateFn != nil {
		return m.BatchCreateFn(ctx, assets)
	}
	return nil
}

func (m *MockAssetRepository) BatchDelete(ctx context.Context, ids []string) (int, int, error) {
	if m.BatchDeleteFn != nil {
		return m.BatchDeleteFn(ctx, ids)
	}
	return 0, 0, nil
}

func (m *MockAssetRepository) UpdateAutoScan(ctx context.Context, id string, autoScan bool) error {
	if m.UpdateAutoScanFn != nil {
		return m.UpdateAutoScanFn(ctx, id, autoScan)
	}
	return nil
}

func (m *MockAssetRepository) GetAutoScanAssets(ctx context.Context) ([]*model.Asset, error) {
	if m.GetAutoScanAssetsFn != nil {
		return m.GetAutoScanAssetsFn(ctx)
	}
	return nil, nil
}

func (m *MockAssetRepository) GetStats(ctx context.Context) (*model.StatsResponse, error) {
	if m.GetStatsFn != nil {
		return m.GetStatsFn(ctx)
	}
	return nil, nil
}

func (m *MockAssetRepository) Count(ctx context.Context, assetType, status string) (int, error) {
	if m.CountFn != nil {
		return m.CountFn(ctx, assetType, status)
	}
	return 0, nil
}

func (m *MockAssetRepository) List(ctx context.Context, page, limit int, assetType, status string) ([]*model.Asset, int, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, page, limit, assetType, status)
	}
	return nil, 0, nil
}

func (m *MockAssetRepository) Search(ctx context.Context, query string, maxResults int) ([]*model.Asset, error) {
	if m.SearchFn != nil {
		return m.SearchFn(ctx, query, maxResults)
	}
	return nil, nil
}

func (m *MockAssetRepository) TotalCount() int {
	if m.TotalCountFn != nil {
		return m.TotalCountFn()
	}
	return 0
}

func TestAssetService_Create(t *testing.T) {
	mockRepo := &MockAssetRepository{
		CreateFn: func(ctx context.Context, asset *model.Asset) error {
			if asset.Name == "fail.com" {
				return errors.New("db error")
			}
			return nil
		},
	}
	service := NewAssetService(mockRepo)

	t.Run("successful create", func(t *testing.T) {
		req := &model.CreateAssetRequest{
			Name:   "example.com",
			Type:   "domain",
			Status: "active",
		}
		res, err := service.CreateAsset(context.Background(), req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if res.Name != "example.com" {
			t.Errorf("expected asset name example.com, got %s", res.Name)
		}
		if res.ID == "" {
			t.Errorf("expected asset ID to be populated")
		}
	})

	t.Run("invalid model validation", func(t *testing.T) {
		req := &model.CreateAssetRequest{
			Name:   "", // empty name is invalid
			Type:   "domain",
			Status: "active",
		}
		_, err := service.CreateAsset(context.Background(), req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("repository error", func(t *testing.T) {
		req := &model.CreateAssetRequest{
			Name:   "fail.com",
			Type:   "domain",
			Status: "active",
		}
		_, err := service.CreateAsset(context.Background(), req)
		if err == nil {
			t.Fatal("expected error from repository, got nil")
		}
	})
}
