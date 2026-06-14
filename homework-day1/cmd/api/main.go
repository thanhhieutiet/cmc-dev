package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	deliveryHttp "homework-day1/internal/delivery/http"
	"homework-day1/internal/repository/memory"
	"homework-day1/internal/usecase"
)

func main() {
	// Ghi nhận thời gian bắt đầu để tính uptime (Bài 5)
	startTime := time.Now()

	// --- Khởi tạo các Layer theo Clean Architecture ---

	// Layer 1: Repository (In-memory storage)
	assetRepo := memory.NewInMemoryAssetRepo()

	// Layer 2: Usecase (Business logic)
	assetUsecase := usecase.NewAssetUsecase(assetRepo)

	// Layer 3: Delivery (HTTP handlers)
	assetHandler := deliveryHttp.NewAssetHandler(assetUsecase)
	healthHandler := deliveryHttp.NewHealthHandler(assetUsecase, startTime)

	// Tạo router
	router := deliveryHttp.NewRouter(assetHandler, healthHandler)

	// --- Khởi chạy HTTP Server ---
	port := ":8080"
	fmt.Println("======================================")
	fmt.Println("  Asset Management API")
	fmt.Println("  Clean Architecture + In-Memory")
	fmt.Printf("  Server started on http://localhost%s\n", port)
	fmt.Println("======================================")
	fmt.Println("")
	fmt.Println("Available endpoints:")
	fmt.Println("  GET    /health          - Health check (Bài 5)")
	fmt.Println("  GET    /assets/stats    - Statistics (Bài 1.1)")
	fmt.Println("  GET    /assets/count    - Count with filters (Bài 1.2)")
	fmt.Println("  POST   /assets/batch    - Batch create (Bài 2)")
	fmt.Println("  DELETE /assets/batch    - Batch delete (Bài 3)")
	fmt.Println("  POST   /assets          - Create single asset")
	fmt.Println("  GET    /assets/{id}     - Get asset by ID")
	fmt.Println("  GET    /assets          - List with pagination (Bài 6)")
	fmt.Println("  GET    /assets/search   - Search by name (Bài 7)")
	fmt.Println("")

	log.Fatal(http.ListenAndServe(port, router))
}
