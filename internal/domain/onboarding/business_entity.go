package onboarding

import (
	"time"

	"gorm.io/gorm"
)

// BusinessEntityMerchant stores merchant-level business data in draft state.
// Kept separate to avoid one giant table and to match per-step saving.
type BusinessEntityMerchant struct {
	// ID is the primary key for this draft record.
	ID string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	// MerchantOnboardingID links to the root onboarding flow.
	MerchantOnboardingID string `gorm:"type:uuid;not null;index"`
	// Payload keeps flexible fields from the contract (address, category, etc).
	PayloadJSON string `gorm:"type:jsonb;not null"`
	// Status remains DRAFT and must not change onboarding status.
	Status string `gorm:"type:varchar(50);not null;default:DRAFT"`
	// CreatedAt and UpdatedAt provide audit timestamps.
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BusinessEntityOwner stores owner data in draft state.
type BusinessEntityOwner struct {
	ID                   string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	MerchantOnboardingID string `gorm:"type:uuid;not null;index"`
	PayloadJSON          string `gorm:"type:jsonb;not null"`
	Status               string `gorm:"type:varchar(50);not null;default:DRAFT"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt `gorm:"index"`
}

// BusinessEntityPIC stores PIC admin data in draft state.
type BusinessEntityPIC struct {
	ID                   string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	MerchantOnboardingID string `gorm:"type:uuid;not null;index"`
	PayloadJSON          string `gorm:"type:jsonb;not null"`
	Status               string `gorm:"type:varchar(50);not null;default:DRAFT"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt `gorm:"index"`
}

// BusinessEntitySettlement stores settlement bank data in draft state.
type BusinessEntitySettlement struct {
	ID                   string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	MerchantOnboardingID string `gorm:"type:uuid;not null;index"`
	PayloadJSON          string `gorm:"type:jsonb;not null"`
	Status               string `gorm:"type:varchar(50);not null;default:DRAFT"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt `gorm:"index"`
}
