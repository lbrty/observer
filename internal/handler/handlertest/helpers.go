package handlertest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	validator "github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"

	"github.com/lbrty/observer/internal/config"
	"github.com/lbrty/observer/internal/middleware"
)

func init() {
	gin.SetMode(gin.TestMode)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("strongpassword", func(fl validator.FieldLevel) bool { //nolint:errcheck
			p := fl.Field().String()
			hasDigit := strings.ContainsAny(p, "0123456789")
			hasSpecial := strings.ContainsAny(p, "!@#$%^&*()-_=+[]{}|;:',.<>?/`~")
			return hasDigit && hasSpecial
		})
	}
}

// NewTestContext creates a gin context and response recorder for tests.
func NewTestContext(method, path string, body any) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	var req *http.Request
	if body != nil {
		b, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c, w
}

// NewTestContextWithParams creates a gin context with URL params.
func NewTestContextWithParams(method, path string, body any, params gin.Params) (*gin.Context, *httptest.ResponseRecorder) {
	c, w := NewTestContext(method, path, body)
	c.Params = params
	return c, w
}

// SetAuthContext sets the authenticated user ID in the gin context.
func SetAuthContext(c *gin.Context, userID ulid.ULID) {
	c.Set(string(middleware.CtxUserID), userID)
}

// TestID returns a new random ULID.
func TestID() ulid.ULID {
	return ulid.Make()
}

// ParseResponse unmarshals the response body into T.
func ParseResponse[T any](w *httptest.ResponseRecorder) T {
	var result T
	_ = json.Unmarshal(w.Body.Bytes(), &result)
	return result
}

// TestCookieConfig returns a test cookie configuration.
func TestCookieConfig() config.CookieConfig {
	return config.CookieConfig{
		Domain:   "localhost",
		Secure:   false,
		SameSite: "lax",
		MaxAge:   2 * time.Hour,
	}
}

// TestJWTConfig returns a test JWT configuration.
func TestJWTConfig() config.JWTConfig {
	return config.JWTConfig{
		AccessTTL:  15 * time.Minute,
		RefreshTTL: 168 * time.Hour,
		Issuer:     "observer-test",
	}
}
