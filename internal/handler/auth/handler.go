package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dto "mob-backend/internal/dto"
	repo "mob-backend/internal/repository/auth"
	usecase "mob-backend/internal/usecase/auth"
)

// Handler orchestrates HTTP request parsing and delegates to the auth usecase.
type Handler struct {
	uc usecase.AuthUsecase
}

// NewHandler builds an auth handler with the provided usecase.
func NewHandler(uc usecase.AuthUsecase) *Handler {
	return &Handler{uc: uc}
}

// Register handles POST /api/auth/register with validation and response mapping.
// @Summary Register merchant user
// @Description Create user account and emit email verification token (status PENDING_EMAIL_VERIFICATION)
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body dto.RegisterRequest true "Register request payload"
// @Success 201 {object} dto.APIResponse{data=dto.RegisterResponseData}
// @Failure 400 {object} dto.APIResponse
// @Failure 409 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid request payload",
		})
		return
	}

	data, err := h.uc.Register(req)
	if err != nil {
		switch err {
		case usecase.ErrPasswordMismatch:
			c.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: "Password and confirm_password do not match"})
			return
		case usecase.ErrEmailAlreadyRegistered:
			c.JSON(http.StatusConflict, dto.APIResponse{Success: false, Message: "Email already registered"})
			return
		default:
			c.JSON(http.StatusInternalServerError, dto.APIResponse{Success: false, Message: "Internal server error"})
			return
		}
	}

	c.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "Registration successful. Please verify your email.",
		Data:    data,
	})
}

// VerifyEmail handles GET /api/auth/verify-email?token=... to activate accounts.
// @Summary Verify email link
// @Description Activate account by validating one-time token
// @Tags Auth
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Router /api/auth/verify-email [get]
func (h *Handler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: "Invalid or expired verification link"})
		return
	}

	if err := h.uc.VerifyEmail(token); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: "Invalid or expired verification link"})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{Success: true, Message: "Email verified successfully"})
}

// Login handles POST /api/auth/login with credential checks and JWT issuance.
// @Summary User login
// @Description Authenticate verified user and issue JWT access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body dto.LoginRequest true "Login payload"
// @Success 200 {object} dto.APIResponse{data=dto.LoginResponseData}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 403 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: "Invalid request payload"})
		return
	}

	data, err := h.uc.Login(req)
	if err != nil {
		switch err {
		case usecase.ErrEmailNotVerified:
			c.JSON(http.StatusForbidden, dto.APIResponse{Success: false, Message: "Please verify your email before signing in", Code: "EMAIL_NOT_VERIFIED"})
			return
		case usecase.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, dto.APIResponse{Success: false, Message: "Invalid email or password", Code: "INVALID_CREDENTIALS"})
			return
		case repo.ErrUserNotFound:
			c.JSON(http.StatusNotFound, dto.APIResponse{Success: false, Message: "User not found", Code: "USER_NOT_FOUND"})
			return
		default:
			c.JSON(http.StatusInternalServerError, dto.APIResponse{Success: false, Message: "Internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, dto.APIResponse{Success: true, Message: "Login successful", Data: data})
}

// Me handles GET /api/auth/me using the authenticated user ID.
// @Summary Get current user
// @Description Return logged-in merchant profile
// @Tags Auth
// @Produce json
// @Success 200 {object} dto.APIResponse{data=dto.UserResponse}
// @Failure 401 {object} dto.APIResponse
// @Security BearerAuth
// @Router /api/auth/me [get]
func (h *Handler) Me(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.APIResponse{Success: false, Message: "Unauthorized. Invalid or expired token", Code: "UNAUTHORIZED"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		c.JSON(http.StatusUnauthorized, dto.APIResponse{Success: false, Message: "Unauthorized. Invalid or expired token", Code: "UNAUTHORIZED"})
		return
	}

	data, err := h.uc.GetMe(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.APIResponse{Success: false, Message: "Unauthorized. Invalid or expired token", Code: "UNAUTHORIZED"})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{Success: true, Message: "Current user retrieved", Data: data})
}
