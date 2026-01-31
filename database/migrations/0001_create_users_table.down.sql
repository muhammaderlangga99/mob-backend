-- Migration: Drop users table
-- Purpose: rollback for users table creation.

DROP TABLE IF EXISTS users;

ALTER TABLE users
DROP COLUMN deleted_at;