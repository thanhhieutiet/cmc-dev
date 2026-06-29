-- Migration 001: Tạo bảng assets
-- Nâng cấp từ in-memory storage sang PostgreSQL

CREATE TABLE IF NOT EXISTS assets (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('domain', 'ip', 'service')),
    status VARCHAR(50) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Indexes cho các query phổ biến
CREATE INDEX IF NOT EXISTS idx_assets_type ON assets(type);
CREATE INDEX IF NOT EXISTS idx_assets_status ON assets(status);
CREATE INDEX IF NOT EXISTS idx_assets_name ON assets(name);
CREATE INDEX IF NOT EXISTS idx_assets_created_at ON assets(created_at DESC);
