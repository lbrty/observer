package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"

	"github.com/lbrty/observer/internal/handler"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
)

// AdminHandler exposes admin user-management HTTP endpoints.
type AdminHandler struct {
	userUC *ucadmin.UserUseCase
}

// NewAdminHandler creates an AdminHandler.
func NewAdminHandler(userUC *ucadmin.UserUseCase) *AdminHandler {
	return &AdminHandler{userUC: userUC}
}

// ListUsers handles GET /admin/users.
func (h *AdminHandler) ListUsers(c *gin.Context) {
	var input ucadmin.ListUsersInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrJSON("errors.validation", err.Error()))
		return
	}

	out, err := h.userUC.List(c.Request.Context(), input)
	if err != nil {
		handler.InternalError(c, "list users", err)
		return
	}

	c.JSON(http.StatusOK, out)
}

// GetUser handles GET /admin/users/:id.
func (h *AdminHandler) GetUser(c *gin.Context) {
	id, err := ulid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrJSON("errors.validation", "invalid user ID"))
		return
	}

	out, err := h.userUC.Get(c.Request.Context(), id)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, out)
}

// UpdateUser handles PATCH /admin/users/:id.
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	id, err := ulid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrJSON("errors.validation", "invalid user ID"))
		return
	}

	input, ok := handler.BindJSON[ucadmin.UpdateUserInput](c)
	if !ok {
		return
	}

	out, err := h.userUC.Update(c.Request.Context(), id, input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, out)
}

// CreateUser handles POST /admin/users.
func (h *AdminHandler) CreateUser(c *gin.Context) {
	input, ok := handler.BindJSON[ucadmin.CreateUserInput](c)
	if !ok {
		return
	}

	out, err := h.userUC.Create(c.Request.Context(), input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, out)
}

// ResetPassword handles POST /admin/users/:id/reset-password.
func (h *AdminHandler) ResetPassword(c *gin.Context) {
	id, err := ulid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrJSON("errors.validation", "invalid user ID"))
		return
	}

	input, ok := handler.BindJSON[ucadmin.ResetPasswordInput](c)
	if !ok {
		return
	}

	if err := h.userUC.ResetPassword(c.Request.Context(), id, input); err != nil {
		handler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}

// UnlockAccount handles POST /admin/users/:id/unlock.
func (h *AdminHandler) UnlockAccount(c *gin.Context) {
	id, err := ulid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrJSON("errors.validation", "invalid user ID"))
		return
	}
	if err := h.userUC.UnlockAccount(c.Request.Context(), id); err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "account unlocked"})
}

// DeactivateUser handles PATCH /admin/users/:id/deactivate.
func (h *AdminHandler) DeactivateUser(c *gin.Context) {
	id, err := ulid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrJSON("errors.validation", "invalid user ID"))
		return
	}
	out, err := h.userUC.DeactivateUser(c.Request.Context(), id)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// ReactivateUser handles PATCH /admin/users/:id/reactivate.
func (h *AdminHandler) ReactivateUser(c *gin.Context) {
	id, err := ulid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrJSON("errors.validation", "invalid user ID"))
		return
	}
	out, err := h.userUC.ReactivateUser(c.Request.Context(), id)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}
