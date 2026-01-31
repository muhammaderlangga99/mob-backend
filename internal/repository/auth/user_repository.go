package auth

import (
	"strings"

	"gorm.io/gorm"

	domain "mob-backend/internal/domain/auth"
)

// userRepository is a thin GORM-backed implementation of UserRepository.
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository builds a GORM repository for users.
func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

// Create inserts a new user record.
func (r *userRepository) Create(user domain.User) error {
	return r.db.Create(&user).Error
}

// GetByEmail loads a user by email.
func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByID loads a user by primary key.
func (r *userRepository) GetByID(id string) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update persists changes to an existing user.
func (r *userRepository) Update(user domain.User) error {
	return r.db.Save(&user).Error
}
