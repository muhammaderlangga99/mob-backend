-- Rollback migration: remove used column
ALTER TABLE email_verification_tokens
DROP COLUMN used;
