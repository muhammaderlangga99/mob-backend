package file

import (
	"time"

	"gorm.io/gorm"
)

// UploadedFile stores metadata for files uploaded during onboarding.
// Linked to a merchant onboarding flow to support document mapping.
type UploadedFile struct {
	ID                   string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	MerchantOnboardingID string `gorm:"type:uuid;not null;index"`
	// Category matches contract values (e.g., owner_ktp, npwp).
	Category string `gorm:"type:varchar(80);not null"`
	// FileURL points to object storage or CDN location.
	FileURL string `gorm:"type:text;not null"`
	// MimeType is returned to clients for preview/validation.
	MimeType string `gorm:"type:varchar(120);not null"`
	// CreatedAt provides upload timestamp for audit.
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
