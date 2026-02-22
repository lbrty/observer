package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"

	"github.com/lbrty/observer/internal/crypto"
	"github.com/lbrty/observer/internal/domain/user"
)

// ctxKey is a typed key for Gin context values.
type ctxKey string

const (
	CtxUserID           ctxKey = "user_id"
	CtxUserRole         ctxKey = "user_role"
	CtxProjectID        ctxKey = "project_id"
	CtxProjectRole      ctxKey = "project_role"
	CtxCanViewContact   ctxKey = "can_view_contact"
	CtxCanViewPersonal  ctxKey = "can_view_personal"
	CtxCanViewDocuments ctxKey = "can_view_documents"
)

// AuthMiddleware provides JWT-based authentication handlers.
type AuthMiddleware struct {
	tokenGen crypto.TokenGenerator
}

// NewAuthMiddleware creates an AuthMiddleware.
func NewAuthMiddleware(tokenGen crypto.TokenGenerator) *AuthMiddleware {
	return &AuthMiddleware{tokenGen: tokenGen}
}

// Authenticate validates the Bearer token and sets user_id / user_role in context.
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			c.Abort()
			return
		}

		claims, err := m.tokenGen.ValidateAccessToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		userID, err := ulid.Parse(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID in token"})
			c.Abort()
			return
		}

		c.Set(string(CtxUserID), userID)
		c.Set(string(CtxUserRole), claims.Role)
		c.Next()
	}
}

// RequireRole checks that the authenticated user has one of the allowed platform roles.
func (m *AuthMiddleware) RequireRole(roles ...user.Role) gin.HandlerFunc {
	allowed := make(map[user.Role]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *gin.Context) {
		roleVal, exists := c.Get(string(CtxUserRole))
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "role not found"})
			c.Abort()
			return
		}
		role := user.Role(roleVal.(string))
		if _, ok := allowed[role]; !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// UserIDFrom extracts the authenticated user's ULID from the Gin context.
func UserIDFrom(c *gin.Context) (ulid.ULID, bool) {
	val, exists := c.Get(string(CtxUserID))
	if !exists {
		return ulid.ULID{}, false
	}
	id, ok := val.(ulid.ULID)
	return id, ok
}

// UserRoleFrom extracts the authenticated user's platform role from the Gin context.
func UserRoleFrom(c *gin.Context) (user.Role, bool) {
	val, exists := c.Get(string(CtxUserRole))
	if !exists {
		return "", false
	}
	role, ok := val.(string)
	return user.Role(role), ok
}
