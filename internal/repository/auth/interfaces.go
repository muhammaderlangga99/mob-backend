package auth

import "mob-backend/internal/domain/auth"

// UserRepository defines persistence operations for auth users.
type UserRepository interface {
	Create(user auth.User) error
	GetByEmail(email string) (*auth.User, error)
	GetByID(id string) (*auth.User, error)
	Update(user auth.User) error
}

// VerificationTokenRepository defines storage for email verification tokens.
type VerificationTokenRepository interface {
	Create(token auth.EmailVerificationToken) error
	Get(token string) (*auth.EmailVerificationToken, error)
	Update(token auth.EmailVerificationToken) error
}
