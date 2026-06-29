package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"homework-day1/internal/model"
)

type MockAssetService struct {
	CreateAssetFn       func(ctx context.Context, req *model.CreateAssetRequest) (*model.Asset, error)
	GetAssetByIDFn      func(ctx context.Context, id string) (*model.Asset, error)
	BatchCreateAssetsFn func(ctx context.Context, req *model.BatchCreateRequest) (*model.BatchCreateResponse, error)
	BatchDeleteAssetsFn func(ctx context.Context, ids []string) (*model.BatchDeleteResponse, error)
	ToggleAutoScanFn    func(ctx context.Context, id string, autoScan bool) error
	GetStatsFn          func(ctx context.Context) (*model.StatsResponse, error)
	CountAssetsFn       func(ctx context.Context, assetType, status string) (*model.CountResponse, error)
	ListAssetsFn        func(ctx context.Context, page, limit int, assetType, status string) (*model.ListResponse, error)
	SearchAssetsFn      func(ctx context.Context, query string) ([]*model.Asset, error)
	GetAssetCountFn     func() int
}

func (m *MockAssetService) CreateAsset(ctx context.Context, req *model.CreateAssetRequest) (*model.Asset, error) {
	if m.CreateAssetFn != nil {
		return m.CreateAssetFn(ctx, req)
	}
	return nil, nil
}

func (m *MockAssetService) GetAssetByID(ctx context.Context, id string) (*model.Asset, error) {
	if m.GetAssetByIDFn != nil {
		return m.GetAssetByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockAssetService) BatchCreateAssets(ctx context.Context, req *model.BatchCreateRequest) (*model.BatchCreateResponse, error) {
	if m.BatchCreateAssetsFn != nil {
		return m.BatchCreateAssetsFn(ctx, req)
	}
	return nil, nil
}

func (m *MockAssetService) BatchDeleteAssets(ctx context.Context, ids []string) (*model.BatchDeleteResponse, error) {
	if m.BatchDeleteAssetsFn != nil {
		return m.BatchDeleteAssetsFn(ctx, ids)
	}
	return nil, nil
}

func (m *MockAssetService) ToggleAutoScan(ctx context.Context, id string, autoScan bool) error {
	if m.ToggleAutoScanFn != nil {
		return m.ToggleAutoScanFn(ctx, id, autoScan)
	}
	return nil
}

func (m *MockAssetService) GetStats(ctx context.Context) (*model.StatsResponse, error) {
	if m.GetStatsFn != nil {
		return m.GetStatsFn(ctx)
	}
	return nil, nil
}

func (m *MockAssetService) CountAssets(ctx context.Context, assetType, status string) (*model.CountResponse, error) {
	if m.CountAssetsFn != nil {
		return m.CountAssetsFn(ctx, assetType, status)
	}
	return nil, nil
}

func (m *MockAssetService) ListAssets(ctx context.Context, page, limit int, assetType, status string) (*model.ListResponse, error) {
	if m.ListAssetsFn != nil {
		return m.ListAssetsFn(ctx, page, limit, assetType, status)
	}
	return nil, nil
}

func (m *MockAssetService) SearchAssets(ctx context.Context, query string) ([]*model.Asset, error) {
	if m.SearchAssetsFn != nil {
		return m.SearchAssetsFn(ctx, query)
	}
	return nil, nil
}

func (m *MockAssetService) GetAssetCount() int {
	if m.GetAssetCountFn != nil {
		return m.GetAssetCountFn()
	}
	return 0
}

func TestCreateAssetHandler(t *testing.T) {
	mockService := &MockAssetService{
		CreateAssetFn: func(ctx context.Context, req *model.CreateAssetRequest) (*model.Asset, error) {
			if req.Name == "" {
				return nil, model.ErrInvalidName
			}
			return &model.Asset{
				ID:        "test-id",
				Name:      req.Name,
				Type:      req.Type,
				Status:    req.Status,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}
	handler := NewAssetHandler(mockService)

	t.Run("valid domain asset", func(t *testing.T) {
		body := `{"name":"test.com","type":"domain","status":"active"}`
		req := httptest.NewRequest("POST", "/assets", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.CreateAsset(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", rr.Code)
		}

		var res model.Asset
		if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		if res.Name != "test.com" {
			t.Errorf("expected asset name test.com, got %s", res.Name)
		}
	})

	t.Run("invalid asset name", func(t *testing.T) {
		body := `{"name":"","type":"domain","status":"active"}`
		req := httptest.NewRequest("POST", "/assets", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.CreateAsset(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rr.Code)
		}
	})
}

func TestGetAssetByIDHandler(t *testing.T) {
	mockService := &MockAssetService{
		GetAssetByIDFn: func(ctx context.Context, id string) (*model.Asset, error) {
			if id == "not-found" {
				return nil, model.ErrAssetNotFound
			}
			return &model.Asset{
				ID:   id,
				Name: "example.com",
				Type: "domain",
			}, nil
		},
	}
	handler := NewAssetHandler(mockService)

	t.Run("found asset", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/assets/test-id", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "test-id")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		rr := httptest.NewRecorder()

		handler.GetAssetByID(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rr.Code)
		}
	})

	t.Run("asset not found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/assets/not-found", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "not-found")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		rr := httptest.NewRecorder()

		handler.GetAssetByID(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rr.Code)
		}
	})
}
