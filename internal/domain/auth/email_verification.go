package auth

import "time"

// EmailVerificationToken stores one-time verification links with expiry.
// Tokens are one-to-many per user to support resend flows.
type EmailVerificationToken struct {
	// Token is the primary key used in the verification link.
	Token string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	// UserID references the owner of the token.
	UserID string `gorm:"type:uuid;not null;index"`
	// ExpiresAt enforces 24-hour TTL per contract.
	ExpiresAt time.Time `gorm:"column:expired_at;not null"`
	// Used marks a token as one-time use.
	Used bool `gorm:"not null;default:false"`
	// CreatedAt records issuance time for auditing.
	CreatedAt time.Time
	// User relation for ORM joins (optional).
	User User `gorm:"foreignKey:UserID"`
}
