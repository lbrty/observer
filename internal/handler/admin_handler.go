package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"

	"github.com/lbrty/observer/internal/domain/user"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
)

// AdminHandler exposes admin user-management HTTP endpoints.
type AdminHandler struct {
	listUsersUC  *ucadmin.ListUsersUseCase
	getUserUC    *ucadmin.GetUserUseCase
	updateUserUC *ucadmin.UpdateUserUseCase
}

// NewAdminHandler creates an AdminHandler.
func NewAdminHandler(
	listUsersUC *ucadmin.ListUsersUseCase,
	getUserUC *ucadmin.GetUserUseCase,
	updateUserUC *ucadmin.UpdateUserUseCase,
) *AdminHandler {
	return &AdminHandler{
		listUsersUC:  listUsersUC,
		getUserUC:    getUserUC,
		updateUserUC: updateUserUC,
	}
}

// ListUsers handles GET /admin/users.
func (h *AdminHandler) ListUsers(c *gin.Context) {
	var input ucadmin.ListUsersInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.listUsersUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, out)
}

// GetUser handles GET /admin/users/:id.
func (h *AdminHandler) GetUser(c *gin.Context) {
	id, err := ulid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	out, err := h.getUserUC.Execute(c.Request.Context(), id)
	if err != nil {
		h.handleUserError(c, err)
		return
	}

	c.JSON(http.StatusOK, out)
}

// UpdateUser handles PATCH /admin/users/:id.
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	id, err := ulid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var input ucadmin.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.updateUserUC.Execute(c.Request.Context(), id, input)
	if err != nil {
		h.handleUserError(c, err)
		return
	}

	c.JSON(http.StatusOK, out)
}

func (h *AdminHandler) handleUserError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, user.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, user.ErrInvalidRole):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, user.ErrEmailExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, user.ErrPhoneExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
