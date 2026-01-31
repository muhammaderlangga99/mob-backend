-- Migration: Add used column to email_verification_tokens
-- Ensures ORM field matches schema for one-time tokens.

ALTER TABLE email_verification_tokens
ADD COLUMN used BOOLEAN NOT NULL DEFAULT FALSE;
