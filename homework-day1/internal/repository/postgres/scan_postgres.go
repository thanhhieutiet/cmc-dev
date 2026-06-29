package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"homework-day1/internal/model"
)

type PostgresScanRepo struct {
	db *sql.DB
}

func NewPostgresScanRepo(db *sql.DB) *PostgresScanRepo {
	return &PostgresScanRepo{db: db}
}

// CreateScanJob tạo một scan job mới
func (r *PostgresScanRepo) CreateScanJob(ctx context.Context, job *model.ScanJob) error {
	query := `INSERT INTO scan_jobs (id, asset_id, scan_type, status, started_at, ended_at, error, results, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, query, job.ID, job.AssetID, job.ScanType, job.Status, job.StartedAt, job.EndedAt, job.Error, job.Results, job.CreatedAt)
	return err
}

// GetScanJob tìm scan job theo ID
func (r *PostgresScanRepo) GetScanJob(ctx context.Context, id string) (*model.ScanJob, error) {
	query := `SELECT id, asset_id, scan_type, status, started_at, ended_at, error, results, created_at 
	          FROM scan_jobs WHERE id = $1`
	var job model.ScanJob
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&job.ID, &job.AssetID, &job.ScanType, &job.Status, &job.StartedAt, &job.EndedAt, &job.Error, &job.Results, &job.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, model.ErrScanJobNotFound
	}
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// UpdateScanJob cập nhật thông tin scan job
func (r *PostgresScanRepo) UpdateScanJob(ctx context.Context, job *model.ScanJob) error {
	query := `UPDATE scan_jobs 
	          SET status = $1, ended_at = $2, error = $3, results = $4 
	          WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, job.Status, job.EndedAt, job.Error, job.Results, job.ID)
	return err
}

// ListScanJobsByAsset lấy danh sách scan job của một asset
func (r *PostgresScanRepo) ListScanJobsByAsset(ctx context.Context, assetID string) ([]*model.ScanJob, error) {
	query := `SELECT id, asset_id, scan_type, status, started_at, ended_at, error, results, created_at 
	          FROM scan_jobs WHERE asset_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*model.ScanJob
	for rows.Next() {
		var job model.ScanJob
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
func (r *PostgresScanRepo) CreateDNSRecord(ctx context.Context, record *model.DNSRecord) error {
	query := `INSERT INTO dns_records (id, asset_id, scan_job_id, record_type, name, value, ttl, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, record.ID, record.AssetID, record.ScanJobID, record.RecordType, record.Name, record.Value, record.TTL, record.CreatedAt)
	return err
}

// GetDNSRecordsByAsset lấy danh sách DNS records của asset
func (r *PostgresScanRepo) GetDNSRecordsByAsset(ctx context.Context, assetID string) ([]*model.DNSRecord, error) {
	query := `SELECT id, asset_id, scan_job_id, record_type, name, value, ttl, created_at 
	          FROM dns_records WHERE asset_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*model.DNSRecord
	for rows.Next() {
		var record model.DNSRecord
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
func (r *PostgresScanRepo) GetDNSRecordsByScan(ctx context.Context, scanJobID string) ([]*model.DNSRecord, error) {
	query := `SELECT id, asset_id, scan_job_id, record_type, name, value, ttl, created_at 
	          FROM dns_records WHERE scan_job_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, scanJobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*model.DNSRecord
	for rows.Next() {
		var record model.DNSRecord
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
func (r *PostgresScanRepo) CreateWHOISRecord(ctx context.Context, record *model.WHOISRecord) error {
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
func (r *PostgresScanRepo) GetWHOISRecordByAsset(ctx context.Context, assetID string) (*model.WHOISRecord, error) {
	query := `SELECT id, asset_id, scan_job_id, registrar, created_date, expiry_date, name_servers, status, emails, raw_data, created_at 
	          FROM whois_records WHERE asset_id = $1 ORDER BY created_at DESC LIMIT 1`
	var record model.WHOISRecord
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
func (r *PostgresScanRepo) GetWHOISRecordsByScan(ctx context.Context, scanJobID string) ([]*model.WHOISRecord, error) {
	query := `SELECT id, asset_id, scan_job_id, registrar, created_date, expiry_date, name_servers, status, emails, raw_data, created_at 
	          FROM whois_records WHERE scan_job_id = $1`
	rows, err := r.db.QueryContext(ctx, query, scanJobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*model.WHOISRecord
	for rows.Next() {
		var record model.WHOISRecord
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
func (r *PostgresScanRepo) CreateSubdomain(ctx context.Context, subdomain *model.Subdomain) error {
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
func (r *PostgresScanRepo) GetSubdomainsByAsset(ctx context.Context, assetID string) ([]*model.Subdomain, error) {
	query := `SELECT id, asset_id, scan_job_id, name, source, is_active, created_at 
	          FROM subdomains WHERE asset_id = $1 ORDER BY name ASC`
	rows, err := r.db.QueryContext(ctx, query, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subdomains []*model.Subdomain
	for rows.Next() {
		var sub model.Subdomain
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
func (r *PostgresScanRepo) GetSubdomainsByScan(ctx context.Context, scanJobID string) ([]*model.Subdomain, error) {
	query := `SELECT id, asset_id, scan_job_id, name, source, is_active, created_at 
	          FROM subdomains WHERE scan_job_id = $1 ORDER BY name ASC`
	rows, err := r.db.QueryContext(ctx, query, scanJobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subdomains []*model.Subdomain
	for rows.Next() {
		var sub model.Subdomain
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
func (r *PostgresScanRepo) CreateScanResult(ctx context.Context, result *model.ScanResult) error {
	query := `INSERT INTO scan_results (id, scan_job_id, asset_id, scan_type, data, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, result.ID, result.ScanJobID, result.AssetID, result.ScanType, result.Data, result.CreatedAt)
	return err
}

// GetScanResultsByScan lấy scan results của scan job
func (r *PostgresScanRepo) GetScanResultsByScan(ctx context.Context, scanJobID string) ([]*model.ScanResult, error) {
	query := `SELECT id, scan_job_id, asset_id, scan_type, data, created_at 
	          FROM scan_results WHERE scan_job_id = $1`
	rows, err := r.db.QueryContext(ctx, query, scanJobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*model.ScanResult
	for rows.Next() {
		var res model.ScanResult
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
func (r *PostgresScanRepo) GetScanResultsByAsset(ctx context.Context, assetID string, scanType model.ScanType) ([]*model.ScanResult, error) {
	query := `SELECT id, scan_job_id, asset_id, scan_type, data, created_at 
	          FROM scan_results WHERE asset_id = $1 AND scan_type = $2 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, assetID, scanType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*model.ScanResult
	for rows.Next() {
		var res model.ScanResult
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

// ListAllScanJobs lấy danh sách tất cả các scan jobs kèm thông tin asset và phân trang
func (r *PostgresScanRepo) ListAllScanJobs(ctx context.Context, page, limit int, scanType, status, assetQuery string) ([]*model.ScanJobDetail, int, error) {
	whereClauses := []string{"1=1"}
	args := []interface{}{}
	argIndex := 1

	if scanType != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("sj.scan_type = $%d", argIndex))
		args = append(args, scanType)
		argIndex++
	}
	if status != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("sj.status = $%d", argIndex))
		args = append(args, status)
		argIndex++
	}
	if assetQuery != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("a.name ILIKE $%d", argIndex))
		args = append(args, "%"+assetQuery+"%")
		argIndex++
	}

	whereSQL := strings.Join(whereClauses, " AND ")

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM scan_jobs sj
		JOIN assets a ON sj.asset_id = a.id
		WHERE %s`, whereSQL)

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	dataQuery := fmt.Sprintf(`
		SELECT sj.id, sj.asset_id, sj.scan_type, sj.status, sj.started_at, sj.ended_at, sj.error, sj.results, sj.created_at,
		       a.name AS asset_name, a.type AS asset_type
		FROM scan_jobs sj
		JOIN assets a ON sj.asset_id = a.id
		WHERE %s
		ORDER BY sj.created_at DESC
		LIMIT $%d OFFSET $%d`, whereSQL, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var details []*model.ScanJobDetail
	for rows.Next() {
		var detail model.ScanJobDetail
		err := rows.Scan(
			&detail.ID, &detail.AssetID, &detail.ScanType, &detail.Status, &detail.StartedAt, &detail.EndedAt, &detail.Error, &detail.Results, &detail.CreatedAt,
			&detail.AssetName, &detail.AssetType,
		)
		if err != nil {
			return nil, 0, err
		}
		details = append(details, &detail)
	}

	return details, total, nil
}

// Helper to make sure it satisfies interface
var _ model.ScanRepository = (*PostgresScanRepo)(nil)
