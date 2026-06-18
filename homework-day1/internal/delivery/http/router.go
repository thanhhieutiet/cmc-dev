package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// CORSMiddleware xử lý chia sẻ tài nguyên nguồn gốc chéo (CORS) cho phép frontend gọi API
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func NewRouter(assetHandler *AssetHandler, scanHandler *ScanHandler, healthHandler *HealthHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(CORSMiddleware)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	r.Get("/health", healthHandler.HealthCheck)

	r.Route("/assets", func(r chi.Router) {
		r.Get("/stats", assetHandler.GetStats)
		r.Get("/count", assetHandler.CountAssets)
		r.Get("/search", assetHandler.SearchAssets)
		r.Post("/batch", assetHandler.BatchCreate)
		r.Delete("/batch", assetHandler.BatchDelete)
		r.Get("/", assetHandler.ListAssets)
		r.Post("/", assetHandler.CreateAsset)
		r.Get("/{id}", assetHandler.GetAssetByID)

		// Scan routes
		r.Post("/{id}/scan", scanHandler.StartScan)
		r.Get("/{id}/scans", scanHandler.ListScanJobs)
		r.Get("/{id}/results", scanHandler.GetAssetAllResults)
		r.Get("/{id}/dns", scanHandler.GetAssetDNS)
		r.Get("/{id}/whois", scanHandler.GetAssetWHOIS)
		r.Get("/{id}/subdomains", scanHandler.GetAssetSubdomains)
	})

	r.Route("/scan-jobs", func(r chi.Router) {
		r.Get("/{id}", scanHandler.GetScanJob)
		r.Get("/{id}/results", scanHandler.GetScanResults)
	})

	return r
}
