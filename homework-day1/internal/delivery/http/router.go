package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter tạo chi router và đăng ký tất cả các routes
func NewRouter(assetHandler *AssetHandler, healthHandler *HealthHandler) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	// Health check (Bài 5)
	r.Get("/health", healthHandler.HealthCheck)

	// Asset routes
	r.Route("/assets", func(r chi.Router) {
		// Các route cụ thể phải đặt TRƯỚC route có param {id}
		r.Get("/stats", assetHandler.GetStats)       // Bài 1.1
		r.Get("/count", assetHandler.CountAssets)     // Bài 1.2
		r.Get("/search", assetHandler.SearchAssets)   // Bài 7 - Bonus
		r.Post("/batch", assetHandler.BatchCreate)    // Bài 2
		r.Delete("/batch", assetHandler.BatchDelete)  // Bài 3

		// CRUD cơ bản
		r.Get("/", assetHandler.ListAssets)           // Bài 6 - Bonus
		r.Post("/", assetHandler.CreateAsset)         // Tạo 1 asset
		r.Get("/{id}", assetHandler.GetAssetByID)     // Lấy asset theo ID
	})

	return r
}
