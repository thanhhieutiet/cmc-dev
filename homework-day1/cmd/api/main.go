package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"homework-day1/internal/config"
	deliveryHttp "homework-day1/internal/delivery/http"
	"homework-day1/internal/repository/postgres"
	"homework-day1/internal/usecase"
)

func main() {
	startTime := time.Now()

	// Load configuration
	cfg := config.Load()

	// Connect to PostgreSQL with retry
	var db *sql.DB
	var err error
	dsn := cfg.DB.DSN()

	log.Printf("Connecting to database at %s:%s...", cfg.DB.Host, cfg.DB.Port)
	for i := 1; i <= 5; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			break
		}
		log.Printf("[Attempt %d/5] Failed to connect to DB: %v. Retrying in 2 seconds...", i, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Could not connect to database after 5 attempts: %v", err)
	}
	defer db.Close()
	log.Println("Database connection established successfully.")

	// Auto-run migrations
	migrationsDir := "./migrations"
	if err := runMigrations(db, migrationsDir); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// --- Initialize Clean Architecture Layers ---

	// Repositories
	assetRepo := postgres.NewPostgresAssetRepo(db)
	scanRepo := postgres.NewPostgresScanRepo(db)

	// Usecases
	assetUsecase := usecase.NewAssetUsecase(assetRepo)
	scanUsecase := usecase.NewScanUsecase(assetRepo, scanRepo)

	// Handlers
	assetHandler := deliveryHttp.NewAssetHandler(assetUsecase)
	scanHandler := deliveryHttp.NewScanHandler(scanUsecase)
	healthHandler := deliveryHttp.NewHealthHandler(assetUsecase, startTime)

	// Router
	router := deliveryHttp.NewRouter(assetHandler, scanHandler, healthHandler)

	// Start server
	port := ":" + cfg.Server.Port
	fmt.Println("==================================================")
	fmt.Println("  EASM Digital Asset Management & Scanner API")
	fmt.Println("  Clean Architecture + PostgreSQL Storage")
	fmt.Printf("  Server started on http://localhost%s\n", port)
	fmt.Println("==================================================")
	fmt.Println("")
	fmt.Println("Available endpoints:")
	fmt.Println("  [Health & Stats]")
	fmt.Println("    GET    /health                - Health check & database stats")
	fmt.Println("    GET    /assets/stats          - Statistics by type and status")
	fmt.Println("    GET    /assets/count          - Count filtered assets")
	fmt.Println("")
	fmt.Println("  [Asset Operations]")
	fmt.Println("    POST   /assets                - Create single asset")
	fmt.Println("    GET    /assets/{id}           - Get asset by ID")
	fmt.Println("    GET    /assets                - List assets with pagination and filters")
	fmt.Println("    GET    /assets/search         - Search assets by name")
	fmt.Println("    POST   /assets/batch          - Batch create assets")
	fmt.Println("    DELETE /assets/batch          - Batch delete assets")
	fmt.Println("")
	fmt.Println("  [Scan Operations]")
	fmt.Println("    POST   /assets/{id}/scan      - Trigger scanning for an asset")
	fmt.Println("    GET    /assets/{id}/scans     - List all scans of an asset")
	fmt.Println("    GET    /assets/{id}/results   - Get combined scan results of an asset")
	fmt.Println("    GET    /assets/{id}/dns       - Get DNS scan records for asset")
	fmt.Println("    GET    /assets/{id}/whois     - Get WHOIS scan record for asset")
	fmt.Println("    GET    /assets/{id}/subdomains- Get Subdomain scan records for asset")
	fmt.Println("    GET    /scan-jobs/{id}        - Check scan job status")
	fmt.Println("    GET    /scan-jobs/{id}/results- Get results for a specific scan job")
	fmt.Println("==================================================")
	fmt.Println("")

	log.Fatal(http.ListenAndServe(port, router))
}

func runMigrations(db *sql.DB, migrationsDir string) error {
	// Create schema_migrations table if not exists
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	)`)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	// Read migration files
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var upMigrationFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".up.sql") {
			upMigrationFiles = append(upMigrationFiles, f.Name())
		}
	}
	sort.Strings(upMigrationFiles)

	for _, filename := range upMigrationFiles {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", filename).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if exists {
			continue
		}

		log.Printf("Applying migration: %s", filename)
		content, err := os.ReadFile(filepath.Join(migrationsDir, filename))
		if err != nil {
			return fmt.Errorf("failed to read migration content: %w", err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to start migration transaction: %w", err)
		}

		_, err = tx.Exec(string(content))
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration SQL: %w", err)
		}

		_, err = tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", filename)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to log migration: %w", err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration transaction: %w", err)
		}
		log.Printf("Successfully applied migration: %s", filename)
	}

	return nil
}
