package auth

import (
	"errors"
	"sync"

	domain "mob-backend/internal/domain/auth"
)

// ErrTokenNotFound is returned when a verification token is missing.
var ErrTokenNotFound = errors.New("verification token not found")

// MemoryVerificationTokenRepository stores verification tokens in memory.
type MemoryVerificationTokenRepository struct {
	mu      sync.RWMutex
	byToken map[string]domain.EmailVerificationToken
}

// NewMemoryVerificationTokenRepository builds an empty in-memory token repository.
func NewMemoryVerificationTokenRepository() *MemoryVerificationTokenRepository {
	return &MemoryVerificationTokenRepository{
		byToken: make(map[string]domain.EmailVerificationToken),
	}
}

// Create inserts a new verification token.
func (r *MemoryVerificationTokenRepository) Create(token domain.EmailVerificationToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byToken[token.Token] = token
	return nil
}

// Get fetches a token by its value.
func (r *MemoryVerificationTokenRepository) Get(token string) (*domain.EmailVerificationToken, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stored, ok := r.byToken[token]
	if !ok {
		return nil, ErrTokenNotFound
	}
	return &stored, nil
}

// Update replaces the stored token (for one-time use toggling).
func (r *MemoryVerificationTokenRepository) Update(token domain.EmailVerificationToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byToken[token.Token]; !ok {
		return ErrTokenNotFound
	}
	r.byToken[token.Token] = token
	return nil
}
