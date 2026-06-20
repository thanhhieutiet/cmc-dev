package postgres

import (
	"context"
	"database/sql"

	"homework-day1/internal/model"
)

type PostgresAlertRepo struct {
	db *sql.DB
}

func NewPostgresAlertRepo(db *sql.DB) *PostgresAlertRepo {
	return &PostgresAlertRepo{db: db}
}

func (r *PostgresAlertRepo) Create(ctx context.Context, alert *model.Alert) error {
	query := `INSERT INTO alerts (id, asset_id, type, message, is_read, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, alert.ID, alert.AssetID, alert.Type, alert.Message, alert.IsRead, alert.CreatedAt)
	return err
}

func (r *PostgresAlertRepo) ListUnread(ctx context.Context) ([]*model.Alert, error) {
	query := `SELECT id, asset_id, type, message, is_read, created_at 
	          FROM alerts WHERE is_read = false ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []*model.Alert
	for rows.Next() {
		var a model.Alert
		if err := rows.Scan(&a.ID, &a.AssetID, &a.Type, &a.Message, &a.IsRead, &a.CreatedAt); err != nil {
			return nil, err
		}
		alerts = append(alerts, &a)
	}
	return alerts, nil
}

func (r *PostgresAlertRepo) ListAll(ctx context.Context) ([]*model.Alert, error) {
	query := `SELECT id, asset_id, type, message, is_read, created_at 
	          FROM alerts ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []*model.Alert
	for rows.Next() {
		var a model.Alert
		if err := rows.Scan(&a.ID, &a.AssetID, &a.Type, &a.Message, &a.IsRead, &a.CreatedAt); err != nil {
			return nil, err
		}
		alerts = append(alerts, &a)
	}
	return alerts, nil
}

func (r *PostgresAlertRepo) ListByAsset(ctx context.Context, assetID string) ([]*model.Alert, error) {
	query := `SELECT id, asset_id, type, message, is_read, created_at 
	          FROM alerts WHERE asset_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []*model.Alert
	for rows.Next() {
		var a model.Alert
		if err := rows.Scan(&a.ID, &a.AssetID, &a.Type, &a.Message, &a.IsRead, &a.CreatedAt); err != nil {
			return nil, err
		}
		alerts = append(alerts, &a)
	}
	return alerts, nil
}

func (r *PostgresAlertRepo) MarkAsRead(ctx context.Context, id string) error {
	query := `UPDATE alerts SET is_read = true WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresAlertRepo) MarkAllAsRead(ctx context.Context) error {
	query := `UPDATE alerts SET is_read = true WHERE is_read = false`
	_, err := r.db.ExecContext(ctx, query)
	return err
}
