package http

import (
	"net/http"
	"time"

	"homework-day1/internal/domain"
)

// HealthHandler chứa handler cho health check endpoint (Bài 5)
type HealthHandler struct {
	usecase   domain.AssetUsecase
	startTime time.Time
}

// NewHealthHandler tạo instance mới của HealthHandler
func NewHealthHandler(uc domain.AssetUsecase, startTime time.Time) *HealthHandler {
	return &HealthHandler{
		usecase:   uc,
		startTime: startTime,
	}
}

// HealthCheck xử lý GET /health
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(h.startTime).Seconds()

	response := domain.HealthResponse{
		Status: "ok",
		Storage: domain.StorageHealth{
			Type:       "in-memory",
			AssetCount: h.usecase.GetAssetCount(),
		},
		Uptime:    uptime,
		Timestamp: time.Now().UTC(),
	}

	respondJSON(w, http.StatusOK, response)
}
