package onboarding

import (
	"time"

	"gorm.io/gorm"
)

// TermsAcceptance stores legal acceptance records for onboarding.
// This is immutable history; new acceptance records should be appended.
type TermsAcceptance struct {
	ID                   string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	MerchantOnboardingID string `gorm:"type:uuid;not null;index"`
	// PayloadJSON stores accepted_at, accepted_by, terms_version, user_agent, etc.
	PayloadJSON string `gorm:"type:jsonb;not null"`
	CreatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	// UpdatedAt omitted on purpose to emphasize immutability of records.
}
