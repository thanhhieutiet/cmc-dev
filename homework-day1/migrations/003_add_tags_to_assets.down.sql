-- Rollback Migration 003: Xóa cột tags khỏi bảng assets

ALTER TABLE assets DROP COLUMN IF EXISTS tags;
