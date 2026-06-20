package handler

import (
	"net/http"
	"time"

	"homework-day1/internal/model"
)

// HealthHandler chứa handler cho health check endpoint (Bài 5)
type HealthHandler struct {
	usecase   model.AssetService
	startTime time.Time
}

// NewHealthHandler tạo instance mới của HealthHandler
func NewHealthHandler(uc model.AssetService, startTime time.Time) *HealthHandler {
	return &HealthHandler{
		usecase:   uc,
		startTime: startTime,
	}
}

// HealthCheck xử lý GET /health
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(h.startTime).Seconds()

	storageType := "in-memory"
	if st, ok := h.usecase.(interface{ GetStorageType() string }); ok {
		storageType = st.GetStorageType()
	}

	response := model.HealthResponse{
		Status: "ok",
		Storage: model.StorageHealth{
			Type:       storageType,
			AssetCount: h.usecase.GetAssetCount(),
		},
		Uptime:    uptime,
		Timestamp: time.Now().UTC(),
	}

	respondJSON(w, http.StatusOK, response)
}
