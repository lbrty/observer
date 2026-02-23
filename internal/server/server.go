package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/lbrty/observer/api/swagger"
	"github.com/lbrty/observer/internal/app"
	"github.com/lbrty/observer/internal/config"
	"github.com/lbrty/observer/internal/database"
	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/domain/user"
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

func (s *Server) setupMiddleware(log *slog.Logger) {
	s.router.Use(requestIDMiddleware())
	s.router.Use(logger.GinMiddleware(log))
	s.router.Use(gin.Recovery())
}

func (s *Server) setupRoutes(db database.DB, container *app.Container) {
	healthHandler := health.NewHandler(db)
	s.router.GET("/health", healthHandler.Health)

	authMW := middleware.NewAuthMiddleware(container.TokenGenerator)
	projectAuthMW := middleware.NewProjectAuthMiddleware(container.PermissionRepo)

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

	// Admin endpoints — requires authentication + admin role
	adminHandler := handler.NewAdminHandler(container.ListUsersUC, container.GetUserUC, container.UpdateUserUC)
	permHandler := handler.NewPermissionHandler(container.ListPermsUC, container.AssignPermUC, container.UpdatePermUC, container.RevokePermUC)
	countryHandler := handler.NewCountryHandler(container.CountryUC)
	stateHandler := handler.NewStateHandler(container.StateUC)
	placeHandler := handler.NewPlaceHandler(container.PlaceUC)
	officeHandler := handler.NewOfficeHandler(container.OfficeUC)
	categoryHandler := handler.NewCategoryHandler(container.CategoryUC)
	projectHandler := handler.NewProjectHandler(container.ProjectUC)

	admin := s.router.Group("/admin", authMW.Authenticate(), authMW.RequireRole(user.RoleAdmin))
	{
		admin.GET("/users", adminHandler.ListUsers)
		admin.GET("/users/:id", adminHandler.GetUser)
		admin.PATCH("/users/:id", adminHandler.UpdateUser)

		admin.GET("/projects", projectHandler.List)
		admin.GET("/projects/:project_id", projectHandler.Get)
		admin.POST("/projects", projectHandler.Create)
		admin.PATCH("/projects/:project_id", projectHandler.Update)

		admin.GET("/projects/:project_id/permissions", permHandler.ListPermissions)
		admin.POST("/projects/:project_id/permissions", permHandler.AssignPermission)
		admin.PATCH("/projects/:project_id/permissions/:id", permHandler.UpdatePermission)
		admin.DELETE("/projects/:project_id/permissions/:id", permHandler.RevokePermission)

		admin.GET("/countries", countryHandler.List)
		admin.GET("/countries/:id", countryHandler.Get)
		admin.POST("/countries", countryHandler.Create)
		admin.PATCH("/countries/:id", countryHandler.Update)
		admin.DELETE("/countries/:id", countryHandler.Delete)

		admin.GET("/states", stateHandler.List)
		admin.GET("/states/:id", stateHandler.Get)
		admin.POST("/states", stateHandler.Create)
		admin.PATCH("/states/:id", stateHandler.Update)
		admin.DELETE("/states/:id", stateHandler.Delete)

		admin.GET("/places", placeHandler.List)
		admin.GET("/places/:id", placeHandler.Get)
		admin.POST("/places", placeHandler.Create)
		admin.PATCH("/places/:id", placeHandler.Update)
		admin.DELETE("/places/:id", placeHandler.Delete)

		admin.GET("/offices", officeHandler.List)
		admin.GET("/offices/:id", officeHandler.Get)
		admin.POST("/offices", officeHandler.Create)
		admin.PATCH("/offices/:id", officeHandler.Update)
		admin.DELETE("/offices/:id", officeHandler.Delete)

		admin.GET("/categories", categoryHandler.List)
		admin.GET("/categories/:id", categoryHandler.Get)
		admin.POST("/categories", categoryHandler.Create)
		admin.PATCH("/categories/:id", categoryHandler.Update)
		admin.DELETE("/categories/:id", categoryHandler.Delete)
	}

	// Project-scoped endpoints — requires authentication + project role
	tagHandler := handler.NewTagHandler(container.TagUC)
	personHandler := handler.NewPersonHandler(container.PersonUC, container.PersonCategoryUC, container.PersonTagUC)
	supportHandler := handler.NewSupportRecordHandler(container.SupportRecordUC)
	migrationHandler := handler.NewMigrationRecordHandler(container.MigrationRecordUC)
	householdHandler := handler.NewHouseholdHandler(container.HouseholdUC)
	noteHandler := handler.NewNoteHandler(container.NoteUC)
	documentHandler := handler.NewDocumentHandler(container.DocumentUC)
	petHandler := handler.NewPetHandler(container.PetUC)

	proj := s.router.Group("/projects/:project_id", authMW.Authenticate())
	{
		// Read-level access
		read := proj.Group("", projectAuthMW.RequireProjectRole(project.ActionRead))
		{
			read.GET("/tags", tagHandler.List)
			read.GET("/people", personHandler.List)
			read.GET("/people/:person_id", personHandler.Get)
			read.GET("/people/:person_id/categories", personHandler.ListCategories)
			read.GET("/people/:person_id/tags", personHandler.ListTags)
			read.GET("/people/:person_id/migration-records", migrationHandler.List)
			read.GET("/people/:person_id/migration-records/:id", migrationHandler.Get)
			read.GET("/people/:person_id/notes", noteHandler.List)
			read.GET("/people/:person_id/documents", documentHandler.List)
			read.GET("/support-records", supportHandler.List)
			read.GET("/support-records/:id", supportHandler.Get)
			read.GET("/households", householdHandler.List)
			read.GET("/households/:id", householdHandler.Get)
			read.GET("/documents/:id", documentHandler.Get)
			read.GET("/pets", petHandler.List)
			read.GET("/pets/:id", petHandler.Get)
		}

		// Create-level access
		create := proj.Group("", projectAuthMW.RequireProjectRole(project.ActionCreate))
		{
			create.POST("/tags", tagHandler.Create)
			create.POST("/people", personHandler.Create)
			create.PUT("/people/:person_id/categories", personHandler.ReplaceCategories)
			create.PUT("/people/:person_id/tags", personHandler.ReplaceTags)
			create.POST("/people/:person_id/migration-records", migrationHandler.Create)
			create.POST("/people/:person_id/notes", noteHandler.Create)
			create.POST("/support-records", supportHandler.Create)
			create.POST("/households", householdHandler.Create)
			create.POST("/households/:id/members", householdHandler.AddMember)
			create.POST("/documents", documentHandler.Create)
			create.POST("/pets", petHandler.Create)
		}

		// Update-level access
		update := proj.Group("", projectAuthMW.RequireProjectRole(project.ActionUpdate))
		{
			update.PATCH("/people/:person_id", personHandler.Update)
			update.PATCH("/support-records/:id", supportHandler.Update)
			update.PATCH("/households/:id", householdHandler.Update)
			update.PATCH("/pets/:id", petHandler.Update)
		}

		// Delete-level access
		del := proj.Group("", projectAuthMW.RequireProjectRole(project.ActionDelete))
		{
			del.DELETE("/tags/:id", tagHandler.Delete)
			del.DELETE("/people/:person_id", personHandler.Delete)
			del.DELETE("/people/:person_id/notes/:id", noteHandler.Delete)
			del.DELETE("/support-records/:id", supportHandler.Delete)
			del.DELETE("/households/:id", householdHandler.Delete)
			del.DELETE("/households/:id/members/:person_id", householdHandler.RemoveMember)
			del.DELETE("/documents/:id", documentHandler.Delete)
			del.DELETE("/pets/:id", petHandler.Delete)
		}
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
