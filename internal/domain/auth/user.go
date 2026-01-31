package auth

import (
	"time"

	"gorm.io/gorm"
)

// UserStatus defines lifecycle state for a merchant user in the auth flow.
type UserStatus string

const (
	// StatusPendingEmailVerification blocks login until email is verified.
	StatusPendingEmailVerification UserStatus = "PENDING_EMAIL_VERIFICATION"
	// StatusActive allows login and dashboard access after verification.
	StatusActive UserStatus = "ACTIVE"
)

// User stores identity and auth profile for merchant onboarding.
// This table is the single source of truth for auth status and email verification.
type User struct {
	// ID is the primary key for user identity (UUID).
	ID string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	// FullName is required by the register payload and shown in dashboard.
	FullName string `gorm:"type:varchar(150);not null"`
	// BusinessName is required by the register payload and shown in dashboard.
	BusinessName string `gorm:"type:varchar(200);not null"`
	// Email is used for login and verification flow; must be unique.
	Email string `gorm:"type:varchar(120);not null;uniqueIndex"`
	// PhoneNumber is stored for dashboard profile.
	PhoneNumber string `gorm:"type:varchar(30);not null"`
	// PasswordHash stores bcrypt hash; raw password never persisted.
	PasswordHash string `gorm:"type:text;not null"`
	// Status tracks PENDING_EMAIL_VERIFICATION vs ACTIVE per contract.
	Status UserStatus `gorm:"type:varchar(50);not null"`
	// EmailVerified mirrors verification state for quick checks.
	EmailVerified bool `gorm:"not null;default:false"`
	// VerificationTokens links to one-time email verification tokens.
	VerificationTokens []EmailVerificationToken `gorm:"foreignKey:UserID"`
	// Timestamps follow standard auditing.
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
