-- Migration 003: Thêm cột tags vào bảng assets
-- Cho phép gắn nhãn/tag cho mỗi asset để phân nhóm và lọc

ALTER TABLE assets ADD COLUMN IF NOT EXISTS tags TEXT NOT NULL DEFAULT '';

-- Index cho tìm kiếm theo tags
CREATE INDEX IF NOT EXISTS idx_assets_tags ON assets(tags);
