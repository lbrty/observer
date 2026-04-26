package my

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/handler"
	"github.com/lbrty/observer/internal/middleware"
	ucmy "github.com/lbrty/observer/internal/usecase/my"
)

// MyHandler exposes endpoints for the current user's data.
type MyHandler struct {
	projectsUC *ucmy.MyProjectsUseCase
}

// NewMyHandler creates a MyHandler.
func NewMyHandler(projectsUC *ucmy.MyProjectsUseCase) *MyHandler {
	return &MyHandler{projectsUC: projectsUC}
}

// Projects handles GET /my/projects.
func (h *MyHandler) Projects(c *gin.Context) {
	userID, ok := middleware.UserIDFrom(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, handler.ErrJSON("errors.auth.missingUser", "missing user"))
		return
	}

	role, _ := middleware.UserRoleFrom(c)

	out, err := h.projectsUC.Execute(c.Request.Context(), userID.String(), role)
	if err != nil {
		handler.InternalError(c, "list my projects", err)
		return
	}

	c.JSON(http.StatusOK, out)
}
