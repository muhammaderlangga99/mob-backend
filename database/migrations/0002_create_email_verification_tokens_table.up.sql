-- Migration: Create email_verification_tokens table
-- Purpose: store one-time verification tokens with expiry for email verification.
-- Notes: token unique, used_at nullable for one-time use tracking.

CREATE TABLE IF NOT EXISTS email_verification_tokens (
  id CHAR(36) NOT NULL,
  user_id CHAR(36) NOT NULL,
  token CHAR(36) NOT NULL,
  expired_at DATETIME(6) NOT NULL,
  used_at DATETIME(6) NULL,
  created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  PRIMARY KEY (id),
  UNIQUE KEY uk_email_verification_tokens_token (token),
  KEY idx_email_verification_tokens_user_id (user_id),
  CONSTRAINT fk_email_verification_tokens_user
    FOREIGN KEY (user_id) REFERENCES users(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
