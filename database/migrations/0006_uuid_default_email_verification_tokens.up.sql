-- Migration: ensure id column auto-generated via UUID
ALTER TABLE email_verification_tokens
MODIFY COLUMN id CHAR(36) NOT NULL DEFAULT (UUID());
