package auth

import (
	"errors"
	"strings"
	"sync"

	domain "mob-backend/internal/domain/auth"
)

// ErrUserNotFound is returned when a user lookup fails.
var ErrUserNotFound = errors.New("user not found")

// MemoryUserRepository stores users in memory for development and testing.
type MemoryUserRepository struct {
	mu      sync.RWMutex
	byID    map[string]domain.User
	byEmail map[string]string
}

// NewMemoryUserRepository builds an empty in-memory repository.
func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		byID:    make(map[string]domain.User),
		byEmail: make(map[string]string),
	}
}

// Create inserts a new user and indexes by email.
func (r *MemoryUserRepository) Create(user domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	emailKey := strings.ToLower(user.Email)
	r.byID[user.ID] = user
	r.byEmail[emailKey] = user.ID
	return nil
}

// GetByEmail fetches a user by email address.
func (r *MemoryUserRepository) GetByEmail(email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.byEmail[strings.ToLower(email)]
	if !ok {
		return nil, ErrUserNotFound
	}
	user := r.byID[id]
	return &user, nil
}

// GetByID fetches a user by ID.
func (r *MemoryUserRepository) GetByID(id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.byID[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

// Update replaces an existing user record.
func (r *MemoryUserRepository) Update(user domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[user.ID]; !ok {
		return ErrUserNotFound
	}
	r.byID[user.ID] = user
	r.byEmail[strings.ToLower(user.Email)] = user.ID
	return nil
}
