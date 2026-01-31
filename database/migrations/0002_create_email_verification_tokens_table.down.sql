-- Migration: Drop email_verification_tokens table
-- Purpose: rollback for verification token storage.

DROP TABLE IF EXISTS email_verification_tokens;
