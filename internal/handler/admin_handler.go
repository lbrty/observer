package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"

	"github.com/lbrty/observer/internal/repository"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
)

// AdminHandler exposes admin user-management HTTP endpoints.
type AdminHandler struct {
	userUC        *ucadmin.UserUseCase
	loginAttempts repository.LoginAttemptStore
}

// NewAdminHandler creates an AdminHandler.
func NewAdminHandler(userUC *ucadmin.UserUseCase, loginAttempts repository.LoginAttemptStore) *AdminHandler {
	return &AdminHandler{userUC: userUC, loginAttempts: loginAttempts}
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

	out, err := h.userUC.List(c.Request.Context(), input)
	if err != nil {
		internalError(c, "list users", err)
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

	out, err := h.userUC.Get(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
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

	out, err := h.userUC.Update(c.Request.Context(), id, input)
	if err != nil {
		HandleError(c, err)
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

	out, err := h.userUC.Create(c.Request.Context(), input)
	if err != nil {
		HandleError(c, err)
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

	if err := h.userUC.ResetPassword(c.Request.Context(), id, input); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}

// UnlockAccount handles POST /admin/users/:id/unlock.
func (h *AdminHandler) UnlockAccount(c *gin.Context) {
	id, err := ulid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", "invalid user ID"))
		return
	}

	u, err := h.userUC.Get(c.Request.Context(), id)
	if err != nil {
		HandleError(c, err)
		return
	}

	if err := h.loginAttempts.ClearAttempts(c.Request.Context(), u.Email); err != nil {
		internalError(c, "unlock account", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "account unlocked"})
}

