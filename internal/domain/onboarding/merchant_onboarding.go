package onboarding

import (
	"time"

	"gorm.io/gorm"
)

// MerchantOnboarding is the root record for a merchant onboarding flow.
// It holds the current status and step without embedding all draft payloads.
type MerchantOnboarding struct {
	// ID is the primary key for the onboarding flow (UUID).
	ID string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	// UserID links the onboarding to the authenticated merchant user.
	UserID string `gorm:"type:uuid;not null;index"`
	// Status tracks lifecycle state like DRAFT or SUBMITTED_FOR_APPROVAL.
	Status string `gorm:"type:varchar(50);not null;default:DRAFT"`
	// CurrentStep tracks UI progress to support draft continuation.
	CurrentStep string `gorm:"type:varchar(50);not null;default:business_entity"`
	// CreatedAt and UpdatedAt provide audit timestamps.
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
