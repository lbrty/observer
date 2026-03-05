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
	listUsersUC     *ucadmin.ListUsersUseCase
	getUserUC       *ucadmin.GetUserUseCase
	updateUserUC    *ucadmin.UpdateUserUseCase
	createUserUC    *ucadmin.CreateUserUseCase
	resetPasswordUC *ucadmin.ResetPasswordUseCase
}

// NewAdminHandler creates an AdminHandler.
func NewAdminHandler(
	listUsersUC *ucadmin.ListUsersUseCase,
	getUserUC *ucadmin.GetUserUseCase,
	updateUserUC *ucadmin.UpdateUserUseCase,
	createUserUC *ucadmin.CreateUserUseCase,
	resetPasswordUC *ucadmin.ResetPasswordUseCase,
) *AdminHandler {
	return &AdminHandler{
		listUsersUC:     listUsersUC,
		getUserUC:       getUserUC,
		updateUserUC:    updateUserUC,
		createUserUC:    createUserUC,
		resetPasswordUC: resetPasswordUC,
	}
}

// ListUsers handles GET /admin/users.
// @Summary List users
// @Tags admin-users
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Param search query string false "Search by name or email"
// @Param role query string false "Filter by role"
// @Param is_active query bool false "Filter by active status"
// @Success 200 {object} ucadmin.ListUsersOutput
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /admin/users [get]
func (h *AdminHandler) ListUsers(c *gin.Context) {
	var input ucadmin.ListUsersInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}

	out, err := h.listUsersUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errJSON("errors.internal", "internal server error"))
		return
	}

	c.JSON(http.StatusOK, out)
}

// GetUser handles GET /admin/users/:id.
// @Summary Get user by ID
// @Tags admin-users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} ucadmin.UserDTO
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /admin/users/{id} [get]
func (h *AdminHandler) GetUser(c *gin.Context) {
	id, err := ulid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", "invalid user ID"))
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
// @Summary Update user by ID
// @Tags admin-users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param input body ucadmin.UpdateUserInput true "User update payload"
// @Success 200 {object} ucadmin.UserDTO
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /admin/users/{id} [patch]
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	id, err := ulid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", "invalid user ID"))
		return
	}

	var input ucadmin.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}

	out, err := h.updateUserUC.Execute(c.Request.Context(), id, input)
	if err != nil {
		h.handleUserError(c, err)
		return
	}

	c.JSON(http.StatusOK, out)
}

// CreateUser handles POST /admin/users.
// @Summary Create a new user
// @Tags admin-users
// @Accept json
// @Produce json
// @Param input body ucadmin.CreateUserInput true "User creation payload"
// @Success 201 {object} ucadmin.UserDTO
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /admin/users [post]
func (h *AdminHandler) CreateUser(c *gin.Context) {
	var input ucadmin.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}

	out, err := h.createUserUC.Execute(c.Request.Context(), input)
	if err != nil {
		h.handleUserError(c, err)
		return
	}

	c.JSON(http.StatusCreated, out)
}

// ResetPassword handles POST /admin/users/:id/reset-password.
func (h *AdminHandler) ResetPassword(c *gin.Context) {
	id, err := ulid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", "invalid user ID"))
		return
	}

	var input ucadmin.ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}

	if err := h.resetPasswordUC.Execute(c.Request.Context(), id, input); err != nil {
		h.handleUserError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}

func (h *AdminHandler) handleUserError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, user.ErrUserNotFound):
		c.JSON(http.StatusNotFound, errJSON("errors.user.notFound", err.Error()))
	case errors.Is(err, user.ErrInvalidRole):
		c.JSON(http.StatusBadRequest, errJSON("errors.user.invalidRole", err.Error()))
	case errors.Is(err, user.ErrEmailExists):
		c.JSON(http.StatusConflict, errJSON("errors.user.emailExists", err.Error()))
	case errors.Is(err, user.ErrPhoneExists):
		c.JSON(http.StatusConflict, errJSON("errors.user.phoneExists", err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, errJSON("errors.internal", "internal server error"))
	}
}
