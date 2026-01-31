package auth

import (
	"gorm.io/gorm"

	domain "mob-backend/internal/domain/auth"
)

// emailVerificationRepository is a thin GORM-backed token repository.
type emailVerificationRepository struct {
	db *gorm.DB
}

// NewEmailVerificationRepository builds a GORM repository for verification tokens.
func NewEmailVerificationRepository(db *gorm.DB) *emailVerificationRepository {
	return &emailVerificationRepository{db: db}
}

// Create inserts a new verification token record.
func (r *emailVerificationRepository) Create(token domain.EmailVerificationToken) error {
	return r.db.Create(&token).Error
}

// Get loads a token by its value.
func (r *emailVerificationRepository) Get(token string) (*domain.EmailVerificationToken, error) {
	var stored domain.EmailVerificationToken
	err := r.db.First(&stored, "token = ?", token).Error
	if err != nil {
		return nil, err
	}
	return &stored, nil
}

// Update persists token changes (expiry/used flags).
func (r *emailVerificationRepository) Update(token domain.EmailVerificationToken) error {
	return r.db.Save(&token).Error
}
