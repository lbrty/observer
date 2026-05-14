package project

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/handler"
	"github.com/lbrty/observer/internal/middleware"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

// PersonHandler exposes person CRUD HTTP endpoints.
type PersonHandler struct {
	personUC   *ucproject.PersonUseCase
	categoryUC *ucproject.PersonCategoryUseCase
	tagUC      *ucproject.PersonTagUseCase
}

// NewPersonHandler creates a PersonHandler.
func NewPersonHandler(
	personUC *ucproject.PersonUseCase,
	categoryUC *ucproject.PersonCategoryUseCase,
	tagUC *ucproject.PersonTagUseCase,
) *PersonHandler {
	return &PersonHandler{personUC: personUC, categoryUC: categoryUC, tagUC: tagUC}
}

// List handles GET /projects/:project_id/people.
func (h *PersonHandler) List(c *gin.Context) {
	projectID := c.Param("project_id")
	var input ucproject.ListPeopleInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrJSON("errors.validation", err.Error()))
		return
	}
	canContact := middleware.CanViewContactFrom(c)
	canPersonal := middleware.CanViewPersonalFrom(c)
	out, err := h.personUC.List(c.Request.Context(), projectID, input, canContact, canPersonal)
	if err != nil {
		handler.InternalError(c, "list people", err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Get handles GET /projects/:project_id/people/:person_id.
func (h *PersonHandler) Get(c *gin.Context) {
	canContact := middleware.CanViewContactFrom(c)
	canPersonal := middleware.CanViewPersonalFrom(c)
	out, err := h.personUC.Get(c.Request.Context(), c.Param("project_id"), c.Param("person_id"), canContact, canPersonal)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /projects/:project_id/people.
func (h *PersonHandler) Create(c *gin.Context) {
	projectID := c.Param("project_id")
	input, ok := handler.BindJSON[ucproject.CreatePersonInput](c)
	if !ok {
		return
	}
	out, err := h.personUC.Create(c.Request.Context(), projectID, input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Update handles PATCH /projects/:project_id/people/:person_id.
func (h *PersonHandler) Update(c *gin.Context) {
	input, ok := handler.BindJSON[ucproject.UpdatePersonInput](c)
	if !ok {
		return
	}
	out, err := h.personUC.Update(c.Request.Context(), c.Param("project_id"), c.Param("person_id"), input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Delete handles DELETE /projects/:project_id/people/:person_id.
func (h *PersonHandler) Delete(c *gin.Context) {
	if err := h.personUC.Delete(c.Request.Context(), c.Param("project_id"), c.Param("person_id")); err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "person deleted"})
}

// ListCategories handles GET /projects/:project_id/people/:person_id/categories.
func (h *PersonHandler) ListCategories(c *gin.Context) {
	ids, err := h.categoryUC.List(c.Request.Context(), c.Param("project_id"), c.Param("person_id"))
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"category_ids": ids})
}

// ReplaceCategories handles PUT /projects/:project_id/people/:person_id/categories.
func (h *PersonHandler) ReplaceCategories(c *gin.Context) {
	input, ok := handler.BindJSON[ucproject.ReplaceIDsInput](c)
	if !ok {
		return
	}
	if err := h.categoryUC.Replace(c.Request.Context(), c.Param("project_id"), c.Param("person_id"), input.IDs); err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"category_ids": input.IDs})
}

// ListTags handles GET /projects/:project_id/people/:person_id/tags.
func (h *PersonHandler) ListTags(c *gin.Context) {
	ids, err := h.tagUC.List(c.Request.Context(), c.Param("project_id"), c.Param("person_id"))
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"tag_ids": ids})
}

// ReplaceTags handles PUT /projects/:project_id/people/:person_id/tags.
func (h *PersonHandler) ReplaceTags(c *gin.Context) {
	input, ok := handler.BindJSON[ucproject.ReplaceIDsInput](c)
	if !ok {
		return
	}
	if err := h.tagUC.Replace(c.Request.Context(), c.Param("project_id"), c.Param("person_id"), input.IDs); err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"tag_ids": input.IDs})
}
