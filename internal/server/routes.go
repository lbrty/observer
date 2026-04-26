package server

import (
	"log/slog"

	"github.com/lbrty/observer/internal/app"
	"github.com/lbrty/observer/internal/config"
	"github.com/lbrty/observer/internal/database"
	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/handler"
	adminhandler "github.com/lbrty/observer/internal/handler/admin"
	authhandler "github.com/lbrty/observer/internal/handler/auth"
	audithandler "github.com/lbrty/observer/internal/handler/audit"
	myhandler "github.com/lbrty/observer/internal/handler/my"
	projecthandler "github.com/lbrty/observer/internal/handler/project"
	reporthandler "github.com/lbrty/observer/internal/handler/report"
	searchhandler "github.com/lbrty/observer/internal/handler/search"
	"github.com/lbrty/observer/internal/health"
	"github.com/lbrty/observer/internal/middleware"
	"github.com/lbrty/observer/internal/spa"
)

func (s *Server) setupRoutes(cfg *config.Config, db database.DB, container *app.Container, schemaStatus handler.SchemaStatus) {
	healthHandler := health.NewHandler(db)
	s.router.GET("/health", healthHandler.Health)

	// All API routes live under /api so the SPA can own everything else.
	api := s.router.Group("/api")

	authMW := middleware.NewAuthMiddleware(container.TokenGenerator, container.UserRepo)
	projectAuthMW := middleware.NewProjectAuthMiddleware(container.PermissionRepo)

	authHandler := authhandler.NewAuthHandler(
		container.AuthUC,
		cfg.Cookie,
		cfg.JWT,
	)

	loginRL := middleware.RateLimit(float64(cfg.RateLimit.LoginRate)/60.0, cfg.RateLimit.LoginRate)
	registerRL := middleware.RateLimit(float64(cfg.RateLimit.RegisterRate)/60.0, cfg.RateLimit.RegisterRate)

	auth := api.Group("/auth")
	{
		auth.POST("/register", registerRL, authHandler.Register)
		auth.POST("/login", loginRL, authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/mfa", authHandler.VerifyMFA)
		auth.GET("/me", authMW.Authenticate(), authHandler.Me)
		auth.PATCH("/me", authMW.Authenticate(), authHandler.UpdateProfile)
		auth.POST("/change-password", authMW.Authenticate(), authHandler.ChangePassword)
		auth.POST("/logout", authMW.Authenticate(), authHandler.Logout)
		auth.GET("/mfa/setup", authMW.Authenticate(), authHandler.MFASetup)
		auth.POST("/mfa/enable", authMW.Authenticate(), authHandler.EnableMFA)
		auth.POST("/mfa/disable", authMW.Authenticate(), authHandler.DisableMFA)
	}

	// My endpoints — authenticated user's own data
	myHandler := myhandler.NewMyHandler(container.MyProjectsUC)
	my := api.Group("/my", authMW.Authenticate())
	{
		my.GET("/projects", myHandler.Projects)
	}

	// Global search — any authenticated user, scoped to their authorised projects
	searchHandler := searchhandler.NewSearchHandler(container.SearchUC)
	api.GET("/search", authMW.Authenticate(), searchHandler.Search)

	// Admin endpoints — requires authentication + admin role
	adminHandler := adminhandler.NewAdminHandler(container.UserUC)
	permHandler := adminhandler.NewPermissionHandler(container.PermUC)
	countryHandler := adminhandler.NewCountryHandler(container.CountryUC)
	stateHandler := adminhandler.NewStateHandler(container.StateUC)
	placeHandler := adminhandler.NewPlaceHandler(container.PlaceUC)
	officeHandler := adminhandler.NewOfficeHandler(container.OfficeUC)
	categoryHandler := adminhandler.NewCategoryHandler(container.CategoryUC)
	projectHandler := adminhandler.NewProjectHandler(container.ProjectUC)

	// Admin endpoints readable by admin + staff + consultant
	adminRead := api.Group("/admin", authMW.Authenticate(), authMW.RequireRole(user.RoleAdmin, user.RoleStaff, user.RoleConsultant))
	{
		adminRead.GET("/users", adminHandler.ListUsers)
		adminRead.GET("/users/:id", adminHandler.GetUser)
	}

	// Reference data write endpoints (admin + staff + consultant)
	adminWrite := api.Group("/admin", authMW.Authenticate(), authMW.RequireRole(user.RoleAdmin, user.RoleStaff, user.RoleConsultant))
	{
		adminWrite.POST("/countries", countryHandler.Create)
		adminWrite.PATCH("/countries/:id", countryHandler.Update)

		adminWrite.POST("/states", stateHandler.Create)
		adminWrite.PATCH("/states/:id", stateHandler.Update)

		adminWrite.POST("/places", placeHandler.Create)
		adminWrite.PATCH("/places/:id", placeHandler.Update)

		adminWrite.POST("/categories", categoryHandler.Create)
		adminWrite.PATCH("/categories/:id", categoryHandler.Update)
	}

	// Project & permission read endpoints — any authenticated user.
	// Admin/Staff see all; Consultant/Guest see only what they have permissions for.
	adminAny := api.Group("/admin", authMW.Authenticate())
	{
		adminAny.GET("/projects", projectHandler.List)
		adminAny.GET("/projects/:project_id", projectHandler.Get)
		adminAny.GET("/projects/:project_id/permissions", permHandler.ListPermissions)

		// Reference data — readable by all authenticated users including guests.
		adminAny.GET("/countries", countryHandler.List)
		adminAny.GET("/countries/:id", countryHandler.Get)
		adminAny.GET("/states", stateHandler.List)
		adminAny.GET("/states/:id", stateHandler.Get)
		adminAny.GET("/places", placeHandler.List)
		adminAny.GET("/places/:id", placeHandler.Get)
		adminAny.GET("/offices", officeHandler.List)
		adminAny.GET("/offices/:id", officeHandler.Get)
		adminAny.GET("/categories", categoryHandler.List)
		adminAny.GET("/categories/:id", categoryHandler.Get)
	}

	// Admin-only endpoints
	auditHandler := audithandler.NewAuditHandler(container.AuditUC)
	schemaHandler := handler.NewSchemaHandler(schemaStatus)
	admin := api.Group("/admin", authMW.Authenticate(), authMW.RequireRole(user.RoleAdmin))
	{
		admin.GET("/schema/status", schemaHandler.Status)
		admin.GET("/audit-logs", auditHandler.ListAll)

		admin.POST("/users", adminHandler.CreateUser)
		admin.PATCH("/users/:id", adminHandler.UpdateUser)
		admin.POST("/users/:id/reset-password", adminHandler.ResetPassword)
		admin.POST("/users/:id/unlock", adminHandler.UnlockAccount)
		admin.PATCH("/users/:id/deactivate", adminHandler.DeactivateUser)
		admin.PATCH("/users/:id/reactivate", adminHandler.ReactivateUser)

		admin.POST("/projects", projectHandler.Create)
		admin.PATCH("/projects/:project_id", projectHandler.Update)

		admin.POST("/projects/:project_id/permissions", permHandler.AssignPermission)
		admin.PATCH("/projects/:project_id/permissions/:id", permHandler.UpdatePermission)
		admin.DELETE("/projects/:project_id/permissions/:id", permHandler.RevokePermission)

		admin.DELETE("/countries/:id", countryHandler.Delete)
		admin.DELETE("/states/:id", stateHandler.Delete)
		admin.DELETE("/places/:id", placeHandler.Delete)

		admin.POST("/offices", officeHandler.Create)
		admin.PATCH("/offices/:id", officeHandler.Update)
		admin.DELETE("/offices/:id", officeHandler.Delete)

		admin.DELETE("/categories/:id", categoryHandler.Delete)
	}

	// Project-scoped endpoints — requires authentication + project role
	tagHandler := projecthandler.NewTagHandler(container.TagUC)
	personHandler := projecthandler.NewPersonHandler(container.PersonUC, container.PersonCategoryUC, container.PersonTagUC)
	supportHandler := projecthandler.NewSupportRecordHandler(container.SupportRecordUC)
	migrationHandler := projecthandler.NewMigrationRecordHandler(container.MigrationRecordUC)
	householdHandler := projecthandler.NewHouseholdHandler(container.HouseholdUC)
	noteHandler := projecthandler.NewNoteHandler(container.NoteUC)
	documentHandler := projecthandler.NewDocumentHandler(container.DocumentUC)
	petHandler := projecthandler.NewPetHandler(container.PetUC, container.PetTagUC)
	reportHandler := reporthandler.NewReportHandler(container.ReportUC)
	petReportHandler := reporthandler.NewPetReportHandler(container.PetReportUC)
	exportHandler := reporthandler.NewExportHandler(container.PersonUC, container.SupportRecordUC, container.PetUC, container.HouseholdUC, container.AuditUC)

	proj := api.Group("/projects/:project_id", authMW.Authenticate())
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
			read.GET("/documents/:id/download", documentHandler.Download)
			read.GET("/documents/:id/stream", documentHandler.Stream)
			read.GET("/documents/:id/thumbnail", documentHandler.Thumbnail)
			read.GET("/pets", petHandler.List)
			read.GET("/pets/:id", petHandler.Get)
			read.GET("/pets/:id/tags", petHandler.ListTags)
			read.GET("/reports", reportHandler.Generate)
			read.GET("/reports/custom", reportHandler.GenerateCustom)
			read.GET("/reports/pets", petReportHandler.Generate)
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
			create.POST("/people/:person_id/documents", documentHandler.Upload)
			create.POST("/pets", petHandler.Create)
			create.PUT("/pets/:id/tags", petHandler.ReplaceTags)
		}

		// Update-level access
		update := proj.Group("", projectAuthMW.RequireProjectRole(project.ActionUpdate))
		{
			update.PATCH("/tags/:id", tagHandler.Update)
			update.PATCH("/people/:person_id", personHandler.Update)
			update.PATCH("/support-records/:id", supportHandler.Update)
			update.PATCH("/households/:id", householdHandler.Update)
			update.PATCH("/people/:person_id/migration-records/:id", migrationHandler.Update)
			update.PATCH("/people/:person_id/notes/:id", noteHandler.Update)
			update.PATCH("/documents/:id", documentHandler.Update)
			update.PATCH("/pets/:id", petHandler.Update)
		}

		// Export-level access (consultant+)
		export := proj.Group("", projectAuthMW.RequireProjectRole(project.ActionExport), projectAuthMW.RequireExport())
		{
			export.GET("/export/people", exportHandler.ExportPeople)
			export.GET("/export/support-records", exportHandler.ExportSupportRecords)
			export.GET("/export/pets", exportHandler.ExportPets)
			export.GET("/export/households", exportHandler.ExportHouseholds)
		}

		// Delete-level access (manager+)
		del := proj.Group("", projectAuthMW.RequireProjectRole(project.ActionDelete))
		{
			del.GET("/audit-logs", auditHandler.ListByProject)
			del.DELETE("/tags/:id", tagHandler.Delete)
			del.DELETE("/people/:person_id", personHandler.Delete)
			del.DELETE("/people/:person_id/notes/:id", noteHandler.Delete)
			del.DELETE("/support-records/:id", supportHandler.Delete)
			del.DELETE("/households/:id", householdHandler.Delete)
			del.DELETE("/households/:id/members/:person_id", householdHandler.RemoveMember)
			del.DELETE("/documents/:id", documentHandler.Delete)
			del.DELETE("/pets/:id", petHandler.Delete)
			del.DELETE("/people/:person_id/migration-records/:id", migrationHandler.Delete)
		}
	}

	// Serve embedded SPA in production builds
	if spa.Enabled() {
		spaFS, err := spa.FS()
		if err != nil {
			slog.Error("failed to load embedded SPA", slog.Any("err", err))
		} else {
			slog.Info("serving embedded SPA")
			spa.Mount(s.router, spaFS)
		}
	}
}
