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
	authUC   *ucauth.AuthUseCase
	userRepo repository.UserRepository
	cookie   config.CookieConfig
	jwt      config.JWTConfig
}

// NewAuthHandler creates an AuthHandler.
func NewAuthHandler(
	authUC *ucauth.AuthUseCase,
	userRepo repository.UserRepository,
	cookie config.CookieConfig,
	jwt config.JWTConfig,
) *AuthHandler {
	return &AuthHandler{
		authUC:   authUC,
		userRepo: userRepo,
		cookie:   cookie,
		jwt:      jwt,
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
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}

	out, err := h.authUC.Register(c.Request.Context(), input)
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
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}

	out, err := h.authUC.Login(c.Request.Context(), input, c.GetHeader("User-Agent"), c.ClientIP())
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
		c.JSON(http.StatusUnauthorized, errJSON("errors.auth.missingUser", "missing user identity"))
		return
	}

	u, err := h.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	resp := gin.H{
		"id":          u.ID.String(),
		"first_name":  u.FirstName,
		"last_name":   u.LastName,
		"email":       u.Email,
		"phone":       u.Phone,
		"role":        string(u.Role),
		"is_verified": u.IsVerified,
		"created_at":  u.CreatedAt,
	}
	if u.OfficeID != nil {
		resp["office_id"] = *u.OfficeID
	}
	c.JSON(http.StatusOK, resp)
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
		c.JSON(http.StatusBadRequest, errJSON("errors.auth.refreshTokenRequired", "refresh token is required"))
		return
	}

	tokens, err := h.authUC.RefreshToken(c.Request.Context(), ucauth.RefreshTokenInput{
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
		c.JSON(http.StatusBadRequest, errJSON("errors.auth.refreshTokenRequired", "refresh token is required"))
		return
	}

	if err := h.authUC.Logout(c.Request.Context(), refreshToken); err != nil {
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

// UpdateProfile handles PATCH /auth/me.
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, ok := middleware.UserIDFrom(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errJSON("errors.auth.missingUser", "missing user identity"))
		return
	}

	var input ucauth.UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}

	dto, err := h.authUC.UpdateProfile(c.Request.Context(), userID, input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto)
}

// ChangePassword handles POST /auth/change-password.
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, ok := middleware.UserIDFrom(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errJSON("errors.auth.missingUser", "missing user identity"))
		return
	}

	var input ucauth.ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}

	if err := h.authUC.ChangePassword(c.Request.Context(), userID, input); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}

func (h *AuthHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, user.ErrEmailExists):
		c.JSON(http.StatusConflict, errJSON("errors.user.emailExists", err.Error()))
	case errors.Is(err, user.ErrPhoneExists):
		c.JSON(http.StatusConflict, errJSON("errors.user.phoneExists", err.Error()))
	case errors.Is(err, user.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, errJSON("errors.auth.invalidCredentials", err.Error()))
	case errors.Is(err, user.ErrUserNotActive):
		c.JSON(http.StatusForbidden, errJSON("errors.user.notActive", err.Error()))
	case errors.Is(err, user.ErrUserNotFound):
		c.JSON(http.StatusNotFound, errJSON("errors.user.notFound", err.Error()))
	case errors.Is(err, user.ErrInvalidRole):
		c.JSON(http.StatusBadRequest, errJSON("errors.user.invalidRole", err.Error()))
	case errors.Is(err, domainauth.ErrSessionNotFound):
		c.JSON(http.StatusUnauthorized, errJSON("errors.auth.sessionNotFound", err.Error()))
	case errors.Is(err, domainauth.ErrSessionExpired):
		c.JSON(http.StatusUnauthorized, errJSON("errors.auth.sessionExpired", err.Error()))
	default:
		internalError(c, "handle auth operation", err)
	}
}
