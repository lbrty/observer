package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/domain/reference"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
)

// CategoryHandler exposes category CRUD HTTP endpoints.
type CategoryHandler struct {
	uc *ucadmin.CategoryUseCase
}

// NewCategoryHandler creates a CategoryHandler.
func NewCategoryHandler(uc *ucadmin.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{uc: uc}
}

// List handles GET /admin/categories.
// @Summary List all categories
// @Tags admin-categories
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object "List of categories"
// @Failure 500 {object} ErrorResponse
// @Router /admin/categories [get]
func (h *CategoryHandler) List(c *gin.Context) {
	out, err := h.uc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"categories": out})
}

// Get handles GET /admin/categories/:id.
// @Summary Get a category by ID
// @Tags admin-categories
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} ucadmin.CategoryDTO
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/categories/{id} [get]
func (h *CategoryHandler) Get(c *gin.Context) {
	out, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /admin/categories.
// @Summary Create a new category
// @Tags admin-categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body ucadmin.CreateCategoryInput true "Category payload"
// @Success 201 {object} ucadmin.CategoryDTO
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	var input ucadmin.CreateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.uc.Create(c.Request.Context(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Update handles PATCH /admin/categories/:id.
// @Summary Update a category
// @Tags admin-categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Param input body ucadmin.UpdateCategoryInput true "Update payload"
// @Success 200 {object} ucadmin.CategoryDTO
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/categories/{id} [patch]
func (h *CategoryHandler) Update(c *gin.Context) {
	var input ucadmin.UpdateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.uc.Update(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Delete handles DELETE /admin/categories/:id.
// @Summary Delete a category
// @Tags admin-categories
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} MessageResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "category deleted"})
}

func (h *CategoryHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, reference.ErrCategoryNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, reference.ErrCategoryNameExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
