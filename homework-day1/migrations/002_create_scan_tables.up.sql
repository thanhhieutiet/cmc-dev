-- Migration 002: Tạo các bảng phục vụ EASM Scanning

-- Bảng scan_jobs: theo dõi các tác vụ quét (async job pattern)
CREATE TABLE IF NOT EXISTS scan_jobs (
    id UUID PRIMARY KEY,
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    scan_type VARCHAR(50) NOT NULL
        CHECK (scan_type IN ('dns', 'whois', 'subdomain', 'ip', 'port', 'ssl', 'tech', 'all')),
    status VARCHAR(50) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'running', 'completed', 'failed', 'partial')),
    started_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    ended_at TIMESTAMP WITH TIME ZONE,
    error TEXT DEFAULT '',
    results INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_scan_jobs_asset_id ON scan_jobs(asset_id);
CREATE INDEX IF NOT EXISTS idx_scan_jobs_status ON scan_jobs(status);
CREATE INDEX IF NOT EXISTS idx_scan_jobs_created_at ON scan_jobs(created_at DESC);

-- Bảng subdomains: lưu subdomain phát hiện được
CREATE TABLE IF NOT EXISTS subdomains (
    id UUID PRIMARY KEY,
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    scan_job_id UUID NOT NULL REFERENCES scan_jobs(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    source VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE (asset_id, name)
);

CREATE INDEX IF NOT EXISTS idx_subdomains_asset_id ON subdomains(asset_id);
CREATE INDEX IF NOT EXISTS idx_subdomains_scan_job_id ON subdomains(scan_job_id);

-- Bảng dns_records: lưu DNS record
CREATE TABLE IF NOT EXISTS dns_records (
    id UUID PRIMARY KEY,
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    scan_job_id UUID NOT NULL REFERENCES scan_jobs(id) ON DELETE CASCADE,
    record_type VARCHAR(10) NOT NULL,
    name VARCHAR(255) NOT NULL,
    value TEXT NOT NULL,
    ttl INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_dns_records_asset_id ON dns_records(asset_id);
CREATE INDEX IF NOT EXISTS idx_dns_records_scan_job_id ON dns_records(scan_job_id);

-- Bảng whois_records: lưu WHOIS info
CREATE TABLE IF NOT EXISTS whois_records (
    id UUID PRIMARY KEY,
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    scan_job_id UUID NOT NULL REFERENCES scan_jobs(id) ON DELETE CASCADE,
    registrar TEXT DEFAULT '',
    created_date TIMESTAMP WITH TIME ZONE,
    expiry_date TIMESTAMP WITH TIME ZONE,
    name_servers TEXT DEFAULT '',
    status TEXT DEFAULT '',
    emails TEXT DEFAULT '',
    raw_data TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE (asset_id, scan_job_id)
);

CREATE INDEX IF NOT EXISTS idx_whois_records_asset_id ON whois_records(asset_id);

-- Bảng scan_results: lưu kết quả scan dạng JSONB cho các scan type mới (ip, port, ssl, tech)
CREATE TABLE IF NOT EXISTS scan_results (
    id UUID PRIMARY KEY,
    scan_job_id UUID NOT NULL REFERENCES scan_jobs(id) ON DELETE CASCADE,
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    scan_type VARCHAR(50) NOT NULL,
    data JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_scan_results_scan_job_id ON scan_results(scan_job_id);
CREATE INDEX IF NOT EXISTS idx_scan_results_asset_id ON scan_results(asset_id);
CREATE INDEX IF NOT EXISTS idx_scan_results_scan_type ON scan_results(scan_type);
