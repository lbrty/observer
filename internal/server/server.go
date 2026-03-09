package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	validator "github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/lbrty/observer/api/swagger"
	"github.com/lbrty/observer/internal/app"
	"github.com/lbrty/observer/internal/config"
	"github.com/lbrty/observer/internal/database"
	"github.com/lbrty/observer/internal/logger"
	"github.com/lbrty/observer/internal/middleware"
	"github.com/lbrty/observer/internal/ulid"
)

// Server wraps the Gin engine and HTTP server.
type Server struct {
	router *gin.Engine
	srv    *http.Server
	cfg    *config.ServerConfig
}

// New creates and configures a new Server.
func New(cfg *config.Config, db database.DB, log *slog.Logger, container *app.Container) *Server {
	gin.SetMode(gin.ReleaseMode)
	registerCustomValidators()
	router := gin.New()

	s := &Server{router: router, cfg: &cfg.Server}
	s.setupMiddleware(cfg, log)
	s.setupRoutes(cfg, db, container)

	if cfg.Swagger.Enabled {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	s.srv = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return s
}

// Router returns the underlying Gin engine (useful for testing).
func (s *Server) Router() *gin.Engine {
	return s.router
}

// Start begins listening for HTTP requests.
func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) setupMiddleware(cfg *config.Config, log *slog.Logger) {
	// Limit non-multipart request bodies to 1 MB to prevent memory exhaustion.
	// Multipart (file uploads) are excluded — document_handler enforces 50 MB itself.
	s.router.Use(func(c *gin.Context) {
		if !strings.HasPrefix(c.GetHeader("Content-Type"), "multipart/form-data") {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
		}
		c.Next()
	})
	s.router.Use(requestIDMiddleware())
	s.router.Use(logger.GinMiddleware(log))
	if cfg.Sentry.Enabled() {
		s.router.Use(sentrygin.New(sentrygin.Options{Repanic: true}))
	}
	s.router.Use(gin.Recovery())
	s.router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.Origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-CSRF-Token"},
		ExposeHeaders:    []string{"X-Request-ID"},
		AllowCredentials: true,
	}))
	s.router.Use(middleware.SecurityHeaders())
	s.router.Use(middleware.CSRFProtection())
}

func registerCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("strongpassword", func(fl validator.FieldLevel) bool { //nolint:errcheck
			p := fl.Field().String()
			hasDigit := strings.ContainsAny(p, "0123456789")
			hasSpecial := strings.ContainsAny(p, "!@#$%^&*()-_=+[]{}|;:',.<>?/`~")
			return hasDigit && hasSpecial
		})
	}
}

func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := ulid.NewString()
		c.Request.Header.Set("X-Request-ID", id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}
