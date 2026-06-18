package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"homework-day1/internal/domain"
)

type PostgresAssetRepo struct {
	db *sql.DB
}

func NewPostgresAssetRepo(db *sql.DB) *PostgresAssetRepo {
	return &PostgresAssetRepo{db: db}
}

// Create thêm một asset vào DB
func (r *PostgresAssetRepo) Create(ctx context.Context, asset *domain.Asset) error {
	query := `INSERT INTO assets (id, name, type, status, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, asset.ID, asset.Name, asset.Type, asset.Status, asset.CreatedAt, asset.UpdatedAt)
	return err
}

// GetByID tìm asset theo ID
func (r *PostgresAssetRepo) GetByID(ctx context.Context, id string) (*domain.Asset, error) {
	query := `SELECT id, name, type, status, created_at, updated_at 
	          FROM assets WHERE id = $1`
	var asset domain.Asset
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&asset.ID, &asset.Name, &asset.Type, &asset.Status, &asset.CreatedAt, &asset.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrAssetNotFound
	}
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// Delete xóa một asset theo ID
func (r *PostgresAssetRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM assets WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrAssetNotFound
	}
	return nil
}

// BatchCreate thêm nhiều assets cùng lúc
func (r *PostgresAssetRepo) BatchCreate(ctx context.Context, assets []*domain.Asset) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO assets (id, name, type, status, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, asset := range assets {
		_, err := stmt.ExecContext(ctx, asset.ID, asset.Name, asset.Type, asset.Status, asset.CreatedAt, asset.UpdatedAt)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

// BatchDelete xóa nhiều assets theo danh sách IDs
func (r *PostgresAssetRepo) BatchDelete(ctx context.Context, ids []string) (int, int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()

	deleted := 0
	notFound := 0

	query := `DELETE FROM assets WHERE id = $1`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return 0, 0, err
	}
	defer stmt.Close()

	for _, id := range ids {
		res, err := stmt.ExecContext(ctx, id)
		if err != nil {
			return deleted, notFound, err
		}
		rows, err := res.RowsAffected()
		if err != nil {
			return deleted, notFound, err
		}
		if rows > 0 {
			deleted++
		} else {
			notFound++
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, 0, err
	}
	return deleted, notFound, nil
}

// GetStats trả về thống kê tổng quan về các assets
func (r *PostgresAssetRepo) GetStats(ctx context.Context) (*domain.StatsResponse, error) {
	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM assets").Scan(&total)
	if err != nil {
		return nil, err
	}

	stats := &domain.StatsResponse{
		Total:    total,
		ByType:   make(map[string]int),
		ByStatus: make(map[string]int),
	}

	typeRows, err := r.db.QueryContext(ctx, "SELECT type, COUNT(*) FROM assets GROUP BY type")
	if err != nil {
		return nil, err
	}
	defer typeRows.Close()
	for typeRows.Next() {
		var t string
		var count int
		if err := typeRows.Scan(&t, &count); err != nil {
			return nil, err
		}
		stats.ByType[t] = count
	}

	statusRows, err := r.db.QueryContext(ctx, "SELECT status, COUNT(*) FROM assets GROUP BY status")
	if err != nil {
		return nil, err
	}
	defer statusRows.Close()
	for statusRows.Next() {
		var s string
		var count int
		if err := statusRows.Scan(&s, &count); err != nil {
			return nil, err
		}
		stats.ByStatus[s] = count
	}

	return stats, nil
}

// Count đếm số assets theo bộ lọc type và status
func (r *PostgresAssetRepo) Count(ctx context.Context, assetType, status string) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM assets WHERE 1=1"
	var args []interface{}
	placeholderIdx := 1

	if assetType != "" {
		query += fmt.Sprintf(" AND type = $%d", placeholderIdx)
		args = append(args, assetType)
		placeholderIdx++
	}
	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", placeholderIdx)
		args = append(args, status)
		placeholderIdx++
	}

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

// List trả về danh sách assets có phân trang và lọc
func (r *PostgresAssetRepo) List(ctx context.Context, page, limit int, assetType, status string) ([]*domain.Asset, int, error) {
	total, err := r.Count(ctx, assetType, status)
	if err != nil {
		return nil, 0, err
	}

	query := "SELECT id, name, type, status, created_at, updated_at FROM assets WHERE 1=1"
	var args []interface{}
	placeholderIdx := 1

	if assetType != "" {
		query += fmt.Sprintf(" AND type = $%d", placeholderIdx)
		args = append(args, assetType)
		placeholderIdx++
	}
	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", placeholderIdx)
		args = append(args, status)
		placeholderIdx++
	}

	query += " ORDER BY created_at DESC"

	offset := (page - 1) * limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", placeholderIdx, placeholderIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var assets []*domain.Asset
	for rows.Next() {
		var asset domain.Asset
		err := rows.Scan(
			&asset.ID, &asset.Name, &asset.Type, &asset.Status, &asset.CreatedAt, &asset.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		assets = append(assets, &asset)
	}

	return assets, total, nil
}

// Search tìm kiếm assets theo tên
func (r *PostgresAssetRepo) Search(ctx context.Context, query string, maxResults int) ([]*domain.Asset, error) {
	sqlQuery := "SELECT id, name, type, status, created_at, updated_at FROM assets WHERE name ILIKE $1 ORDER BY created_at DESC LIMIT $2"
	rows, err := r.db.QueryContext(ctx, sqlQuery, "%"+query+"%", maxResults)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []*domain.Asset
	for rows.Next() {
		var asset domain.Asset
		err := rows.Scan(
			&asset.ID, &asset.Name, &asset.Type, &asset.Status, &asset.CreatedAt, &asset.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

// TotalCount trả về tổng số assets hiện có (dùng cho Health Check)
func (r *PostgresAssetRepo) TotalCount() int {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM assets").Scan(&count)
	if err != nil {
		return 0
	}
	return count
}
