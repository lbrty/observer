package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	domainauth "github.com/lbrty/observer/internal/domain/auth"
	"github.com/lbrty/observer/internal/domain/user"
	ucauth "github.com/lbrty/observer/internal/usecase/auth"
)

// AuthHandler exposes auth HTTP endpoints.
type AuthHandler struct {
	registerUC *ucauth.RegisterUseCase
	loginUC    *ucauth.LoginUseCase
	refreshUC  *ucauth.RefreshTokenUseCase
	logoutUC   *ucauth.LogoutUseCase
}

// NewAuthHandler creates an AuthHandler.
func NewAuthHandler(
	registerUC *ucauth.RegisterUseCase,
	loginUC *ucauth.LoginUseCase,
	refreshUC *ucauth.RefreshTokenUseCase,
	logoutUC *ucauth.LogoutUseCase,
) *AuthHandler {
	return &AuthHandler{
		registerUC: registerUC,
		loginUC:    loginUC,
		refreshUC:  refreshUC,
		logoutUC:   logoutUC,
	}
}

// Register handles POST /auth/register.
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param input body ucauth.RegisterInput true "Registration payload"
// @Success 201 {object} ucauth.RegisterOutput
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input ucauth.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.registerUC.Execute(c.Request.Context(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, out)
}

// Login handles POST /auth/login.
// @Summary Log in with credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param input body ucauth.LoginInput true "Login payload"
// @Success 200 {object} ucauth.LoginOutput
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input ucauth.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.loginUC.Execute(c.Request.Context(), input, c.GetHeader("User-Agent"), c.ClientIP())
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, out)
}

// RefreshToken handles POST /auth/refresh.
// @Summary Refresh access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param input body ucauth.RefreshTokenInput true "Refresh token payload"
// @Success 200 {object} ucauth.TokenPair
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var input ucauth.RefreshTokenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.refreshUC.Execute(c.Request.Context(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Logout handles POST /auth/logout.
// @Summary Log out and invalidate refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body ucauth.LogoutInput true "Logout payload"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var input ucauth.LogoutInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.logoutUC.Execute(c.Request.Context(), input.RefreshToken); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func (h *AuthHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, user.ErrEmailExists), errors.Is(err, user.ErrPhoneExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, user.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case errors.Is(err, user.ErrUserNotActive):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, user.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, user.ErrInvalidRole):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, domainauth.ErrSessionNotFound), errors.Is(err, domainauth.ErrSessionExpired):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
