package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"homework-day1/internal/domain"
)

type PostgresScanRepo struct {
	db *sql.DB
}

func NewPostgresScanRepo(db *sql.DB) *PostgresScanRepo {
	return &PostgresScanRepo{db: db}
}

// CreateScanJob tạo một scan job mới
func (r *PostgresScanRepo) CreateScanJob(ctx context.Context, job *domain.ScanJob) error {
	query := `INSERT INTO scan_jobs (id, asset_id, scan_type, status, started_at, ended_at, error, results, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, query, job.ID, job.AssetID, job.ScanType, job.Status, job.StartedAt, job.EndedAt, job.Error, job.Results, job.CreatedAt)
	return err
}

// GetScanJob tìm scan job theo ID
func (r *PostgresScanRepo) GetScanJob(ctx context.Context, id string) (*domain.ScanJob, error) {
	query := `SELECT id, asset_id, scan_type, status, started_at, ended_at, error, results, created_at 
	          FROM scan_jobs WHERE id = $1`
	var job domain.ScanJob
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&job.ID, &job.AssetID, &job.ScanType, &job.Status, &job.StartedAt, &job.EndedAt, &job.Error, &job.Results, &job.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrScanJobNotFound
	}
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// UpdateScanJob cập nhật thông tin scan job
func (r *PostgresScanRepo) UpdateScanJob(ctx context.Context, job *domain.ScanJob) error {
	query := `UPDATE scan_jobs 
	          SET status = $1, ended_at = $2, error = $3, results = $4 
	          WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, job.Status, job.EndedAt, job.Error, job.Results, job.ID)
	return err
}

// ListScanJobsByAsset lấy danh sách scan job của một asset
func (r *PostgresScanRepo) ListScanJobsByAsset(ctx context.Context, assetID string) ([]*domain.ScanJob, error) {
	query := `SELECT id, asset_id, scan_type, status, started_at, ended_at, error, results, created_at 
	          FROM scan_jobs WHERE asset_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*domain.ScanJob
	for rows.Next() {
		var job domain.ScanJob
		err := rows.Scan(
			&job.ID, &job.AssetID, &job.ScanType, &job.Status, &job.StartedAt, &job.EndedAt, &job.Error, &job.Results, &job.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, &job)
	}
	return jobs, nil
}

// CreateDNSRecord lưu DNS record
func (r *PostgresScanRepo) CreateDNSRecord(ctx context.Context, record *domain.DNSRecord) error {
	query := `INSERT INTO dns_records (id, asset_id, scan_job_id, record_type, name, value, ttl, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, record.ID, record.AssetID, record.ScanJobID, record.RecordType, record.Name, record.Value, record.TTL, record.CreatedAt)
	return err
}

// GetDNSRecordsByAsset lấy danh sách DNS records của asset
func (r *PostgresScanRepo) GetDNSRecordsByAsset(ctx context.Context, assetID string) ([]*domain.DNSRecord, error) {
	query := `SELECT id, asset_id, scan_job_id, record_type, name, value, ttl, created_at 
	          FROM dns_records WHERE asset_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*domain.DNSRecord
	for rows.Next() {
		var record domain.DNSRecord
		err := rows.Scan(
			&record.ID, &record.AssetID, &record.ScanJobID, &record.RecordType, &record.Name, &record.Value, &record.TTL, &record.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	return records, nil
}

// GetDNSRecordsByScan lấy danh sách DNS records của scan job
func (r *PostgresScanRepo) GetDNSRecordsByScan(ctx context.Context, scanJobID string) ([]*domain.DNSRecord, error) {
	query := `SELECT id, asset_id, scan_job_id, record_type, name, value, ttl, created_at 
	          FROM dns_records WHERE scan_job_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, scanJobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*domain.DNSRecord
	for rows.Next() {
		var record domain.DNSRecord
		err := rows.Scan(
			&record.ID, &record.AssetID, &record.ScanJobID, &record.RecordType, &record.Name, &record.Value, &record.TTL, &record.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	return records, nil
}

// CreateWHOISRecord lưu/cập nhật WHOIS record
func (r *PostgresScanRepo) CreateWHOISRecord(ctx context.Context, record *domain.WHOISRecord) error {
	query := `INSERT INTO whois_records (id, asset_id, scan_job_id, registrar, created_date, expiry_date, name_servers, status, emails, raw_data, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	          ON CONFLICT (asset_id, scan_job_id) 
	          DO UPDATE SET 
	              registrar = EXCLUDED.registrar, 
	              created_date = EXCLUDED.created_date, 
	              expiry_date = EXCLUDED.expiry_date, 
	              name_servers = EXCLUDED.name_servers, 
	              status = EXCLUDED.status, 
	              emails = EXCLUDED.emails, 
	              raw_data = EXCLUDED.raw_data`
	_, err := r.db.ExecContext(ctx, query, record.ID, record.AssetID, record.ScanJobID, record.Registrar, record.CreatedDate, record.ExpiryDate, record.NameServers, record.Status, record.Emails, record.RawData, record.CreatedAt)
	return err
}

// GetWHOISRecordByAsset lấy WHOIS record mới nhất của asset
func (r *PostgresScanRepo) GetWHOISRecordByAsset(ctx context.Context, assetID string) (*domain.WHOISRecord, error) {
	query := `SELECT id, asset_id, scan_job_id, registrar, created_date, expiry_date, name_servers, status, emails, raw_data, created_at 
	          FROM whois_records WHERE asset_id = $1 ORDER BY created_at DESC LIMIT 1`
	var record domain.WHOISRecord
	err := r.db.QueryRowContext(ctx, query, assetID).Scan(
		&record.ID, &record.AssetID, &record.ScanJobID, &record.Registrar, &record.CreatedDate, &record.ExpiryDate, &record.NameServers, &record.Status, &record.Emails, &record.RawData, &record.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // Return nil, nil when no whois is present
	}
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// GetWHOISRecordsByScan lấy WHOIS records của scan job
func (r *PostgresScanRepo) GetWHOISRecordsByScan(ctx context.Context, scanJobID string) ([]*domain.WHOISRecord, error) {
	query := `SELECT id, asset_id, scan_job_id, registrar, created_date, expiry_date, name_servers, status, emails, raw_data, created_at 
	          FROM whois_records WHERE scan_job_id = $1`
	rows, err := r.db.QueryContext(ctx, query, scanJobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*domain.WHOISRecord
	for rows.Next() {
		var record domain.WHOISRecord
		err := rows.Scan(
			&record.ID, &record.AssetID, &record.ScanJobID, &record.Registrar, &record.CreatedDate, &record.ExpiryDate, &record.NameServers, &record.Status, &record.Emails, &record.RawData, &record.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	return records, nil
}

// CreateSubdomain lưu subdomain phát hiện được
func (r *PostgresScanRepo) CreateSubdomain(ctx context.Context, subdomain *domain.Subdomain) error {
	query := `INSERT INTO subdomains (id, asset_id, scan_job_id, name, source, is_active, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7)
	          ON CONFLICT (asset_id, name) 
	          DO UPDATE SET 
	              scan_job_id = EXCLUDED.scan_job_id, 
	              source = EXCLUDED.source, 
	              is_active = EXCLUDED.is_active, 
	              created_at = EXCLUDED.created_at`
	_, err := r.db.ExecContext(ctx, query, subdomain.ID, subdomain.AssetID, subdomain.ScanJobID, subdomain.Name, subdomain.Source, subdomain.IsActive, subdomain.CreatedAt)
	return err
}

// GetSubdomainsByAsset lấy danh sách subdomains của asset
func (r *PostgresScanRepo) GetSubdomainsByAsset(ctx context.Context, assetID string) ([]*domain.Subdomain, error) {
	query := `SELECT id, asset_id, scan_job_id, name, source, is_active, created_at 
	          FROM subdomains WHERE asset_id = $1 ORDER BY name ASC`
	rows, err := r.db.QueryContext(ctx, query, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subdomains []*domain.Subdomain
	for rows.Next() {
		var sub domain.Subdomain
		err := rows.Scan(
			&sub.ID, &sub.AssetID, &sub.ScanJobID, &sub.Name, &sub.Source, &sub.IsActive, &sub.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		subdomains = append(subdomains, &sub)
	}
	return subdomains, nil
}

// GetSubdomainsByScan lấy danh sách subdomains của scan job
func (r *PostgresScanRepo) GetSubdomainsByScan(ctx context.Context, scanJobID string) ([]*domain.Subdomain, error) {
	query := `SELECT id, asset_id, scan_job_id, name, source, is_active, created_at 
	          FROM subdomains WHERE scan_job_id = $1 ORDER BY name ASC`
	rows, err := r.db.QueryContext(ctx, query, scanJobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subdomains []*domain.Subdomain
	for rows.Next() {
		var sub domain.Subdomain
		err := rows.Scan(
			&sub.ID, &sub.AssetID, &sub.ScanJobID, &sub.Name, &sub.Source, &sub.IsActive, &sub.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		subdomains = append(subdomains, &sub)
	}
	return subdomains, nil
}

// CreateScanResult lưu kết quả scan dạng JSONB (ip, port, ssl, tech)
func (r *PostgresScanRepo) CreateScanResult(ctx context.Context, result *domain.ScanResult) error {
	query := `INSERT INTO scan_results (id, scan_job_id, asset_id, scan_type, data, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, result.ID, result.ScanJobID, result.AssetID, result.ScanType, result.Data, result.CreatedAt)
	return err
}

// GetScanResultsByScan lấy scan results của scan job
func (r *PostgresScanRepo) GetScanResultsByScan(ctx context.Context, scanJobID string) ([]*domain.ScanResult, error) {
	query := `SELECT id, scan_job_id, asset_id, scan_type, data, created_at 
	          FROM scan_results WHERE scan_job_id = $1`
	rows, err := r.db.QueryContext(ctx, query, scanJobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*domain.ScanResult
	for rows.Next() {
		var res domain.ScanResult
		var dataBytes []byte
		err := rows.Scan(
			&res.ID, &res.ScanJobID, &res.AssetID, &res.ScanType, &dataBytes, &res.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		res.Data = json.RawMessage(dataBytes)
		results = append(results, &res)
	}
	return results, nil
}

// GetScanResultsByAsset lấy scan results của asset theo scan type
func (r *PostgresScanRepo) GetScanResultsByAsset(ctx context.Context, assetID string, scanType domain.ScanType) ([]*domain.ScanResult, error) {
	query := `SELECT id, scan_job_id, asset_id, scan_type, data, created_at 
	          FROM scan_results WHERE asset_id = $1 AND scan_type = $2 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, assetID, scanType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*domain.ScanResult
	for rows.Next() {
		var res domain.ScanResult
		var dataBytes []byte
		err := rows.Scan(
			&res.ID, &res.ScanJobID, &res.AssetID, &res.ScanType, &dataBytes, &res.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		res.Data = json.RawMessage(dataBytes)
		results = append(results, &res)
	}
	return results, nil
}

// Helper to make sure it satisfies interface
var _ domain.ScanRepository = (*PostgresScanRepo)(nil)
