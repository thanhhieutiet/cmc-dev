package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"homework-day1/internal/model"
)

type AlertHandler struct {
	repo model.AlertRepository
}

func NewAlertHandler(repo model.AlertRepository) *AlertHandler {
	return &AlertHandler{repo: repo}
}

func (h *AlertHandler) GetUnreadAlerts(w http.ResponseWriter, r *http.Request) {
	alerts, err := h.repo.ListUnread(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}

func (h *AlertHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	err := h.repo.MarkAsRead(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AlertHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	err := h.repo.MarkAllAsRead(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AlertHandler) ListAllAlerts(w http.ResponseWriter, r *http.Request) {
	alerts, err := h.repo.ListAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}
