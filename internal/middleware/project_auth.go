package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/repository"
)

// ProjectAuthMiddleware provides project-scoped authorization.
type ProjectAuthMiddleware struct {
	permLoader repository.PermissionLoader
}

// NewProjectAuthMiddleware creates a ProjectAuthMiddleware.
func NewProjectAuthMiddleware(permLoader repository.PermissionLoader) *ProjectAuthMiddleware {
	return &ProjectAuthMiddleware{permLoader: permLoader}
}

// RequireProjectRole returns middleware that checks the user holds a sufficient
// project role for the given action. The route must contain :project_id.
func (m *ProjectAuthMiddleware) RequireProjectRole(action project.Action) gin.HandlerFunc {
	minRole, ok := project.MinRoleForAction[action]
	if !ok {
		panic("unknown project action: " + string(action))
	}

	return func(c *gin.Context) {
		projectID := c.Param("project_id")
		if projectID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing project_id", "code": "errors.validation"})
			c.Abort()
			return
		}

		userID, ok := UserIDFrom(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated", "code": "errors.auth.missingUser"})
			c.Abort()
			return
		}

		userRole, _ := UserRoleFrom(c)

		// Platform admin bypass — implicit owner on all projects.
		if userRole == "admin" {
			setProjectContext(c, projectID, project.ProjectRoleOwner, true, true, true)
			c.Next()
			return
		}

		// Project owner bypass.
		isOwner, err := m.permLoader.IsProjectOwner(c.Request.Context(), userID, projectID)
		if err != nil && !errors.Is(err, project.ErrProjectNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "authorization check failed", "code": "errors.internal"})
			c.Abort()
			return
		}
		if errors.Is(err, project.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "project not found", "code": "errors.project.notFound"})
			c.Abort()
			return
		}
		if isOwner {
			setProjectContext(c, projectID, project.ProjectRoleOwner, true, true, true)
			c.Next()
			return
		}

		// Load project permission.
		perm, err := m.permLoader.GetPermission(c.Request.Context(), userID, projectID)
		if err != nil {
			if errors.Is(err, project.ErrPermissionNotFound) {
				c.JSON(http.StatusForbidden, gin.H{"error": "no project access", "code": "errors.project.permissionDenied"})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "authorization check failed", "code": "errors.internal"})
			c.Abort()
			return
		}

		if perm.Role.Rank() < minRole.Rank() {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient project permissions", "code": "errors.project.permissionDenied"})
			c.Abort()
			return
		}

		setProjectContext(c, projectID, perm.Role, perm.CanViewContact, perm.CanViewPersonal, perm.CanViewDocuments)
		c.Next()
	}
}

func setProjectContext(c *gin.Context, projectID string, role project.ProjectRole, contact, personal, documents bool) {
	c.Set(string(CtxProjectID), projectID)
	c.Set(string(CtxProjectRole), string(role))
	c.Set(string(CtxCanViewContact), contact)
	c.Set(string(CtxCanViewPersonal), personal)
	c.Set(string(CtxCanViewDocuments), documents)
}

// ProjectRoleFrom extracts the project role from the Gin context.
func ProjectRoleFrom(c *gin.Context) (project.ProjectRole, bool) {
	val, exists := c.Get(string(CtxProjectRole))
	if !exists {
		return "", false
	}
	role, ok := val.(string)
	return project.ProjectRole(role), ok
}

// CanViewContactFrom extracts the can_view_contact flag from the Gin context.
func CanViewContactFrom(c *gin.Context) bool {
	val, _ := c.Get(string(CtxCanViewContact))
	b, _ := val.(bool)
	return b
}

// CanViewPersonalFrom extracts the can_view_personal flag from the Gin context.
func CanViewPersonalFrom(c *gin.Context) bool {
	val, _ := c.Get(string(CtxCanViewPersonal))
	b, _ := val.(bool)
	return b
}

// CanViewDocumentsFrom extracts the can_view_documents flag from the Gin context.
func CanViewDocumentsFrom(c *gin.Context) bool {
	val, _ := c.Get(string(CtxCanViewDocuments))
	b, _ := val.(bool)
	return b
}
