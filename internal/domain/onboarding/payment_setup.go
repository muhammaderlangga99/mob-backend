package onboarding

import (
	"time"

	"gorm.io/gorm"
)

// PaymentSetup stores selected payment features per device in draft state.
// Payload is flexible to align with device → payment features mapping.
type PaymentSetup struct {
	ID                   string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	MerchantOnboardingID string `gorm:"type:uuid;not null;index"`
	PayloadJSON          string `gorm:"type:jsonb;not null"`
	Status               string `gorm:"type:varchar(50);not null;default:DRAFT"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt `gorm:"index"`
}
