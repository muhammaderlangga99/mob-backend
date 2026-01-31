-- Rollback: remove default UUID on id column
ALTER TABLE email_verification_tokens
MODIFY COLUMN id CHAR(36) NOT NULL;
