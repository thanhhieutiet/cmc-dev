package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"homework-day1/internal/domain"
)

type ScanHandler struct {
	usecase domain.ScanUsecase
}

func NewScanHandler(uc domain.ScanUsecase) *ScanHandler {
	return &ScanHandler{usecase: uc}
}

// StartScan xử lý POST /assets/{id}/scan
func (h *ScanHandler) StartScan(w http.ResponseWriter, r *http.Request) {
	assetID := chi.URLParam(r, "id")
	if assetID == "" {
		respondError(w, http.StatusBadRequest, "asset id is required")
		return
	}

	var req domain.StartScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	job, err := h.usecase.StartScan(r.Context(), assetID, req.ScanType)
	if err != nil {
		if err == domain.ErrAssetNotFound || errors.Is(err, domain.ErrAssetNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		if err == domain.ErrInvalidScanType || err == domain.ErrScanNotSupported || err == domain.ErrPortScanUnauthorized ||
			errors.Is(err, domain.ErrInvalidScanType) || errors.Is(err, domain.ErrScanNotSupported) || errors.Is(err, domain.ErrPortScanUnauthorized) {
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
		if err == domain.ErrScanJobNotFound || errors.Is(err, domain.ErrScanJobNotFound) {
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
		if err == domain.ErrScanJobNotFound || errors.Is(err, domain.ErrScanJobNotFound) {
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
		if err == domain.ErrAssetNotFound || errors.Is(err, domain.ErrAssetNotFound) {
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
		if err == domain.ErrAssetNotFound || errors.Is(err, domain.ErrAssetNotFound) {
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
		if err == domain.ErrAssetNotFound || errors.Is(err, domain.ErrAssetNotFound) {
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
		if err == domain.ErrAssetNotFound || errors.Is(err, domain.ErrAssetNotFound) {
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
		if err == domain.ErrAssetNotFound || errors.Is(err, domain.ErrAssetNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to get subdomains: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, subs)
}
