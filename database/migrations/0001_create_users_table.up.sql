-- Migration: Create users table
-- Purpose: store merchant user identities for auth flow (register/login/verify).
-- Notes: UUID PK, email unique, status/email_verified per API contract.

CREATE TABLE IF NOT EXISTS users (
  id CHAR(36) NOT NULL,
  full_name VARCHAR(150) NOT NULL,
  business_name VARCHAR(200) NOT NULL,
  email VARCHAR(120) NOT NULL,
  phone_number VARCHAR(30) NOT NULL,
  password_hash TEXT NOT NULL,
  status VARCHAR(50) NOT NULL,
  email_verified BOOLEAN NOT NULL DEFAULT FALSE,
  created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (id),
  UNIQUE KEY uk_users_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE users
ADD COLUMN deleted_at DATETIME NULL;
