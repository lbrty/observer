package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/handler"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
)

// StateHandler exposes state CRUD HTTP endpoints.
type StateHandler struct {
	uc *ucadmin.StateUseCase
}

// NewStateHandler creates a StateHandler.
func NewStateHandler(uc *ucadmin.StateUseCase) *StateHandler {
	return &StateHandler{uc: uc}
}

// List handles GET /admin/states.
func (h *StateHandler) List(c *gin.Context) {
	countryID := c.Query("country_id")
	var (
		out []ucadmin.StateDTO
		err error
	)
	if countryID != "" {
		out, err = h.uc.List(c.Request.Context(), countryID)
	} else {
		out, err = h.uc.ListAll(c.Request.Context())
	}
	if err != nil {
		handler.InternalError(c, "list states", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"states": out})
}

// Get handles GET /admin/states/:id.
func (h *StateHandler) Get(c *gin.Context) {
	out, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /admin/states.
func (h *StateHandler) Create(c *gin.Context) {
	countryID := c.Query("country_id")
	if countryID == "" {
		c.JSON(http.StatusBadRequest, handler.ErrJSON("errors.validation", "country_id is required"))
		return
	}
	input, ok := handler.BindJSON[ucadmin.CreateStateInput](c)
	if !ok {
		return
	}
	out, err := h.uc.Create(c.Request.Context(), countryID, input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Update handles PATCH /admin/states/:id.
func (h *StateHandler) Update(c *gin.Context) {
	input, ok := handler.BindJSON[ucadmin.UpdateStateInput](c)
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

// Delete handles DELETE /admin/states/:id.
func (h *StateHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "state deleted"})
}
