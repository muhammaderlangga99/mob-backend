package dto

// RegisterRequest represents the payload for user registration.
type RegisterRequest struct {
	FullName        string `json:"full_name" binding:"required"`
	BusinessName    string `json:"business_name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	PhoneNumber     string `json:"phone_number" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

// RegisterResponseData is the data section returned after successful registration.
type RegisterResponseData struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

// VerifyEmailResponse is used for the verification endpoint result.
type VerifyEmailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// LoginRequest represents the payload for user login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponseData is the data section returned after successful login.
type LoginResponseData struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresIn   int          `json:"expires_in"`
	User        UserResponse `json:"user"`
}

// UserResponse represents the user profile payload exposed to clients.
type UserResponse struct {
	ID            string `json:"id"`
	FullName      string `json:"full_name"`
	BusinessName  string `json:"business_name"`
	Email         string `json:"email"`
	PhoneNumber   string `json:"phone_number"`
	EmailVerified bool   `json:"email_verified"`
	Status        string `json:"status"`
}

// MeResponse wraps the user profile for the /me endpoint.
type MeResponse struct {
	Success bool         `json:"success"`
	Data    UserResponse `json:"data"`
}

// APIResponse standardizes the JSON envelope for both success and error cases.
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Code    string      `json:"code,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
