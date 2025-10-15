-- Migration: add index on auth_refresh_tokens for faster lookup by user and expiry
-- This file includes variants for MySQL/MariaDB and PostgreSQL. Apply the one appropriate for your DB.

-- PostgreSQL
-- CREATE INDEX IF NOT EXISTS idx_auth_refresh_tokens_user_expires ON auth_refresh_tokens (user_id, expires_at);

-- MySQL / MariaDB
-- ALTER TABLE auth_refresh_tokens ADD INDEX idx_auth_refresh_tokens_user_expires (user_id, expires_at);

-- Ensure token column can store bcrypt hashes (at least 60 chars). If token is VARCHAR(255) or TEXT, it's already fine.
-- Example for MySQL: ALTER TABLE auth_refresh_tokens MODIFY token VARCHAR(255) NOT NULL;
