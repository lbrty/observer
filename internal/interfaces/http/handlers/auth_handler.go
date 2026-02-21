package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	appauth "github.com/lbrty/observer/internal/application/auth"
	domainauth "github.com/lbrty/observer/internal/domain/auth"
	"github.com/lbrty/observer/internal/domain/user"
)

// AuthHandler exposes auth HTTP endpoints.
type AuthHandler struct {
	registerUC *appauth.RegisterUseCase
	loginUC    *appauth.LoginUseCase
	refreshUC  *appauth.RefreshTokenUseCase
	logoutUC   *appauth.LogoutUseCase
}

// NewAuthHandler creates an AuthHandler.
func NewAuthHandler(
	registerUC *appauth.RegisterUseCase,
	loginUC *appauth.LoginUseCase,
	refreshUC *appauth.RefreshTokenUseCase,
	logoutUC *appauth.LogoutUseCase,
) *AuthHandler {
	return &AuthHandler{
		registerUC: registerUC,
		loginUC:    loginUC,
		refreshUC:  refreshUC,
		logoutUC:   logoutUC,
	}
}

// Register handles POST /auth/register.
func (h *AuthHandler) Register(c *gin.Context) {
	var input appauth.RegisterInput
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
func (h *AuthHandler) Login(c *gin.Context) {
	var input appauth.LoginInput
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
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var input appauth.RefreshTokenInput
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
func (h *AuthHandler) Logout(c *gin.Context) {
	var input appauth.LogoutInput
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
