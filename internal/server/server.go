package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/app"
	"github.com/lbrty/observer/internal/config"
	"github.com/lbrty/observer/internal/database"
	"github.com/lbrty/observer/internal/handler"
	"github.com/lbrty/observer/internal/health"
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
	router := gin.New()

	s := &Server{router: router, cfg: &cfg.Server}
	s.setupMiddleware(log)
	s.setupRoutes(db, container)

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

func (s *Server) setupMiddleware(log *slog.Logger) {
	s.router.Use(requestIDMiddleware())
	s.router.Use(logger.GinMiddleware(log))
	s.router.Use(gin.Recovery())
}

func (s *Server) setupRoutes(db database.DB, container *app.Container) {
	healthHandler := health.NewHandler(db)
	s.router.GET("/health", healthHandler.Health)

	authMW := middleware.NewAuthMiddleware(container.TokenGenerator)
	_ = middleware.NewProjectAuthMiddleware(container.PermissionRepo)

	authHandler := handler.NewAuthHandler(
		container.RegisterUC,
		container.LoginUC,
		container.RefreshTokenUC,
		container.LogoutUC,
	)

	auth := s.router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authMW.Authenticate(), authHandler.Logout)
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
