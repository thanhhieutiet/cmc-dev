package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"homework-day1/internal/model"
)

type ScanHandler struct {
	usecase model.ScanService
}

func NewScanHandler(uc model.ScanService) *ScanHandler {
	return &ScanHandler{usecase: uc}
}

// StartScan xử lý POST /assets/{id}/scan
func (h *ScanHandler) StartScan(w http.ResponseWriter, r *http.Request) {
	assetID := chi.URLParam(r, "id")
	if assetID == "" {
		respondError(w, http.StatusBadRequest, "asset id is required")
		return
	}

	var req model.StartScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	job, err := h.usecase.StartScan(r.Context(), assetID, req.ScanType)
	if err != nil {
		if err == model.ErrAssetNotFound || errors.Is(err, model.ErrAssetNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		if err == model.ErrInvalidScanType || err == model.ErrScanNotSupported || err == model.ErrPortScanUnauthorized ||
			errors.Is(err, model.ErrInvalidScanType) || errors.Is(err, model.ErrScanNotSupported) || errors.Is(err, model.ErrPortScanUnauthorized) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to start scan: "+err.Error())
		return
	}

	respondJSON(w, http.StatusAccepted, job)
}

// GetScanJob xử lý GET /scan-jobs/{id}
func (h *ScanHandler) GetScanJob(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "id")
	if jobID == "" {
		respondError(w, http.StatusBadRequest, "job id is required")
		return
	}

	job, err := h.usecase.GetScanJob(r.Context(), jobID)
	if err != nil {
		if err == model.ErrScanJobNotFound || errors.Is(err, model.ErrScanJobNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to get scan job: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, job)
}

// GetScanResults xử lý GET /scan-jobs/{id}/results
func (h *ScanHandler) GetScanResults(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "id")
	if jobID == "" {
		respondError(w, http.StatusBadRequest, "job id is required")
		return
	}

	res, err := h.usecase.GetScanResults(r.Context(), jobID)
	if err != nil {
		if err == model.ErrScanJobNotFound || errors.Is(err, model.ErrScanJobNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to get scan results: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, res)
}

// ListScanJobs xử lý GET /assets/{id}/scans
func (h *ScanHandler) ListScanJobs(w http.ResponseWriter, r *http.Request) {
	assetID := chi.URLParam(r, "id")
	if assetID == "" {
		respondError(w, http.StatusBadRequest, "asset id is required")
		return
	}

	jobs, err := h.usecase.ListScanJobs(r.Context(), assetID)
	if err != nil {
		if err == model.ErrAssetNotFound || errors.Is(err, model.ErrAssetNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to list scan jobs: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, jobs)
}

// GetAssetAllResults xử lý GET /assets/{id}/results
func (h *ScanHandler) GetAssetAllResults(w http.ResponseWriter, r *http.Request) {
	assetID := chi.URLParam(r, "id")
	if assetID == "" {
		respondError(w, http.StatusBadRequest, "asset id is required")
		return
	}

	res, err := h.usecase.GetAssetAllResults(r.Context(), assetID)
	if err != nil {
		if err == model.ErrAssetNotFound || errors.Is(err, model.ErrAssetNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to get asset results: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, res)
}

// GetAssetDNS xử lý GET /assets/{id}/dns
func (h *ScanHandler) GetAssetDNS(w http.ResponseWriter, r *http.Request) {
	assetID := chi.URLParam(r, "id")
	if assetID == "" {
		respondError(w, http.StatusBadRequest, "asset id is required")
		return
	}

	recs, err := h.usecase.GetAssetDNSRecords(r.Context(), assetID)
	if err != nil {
		if err == model.ErrAssetNotFound || errors.Is(err, model.ErrAssetNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to get dns records: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, recs)
}

// GetAssetWHOIS xử lý GET /assets/{id}/whois
func (h *ScanHandler) GetAssetWHOIS(w http.ResponseWriter, r *http.Request) {
	assetID := chi.URLParam(r, "id")
	if assetID == "" {
		respondError(w, http.StatusBadRequest, "asset id is required")
		return
	}

	rec, err := h.usecase.GetAssetWHOIS(r.Context(), assetID)
	if err != nil {
		if err == model.ErrAssetNotFound || errors.Is(err, model.ErrAssetNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to get whois record: "+err.Error())
		return
	}

	if rec == nil {
		respondJSON(w, http.StatusOK, nil)
		return
	}

	respondJSON(w, http.StatusOK, rec)
}

// GetAssetSubdomains xử lý GET /assets/{id}/subdomains
func (h *ScanHandler) GetAssetSubdomains(w http.ResponseWriter, r *http.Request) {
	assetID := chi.URLParam(r, "id")
	if assetID == "" {
		respondError(w, http.StatusBadRequest, "asset id is required")
		return
	}

	subs, err := h.usecase.GetAssetSubdomains(r.Context(), assetID)
	if err != nil {
		if err == model.ErrAssetNotFound || errors.Is(err, model.ErrAssetNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to get subdomains: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, subs)
}

// ListAllScanJobs xử lý GET /scan-jobs
func (h *ScanHandler) ListAllScanJobs(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	scanType := r.URL.Query().Get("scan_type")
	status := r.URL.Query().Get("status")
	assetQuery := r.URL.Query().Get("q")

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	jobs, total, err := h.usecase.ListAllScanJobs(r.Context(), page, limit, scanType, status, assetQuery)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list scan jobs: "+err.Error())
		return
	}

	type paginatedResponse struct {
		Data       []*model.ScanJobDetail `json:"data"`
		Total      int                    `json:"total"`
		Page       int                    `json:"page"`
		Limit      int                    `json:"limit"`
		TotalPages int                    `json:"total_pages"`
	}

	totalPages := (total + limit - 1) / limit

	respondJSON(w, http.StatusOK, paginatedResponse{
		Data:       jobs,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	})
}
