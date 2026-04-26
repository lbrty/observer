package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/handler"
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
func (h *CategoryHandler) List(c *gin.Context) {
	out, err := h.uc.List(c.Request.Context())
	if err != nil {
		handler.InternalError(c, "list categories", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"categories": out})
}

// Get handles GET /admin/categories/:id.
func (h *CategoryHandler) Get(c *gin.Context) {
	out, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /admin/categories.
func (h *CategoryHandler) Create(c *gin.Context) {
	input, ok := handler.BindJSON[ucadmin.CreateCategoryInput](c)
	if !ok {
		return
	}
	out, err := h.uc.Create(c.Request.Context(), input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Update handles PATCH /admin/categories/:id.
func (h *CategoryHandler) Update(c *gin.Context) {
	input, ok := handler.BindJSON[ucadmin.UpdateCategoryInput](c)
	if !ok {
		return
	}
	out, err := h.uc.Update(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Delete handles DELETE /admin/categories/:id.
func (h *CategoryHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "category deleted"})
}
