package auth

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	domain "mob-backend/internal/domain/auth"
	dto "mob-backend/internal/dto"
	repo "mob-backend/internal/repository/auth"
	email "mob-backend/internal/service/email"
)

// ErrEmailAlreadyRegistered is returned when a user signs up with an existing email.
var ErrEmailAlreadyRegistered = errors.New("email already registered")

// ErrPasswordMismatch is returned when password and confirm_password differ.
var ErrPasswordMismatch = errors.New("password mismatch")

// ErrVerificationInvalid is returned when verification token is invalid or expired.
var ErrVerificationInvalid = errors.New("verification token invalid or expired")

// ErrEmailNotVerified is returned when login is attempted before verification.
var ErrEmailNotVerified = errors.New("email not verified")

// ErrInvalidCredentials is returned when login credentials are incorrect.
var ErrInvalidCredentials = errors.New("invalid credentials")

// AuthUsecase defines the auth flow business rules.
type AuthUsecase interface {
	Register(req dto.RegisterRequest) (*dto.RegisterResponseData, error)
	VerifyEmail(token string) error
	Login(req dto.LoginRequest) (*dto.LoginResponseData, error)
	GetMe(userID string) (*dto.UserResponse, error)
}

// TokenConfig stores JWT parameters to keep business logic testable.
type TokenConfig struct {
	Secret         string
	AccessTokenTTL time.Duration
	Issuer         string
}

// AuthService implements the auth flow using repositories and JWT.
type AuthService struct {
	users        repo.UserRepository
	tokens       repo.VerificationTokenRepository
	emailService email.Emailer
	config       TokenConfig
}

// NewAuthService wires repositories and JWT config for the auth usecase.
func NewAuthService(users repo.UserRepository, tokens repo.VerificationTokenRepository, emailService email.Emailer, config TokenConfig) *AuthService {
	return &AuthService{users: users, tokens: tokens, emailService: emailService, config: config}
}

// Register creates a new user and emits a verification token for email activation.
func (s *AuthService) Register(req dto.RegisterRequest) (*dto.RegisterResponseData, error) {
	if req.Password != req.ConfirmPassword {
		return nil, ErrPasswordMismatch
	}

	if _, err := s.users.GetByEmail(req.Email); err == nil {
		return nil, ErrEmailAlreadyRegistered
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	userID := uuid.NewString()
	now := time.Now()
	user := domain.User{
		ID:            userID,
		FullName:      req.FullName,
		BusinessName:  req.BusinessName,
		Email:         req.Email,
		PhoneNumber:   req.PhoneNumber,
		PasswordHash:  string(passwordHash),
		Status:        domain.StatusPendingEmailVerification,
		EmailVerified: false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.users.Create(user); err != nil {
		return nil, err
	}

	verificationToken := domain.EmailVerificationToken{
		Token:     uuid.NewString(),
		UserID:    userID,
		ExpiresAt: now.Add(24 * time.Hour),
		Used:      false,
		CreatedAt: now,
	}

	if err := s.tokens.Create(verificationToken); err != nil {
		return nil, err
	}

	if s.emailService != nil {
		if err := s.emailService.SendVerificationEmail(user.Email, verificationToken.Token); err != nil {
			log.Printf("failed to send verification email: %v", err)
		}
	}

	// Simulate email delivery for now to unblock frontend integration.
	log.Printf("verification link: /api/auth/verify-email?token=%s", verificationToken.Token)

	return &dto.RegisterResponseData{
		UserID: userID,
		Email:  req.Email,
		Status: string(domain.StatusPendingEmailVerification),
	}, nil
}

// VerifyEmail activates the user account if the verification token is valid.
func (s *AuthService) VerifyEmail(token string) error {
	stored, err := s.tokens.Get(token)
	if err != nil {
		return ErrVerificationInvalid
	}

	if stored.Used || time.Now().After(stored.ExpiresAt) {
		return ErrVerificationInvalid
	}

	user, err := s.users.GetByID(stored.UserID)
	if err != nil {
		return err
	}

	user.Status = domain.StatusActive
	user.EmailVerified = true
	user.UpdatedAt = time.Now()

	if err := s.users.Update(*user); err != nil {
		return err
	}

	stored.Used = true
	if err := s.tokens.Update(*stored); err != nil {
		return err
	}

	return nil
}

// Login validates credentials, checks verification status, and issues a JWT.
func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponseData, error) {
	user, err := s.users.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	if user.Status != domain.StatusActive || !user.EmailVerified {
		return nil, ErrEmailNotVerified
	}

	accessToken, expiresIn, err := s.generateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponseData{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
		User: dto.UserResponse{
			ID:            user.ID,
			FullName:      user.FullName,
			BusinessName:  user.BusinessName,
			Email:         user.Email,
			PhoneNumber:   user.PhoneNumber,
			EmailVerified: user.EmailVerified,
			Status:        string(user.Status),
		},
	}, nil
}

// GetMe returns the current user's profile data for dashboard preload.
func (s *AuthService) GetMe(userID string) (*dto.UserResponse, error) {
	user, err := s.users.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:            user.ID,
		FullName:      user.FullName,
		BusinessName:  user.BusinessName,
		Email:         user.Email,
		PhoneNumber:   user.PhoneNumber,
		EmailVerified: user.EmailVerified,
		Status:        string(user.Status),
	}, nil
}

// generateAccessToken issues an HS256 JWT with the configured TTL.
func (s *AuthService) generateAccessToken(userID string) (string, int, error) {
	expiresAt := time.Now().Add(s.config.AccessTokenTTL)

	claims := jwt.MapClaims{
		"sub": userID,
		"iss": s.config.Issuer,
		"exp": expiresAt.Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", 0, err
	}

	return signed, int(s.config.AccessTokenTTL.Seconds()), nil
}
