package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/config"
	domainauth "github.com/lbrty/observer/internal/domain/auth"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/middleware"
	"github.com/lbrty/observer/internal/repository"
	ucauth "github.com/lbrty/observer/internal/usecase/auth"
)

const (
	accessTokenCookie  = "access_token"
	refreshTokenCookie = "refresh_token"
)

// AuthHandler exposes auth HTTP endpoints.
type AuthHandler struct {
	registerUC *ucauth.RegisterUseCase
	loginUC    *ucauth.LoginUseCase
	refreshUC  *ucauth.RefreshTokenUseCase
	logoutUC   *ucauth.LogoutUseCase
	userRepo   repository.UserRepository
	cookie     config.CookieConfig
	jwt        config.JWTConfig
}

// NewAuthHandler creates an AuthHandler.
func NewAuthHandler(
	registerUC *ucauth.RegisterUseCase,
	loginUC *ucauth.LoginUseCase,
	refreshUC *ucauth.RefreshTokenUseCase,
	logoutUC *ucauth.LogoutUseCase,
	userRepo repository.UserRepository,
	cookie config.CookieConfig,
	jwt config.JWTConfig,
) *AuthHandler {
	return &AuthHandler{
		registerUC: registerUC,
		loginUC:    loginUC,
		refreshUC:  refreshUC,
		logoutUC:   logoutUC,
		userRepo:   userRepo,
		cookie:     cookie,
		jwt:        jwt,
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

	if !out.RequiresMFA && out.Tokens != nil {
		h.setTokenCookies(c, out.Tokens.AccessToken, out.Tokens.RefreshToken)
	}

	c.JSON(http.StatusOK, out)
}

// Me handles GET /auth/me.
// @Summary Get current authenticated user
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ucauth.UserDTO
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userID, ok := middleware.UserIDFrom(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user identity"})
		return
	}

	u, err := h.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          u.ID.String(),
		"email":       u.Email,
		"phone":       u.Phone,
		"role":        string(u.Role),
		"is_verified": u.IsVerified,
		"created_at":  u.CreatedAt,
	})
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
	refreshToken := h.readRefreshToken(c)
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token is required"})
		return
	}

	tokens, err := h.refreshUC.Execute(c.Request.Context(), ucauth.RefreshTokenInput{
		RefreshToken: refreshToken,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.setTokenCookies(c, tokens.AccessToken, tokens.RefreshToken)
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
	refreshToken := h.readRefreshToken(c)
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token is required"})
		return
	}

	if err := h.logoutUC.Execute(c.Request.Context(), refreshToken); err != nil {
		h.handleError(c, err)
		return
	}

	h.clearTokenCookies(c)
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// readRefreshToken reads the refresh token from cookie, falling back to JSON body.
func (h *AuthHandler) readRefreshToken(c *gin.Context) string {
	if token, err := c.Cookie(refreshTokenCookie); err == nil && token != "" {
		return token
	}

	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&body); err == nil && body.RefreshToken != "" {
		return body.RefreshToken
	}

	return ""
}

func (h *AuthHandler) setTokenCookies(c *gin.Context, accessToken, refreshToken string) {
	sameSite := h.cookie.HTTPSameSite()

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     accessTokenCookie,
		Value:    accessToken,
		Path:     "/",
		Domain:   h.cookie.Domain,
		MaxAge:   int(h.cookie.MaxAge.Seconds()),
		HttpOnly: true,
		Secure:   h.cookie.Secure,
		SameSite: sameSite,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     refreshTokenCookie,
		Value:    refreshToken,
		Path:     "/auth",
		Domain:   h.cookie.Domain,
		MaxAge:   int(h.cookie.MaxAge.Seconds()),
		HttpOnly: true,
		Secure:   h.cookie.Secure,
		SameSite: sameSite,
	})
}

func (h *AuthHandler) clearTokenCookies(c *gin.Context) {
	sameSite := h.cookie.HTTPSameSite()

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     accessTokenCookie,
		Value:    "",
		Path:     "/",
		Domain:   h.cookie.Domain,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.cookie.Secure,
		SameSite: sameSite,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     refreshTokenCookie,
		Value:    "",
		Path:     "/auth",
		Domain:   h.cookie.Domain,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.cookie.Secure,
		SameSite: sameSite,
	})
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
