package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"homework-day1/internal/model"
)

// AssetHandler chứa các HTTP handler cho asset endpoints
type AssetHandler struct {
	usecase model.AssetService
}

// NewAssetHandler tạo một instance mới của AssetHandler
func NewAssetHandler(uc model.AssetService) *AssetHandler {
	return &AssetHandler{usecase: uc}
}

// respondJSON gửi JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError gửi error response dạng JSON
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// CreateAsset xử lý POST /assets - tạo 1 asset
func (h *AssetHandler) CreateAsset(w http.ResponseWriter, r *http.Request) {
	var req model.CreateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	asset, err := h.usecase.CreateAsset(r.Context(), &req)
	if err != nil {
		if errors.Is(err, model.ErrInvalidAssetType) ||
			errors.Is(err, model.ErrInvalidName) ||
			errors.Is(err, model.ErrInvalidAssetStatus) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusCreated, asset)
}

// GetAssetByID xử lý GET /assets/{id} - lấy 1 asset theo ID
func (h *AssetHandler) GetAssetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	asset, err := h.usecase.GetAssetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrAssetNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, asset)
}

// BatchCreate xử lý POST /assets/batch - tạo nhiều assets (Bài 2)
func (h *AssetHandler) BatchCreate(w http.ResponseWriter, r *http.Request) {
	var req model.BatchCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	result, err := h.usecase.BatchCreateAssets(r.Context(), &req)
	if err != nil {
		if errors.Is(err, model.ErrInvalidAssetType) ||
			errors.Is(err, model.ErrInvalidName) ||
			errors.Is(err, model.ErrInvalidAssetStatus) ||
			errors.Is(err, model.ErrBatchLimitExceeded) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusCreated, result)
}

// BatchDelete xử lý DELETE /assets/batch - xóa nhiều assets (Bài 3)
func (h *AssetHandler) BatchDelete(w http.ResponseWriter, r *http.Request) {
	idsParam := r.URL.Query().Get("ids")
	if idsParam == "" {
		respondError(w, http.StatusBadRequest, model.ErrIDsRequired.Error())
		return
	}

	ids := strings.Split(idsParam, ",")
	// Loại bỏ các ID rỗng
	var cleanIDs []string
	for _, id := range ids {
		trimmed := strings.TrimSpace(id)
		if trimmed != "" {
			cleanIDs = append(cleanIDs, trimmed)
		}
	}

	if len(cleanIDs) == 0 {
		respondError(w, http.StatusBadRequest, model.ErrIDsRequired.Error())
		return
	}

	result, err := h.usecase.BatchDeleteAssets(r.Context(), cleanIDs)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// GetStats xử lý GET /assets/stats - thống kê (Bài 1.1)
func (h *AssetHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.usecase.GetStats(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, stats)
}

// CountAssets xử lý GET /assets/count - đếm theo bộ lọc (Bài 1.2)
func (h *AssetHandler) CountAssets(w http.ResponseWriter, r *http.Request) {
	assetType := r.URL.Query().Get("type")
	status := r.URL.Query().Get("status")

	result, err := h.usecase.CountAssets(r.Context(), assetType, status)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// ListAssets xử lý GET /assets - danh sách có phân trang (Bài 6 - Bonus)
func (h *AssetHandler) ListAssets(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	assetType := r.URL.Query().Get("type")
	status := r.URL.Query().Get("status")

	result, err := h.usecase.ListAssets(r.Context(), page, limit, assetType, status)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// SearchAssets xử lý GET /assets/search - tìm kiếm theo tên (Bài 7 - Bonus)
func (h *AssetHandler) SearchAssets(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		respondError(w, http.StatusBadRequest, model.ErrSearchQueryRequired.Error())
		return
	}

	results, err := h.usecase.SearchAssets(r.Context(), query)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, results)
}

func (h *AssetHandler) ToggleAutoScan(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "missing asset id")
		return
	}

	var req struct {
		AutoScan bool `json:"auto_scan"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.usecase.ToggleAutoScan(r.Context(), id, req.AutoScan); err != nil {
		if err == model.ErrAssetNotFound {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
