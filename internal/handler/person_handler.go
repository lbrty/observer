package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/domain/person"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	canContact := middleware.CanViewContactFrom(c)
	canPersonal := middleware.CanViewPersonalFrom(c)
	out, err := h.personUC.List(c.Request.Context(), projectID, input, canContact, canPersonal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, out)
}

// Get handles GET /projects/:project_id/people/:id.
func (h *PersonHandler) Get(c *gin.Context) {
	canContact := middleware.CanViewContactFrom(c)
	canPersonal := middleware.CanViewPersonalFrom(c)
	out, err := h.personUC.Get(c.Request.Context(), c.Param("id"), canContact, canPersonal)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /projects/:project_id/people.
func (h *PersonHandler) Create(c *gin.Context) {
	projectID := c.Param("project_id")
	var input ucproject.CreatePersonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.personUC.Create(c.Request.Context(), projectID, input)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Update handles PATCH /projects/:project_id/people/:id.
func (h *PersonHandler) Update(c *gin.Context) {
	var input ucproject.UpdatePersonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.personUC.Update(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Delete handles DELETE /projects/:project_id/people/:id.
func (h *PersonHandler) Delete(c *gin.Context) {
	if err := h.personUC.Delete(c.Request.Context(), c.Param("id")); err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "person deleted"})
}

// ListCategories handles GET /projects/:project_id/people/:person_id/categories.
func (h *PersonHandler) ListCategories(c *gin.Context) {
	ids, err := h.categoryUC.List(c.Request.Context(), c.Param("person_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"category_ids": ids})
}

// ReplaceCategories handles PUT /projects/:project_id/people/:person_id/categories.
func (h *PersonHandler) ReplaceCategories(c *gin.Context) {
	var input ucproject.ReplaceIDsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.categoryUC.Replace(c.Request.Context(), c.Param("person_id"), input.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"category_ids": input.IDs})
}

// ListTags handles GET /projects/:project_id/people/:person_id/tags.
func (h *PersonHandler) ListTags(c *gin.Context) {
	ids, err := h.tagUC.List(c.Request.Context(), c.Param("person_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tag_ids": ids})
}

// ReplaceTags handles PUT /projects/:project_id/people/:person_id/tags.
func (h *PersonHandler) ReplaceTags(c *gin.Context) {
	var input ucproject.ReplaceIDsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.tagUC.Replace(c.Request.Context(), c.Param("person_id"), input.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tag_ids": input.IDs})
}

func (h *PersonHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, person.ErrPersonNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, person.ErrExternalIDExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, person.ErrConsentConstraint),
		errors.Is(err, person.ErrAgeConstraint):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
