-- Migration 004: Rollback alerts table and auto_scan

DROP TABLE IF EXISTS alerts;
ALTER TABLE assets DROP COLUMN IF EXISTS auto_scan;
