package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/middleware"
	ucauth "github.com/lbrty/observer/internal/usecase/auth"
)

// VerifyMFA handles POST /auth/mfa — exchanges an MFA token + TOTP code for real tokens.
func (h *AuthHandler) VerifyMFA(c *gin.Context) {
	var input ucauth.VerifyMFAInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.auth.invalidRequest", err.Error()))
		return
	}
	out, err := h.authUC.VerifyMFA(c.Request.Context(), input, c.GetHeader("User-Agent"), c.ClientIP())
	if err != nil {
		HandleError(c, err)
		return
	}
	if out.Tokens != nil {
		h.setTokenCookies(c, out.Tokens.AccessToken, out.Tokens.RefreshToken)
	}
	c.JSON(http.StatusOK, out)
}

// MFASetup handles GET /auth/mfa/setup — generates a TOTP secret + QR URI.
func (h *AuthHandler) MFASetup(c *gin.Context) {
	userID, ok := middleware.UserIDFrom(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errJSON("errors.auth.missingUser", "missing user identity"))
		return
	}
	out, err := h.authUC.SetupMFA(c.Request.Context(), userID)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// EnableMFA handles POST /auth/mfa/enable — verifies TOTP code and enables MFA.
func (h *AuthHandler) EnableMFA(c *gin.Context) {
	userID, ok := middleware.UserIDFrom(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errJSON("errors.auth.missingUser", "missing user identity"))
		return
	}
	var input ucauth.EnableMFAInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.auth.invalidRequest", err.Error()))
		return
	}
	out, err := h.authUC.EnableMFA(c.Request.Context(), userID, input)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// DisableMFA handles POST /auth/mfa/disable — verifies TOTP code and disables MFA.
func (h *AuthHandler) DisableMFA(c *gin.Context) {
	userID, ok := middleware.UserIDFrom(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errJSON("errors.auth.missingUser", "missing user identity"))
		return
	}
	var input ucauth.DisableMFAInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.auth.invalidRequest", err.Error()))
		return
	}
	if err := h.authUC.DisableMFA(c.Request.Context(), userID, input); err != nil {
		HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
