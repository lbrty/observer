package auth

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/config"
	"github.com/lbrty/observer/internal/handler"
	"github.com/lbrty/observer/internal/middleware"
	ucauth "github.com/lbrty/observer/internal/usecase/auth"
)

const (
	accessTokenCookie  = "access_token"
	refreshTokenCookie = "refresh_token"
)

// AuthHandler exposes auth HTTP endpoints.
type AuthHandler struct {
	authUC *ucauth.AuthUseCase
	cookie config.CookieConfig
	jwt    config.JWTConfig
}

// NewAuthHandler creates an AuthHandler.
func NewAuthHandler(
	authUC *ucauth.AuthUseCase,
	cookie config.CookieConfig,
	jwt config.JWTConfig,
) *AuthHandler {
	return &AuthHandler{
		authUC: authUC,
		cookie: cookie,
		jwt:    jwt,
	}
}

// Register handles POST /auth/register.
func (h *AuthHandler) Register(c *gin.Context) {
	input, ok := handler.BindJSON[ucauth.RegisterInput](c)
	if !ok {
		return
	}

	out, err := h.authUC.Register(c.Request.Context(), input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, out)
}

// Login handles POST /auth/login.
func (h *AuthHandler) Login(c *gin.Context) {
	input, ok := handler.BindJSON[ucauth.LoginInput](c)
	if !ok {
		return
	}

	out, err := h.authUC.Login(c.Request.Context(), input, c.GetHeader("User-Agent"), c.ClientIP())
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	if !out.RequiresMFA && out.Tokens != nil {
		h.setTokenCookies(c, out.Tokens.AccessToken, out.Tokens.RefreshToken)
	}

	c.JSON(http.StatusOK, out)
}

// Me handles GET /auth/me.
func (h *AuthHandler) Me(c *gin.Context) {
	userID, ok := middleware.UserIDFrom(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, handler.ErrJSON("errors.auth.missingUser", "missing user identity"))
		return
	}
	dto, err := h.authUC.Me(c.Request.Context(), userID)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto)
}

// RefreshToken handles POST /auth/refresh.
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken := h.readRefreshToken(c)
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, handler.ErrJSON("errors.auth.refreshTokenRequired", "refresh token is required"))
		return
	}

	tokens, err := h.authUC.RefreshToken(c.Request.Context(), ucauth.RefreshTokenInput{
		RefreshToken: refreshToken,
		UserAgent:    c.GetHeader("User-Agent"),
		IP:           c.ClientIP(),
	})
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	h.setTokenCookies(c, tokens.AccessToken, tokens.RefreshToken)
	c.JSON(http.StatusOK, tokens)
}

// Logout handles POST /auth/logout.
func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken := h.readRefreshToken(c)
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, handler.ErrJSON("errors.auth.refreshTokenRequired", "refresh token is required"))
		return
	}

	if err := h.authUC.Logout(c.Request.Context(), refreshToken); err != nil {
		handler.HandleError(c, err)
		return
	}

	h.clearTokenCookies(c)
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// UpdateProfile handles PATCH /auth/me.
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, ok := middleware.UserIDFrom(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, handler.ErrJSON("errors.auth.missingUser", "missing user identity"))
		return
	}

	input, ok := handler.BindJSON[ucauth.UpdateProfileInput](c)
	if !ok {
		return
	}

	dto, err := h.authUC.UpdateProfile(c.Request.Context(), userID, input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto)
}

// ChangePassword handles POST /auth/change-password.
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, ok := middleware.UserIDFrom(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, handler.ErrJSON("errors.auth.missingUser", "missing user identity"))
		return
	}

	input, ok := handler.BindJSON[ucauth.ChangePasswordInput](c)
	if !ok {
		return
	}

	if err := h.authUC.ChangePassword(c.Request.Context(), userID, input); err != nil {
		handler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}

// VerifyMFA handles POST /auth/mfa.
func (h *AuthHandler) VerifyMFA(c *gin.Context) {
	var input ucauth.VerifyMFAInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrJSON("errors.auth.invalidRequest", err.Error()))
		return
	}
	out, err := h.authUC.VerifyMFA(c.Request.Context(), input, c.GetHeader("User-Agent"), c.ClientIP())
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	if out.Tokens != nil {
		h.setTokenCookies(c, out.Tokens.AccessToken, out.Tokens.RefreshToken)
	}
	c.JSON(http.StatusOK, out)
}

// MFASetup handles GET /auth/mfa/setup.
func (h *AuthHandler) MFASetup(c *gin.Context) {
	userID, ok := middleware.UserIDFrom(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, handler.ErrJSON("errors.auth.missingUser", "missing user identity"))
		return
	}
	out, err := h.authUC.SetupMFA(c.Request.Context(), userID)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// EnableMFA handles POST /auth/mfa/enable.
func (h *AuthHandler) EnableMFA(c *gin.Context) {
	userID, ok := middleware.UserIDFrom(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, handler.ErrJSON("errors.auth.missingUser", "missing user identity"))
		return
	}
	var input ucauth.EnableMFAInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrJSON("errors.auth.invalidRequest", err.Error()))
		return
	}
	out, err := h.authUC.EnableMFA(c.Request.Context(), userID, input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// DisableMFA handles POST /auth/mfa/disable.
func (h *AuthHandler) DisableMFA(c *gin.Context) {
	userID, ok := middleware.UserIDFrom(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, handler.ErrJSON("errors.auth.missingUser", "missing user identity"))
		return
	}
	var input ucauth.DisableMFAInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrJSON("errors.auth.invalidRequest", err.Error()))
		return
	}
	if err := h.authUC.DisableMFA(c.Request.Context(), userID, input); err != nil {
		handler.HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
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
		Path:     "/api/auth",
		Domain:   h.cookie.Domain,
		MaxAge:   int(h.cookie.MaxAge.Seconds()),
		HttpOnly: true,
		Secure:   h.cookie.Secure,
		SameSite: sameSite,
	})

	// CSRF double-submit cookie — HttpOnly: false so JS can read and echo it.
	csrfToken := generateCSRFToken()
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     middleware.CSRFTokenCookie,
		Value:    csrfToken,
		Path:     "/",
		Domain:   h.cookie.Domain,
		MaxAge:   int(h.cookie.MaxAge.Seconds()),
		HttpOnly: false,
		Secure:   h.cookie.Secure,
		SameSite: sameSite,
	})
}

func generateCSRFToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b) // crypto/rand.Read never fails; panics internally on catastrophic entropy failure
	return hex.EncodeToString(b)
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
		Path:     "/api/auth",
		Domain:   h.cookie.Domain,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.cookie.Secure,
		SameSite: sameSite,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     middleware.CSRFTokenCookie,
		Value:    "",
		Path:     "/",
		Domain:   h.cookie.Domain,
		MaxAge:   -1,
		HttpOnly: false,
		Secure:   h.cookie.Secure,
		SameSite: sameSite,
	})
}
