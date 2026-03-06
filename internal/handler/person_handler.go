package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

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
// @Summary List people in a project
// @Tags project-people
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param consultant_id query string false "Filter by consultant ID"
// @Param office_id query string false "Filter by office ID"
// @Param case_status query string false "Filter by case status"
// @Param search query string false "Search term"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} ucproject.ListPeopleOutput
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/people [get]
func (h *PersonHandler) List(c *gin.Context) {
	projectID := c.Param("project_id")
	var input ucproject.ListPeopleInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	canContact := middleware.CanViewContactFrom(c)
	canPersonal := middleware.CanViewPersonalFrom(c)
	out, err := h.personUC.List(c.Request.Context(), projectID, input, canContact, canPersonal)
	if err != nil {
		internalError(c, "list people", err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Get handles GET /projects/:project_id/people/:person_id.
// @Summary Get a person by ID
// @Tags project-people
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param person_id path string true "Person ID"
// @Success 200 {object} ucproject.PersonDTO
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/people/{person_id} [get]
func (h *PersonHandler) Get(c *gin.Context) {
	canContact := middleware.CanViewContactFrom(c)
	canPersonal := middleware.CanViewPersonalFrom(c)
	out, err := h.personUC.Get(c.Request.Context(), c.Param("person_id"), canContact, canPersonal)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /projects/:project_id/people.
// @Summary Create a person in a project
// @Tags project-people
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param input body ucproject.CreatePersonInput true "Person payload"
// @Success 201 {object} ucproject.PersonDTO
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/people [post]
func (h *PersonHandler) Create(c *gin.Context) {
	projectID := c.Param("project_id")
	var input ucproject.CreatePersonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	out, err := h.personUC.Create(c.Request.Context(), projectID, input)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Update handles PATCH /projects/:project_id/people/:person_id.
// @Summary Update a person
// @Tags project-people
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param person_id path string true "Person ID"
// @Param input body ucproject.UpdatePersonInput true "Update payload"
// @Success 200 {object} ucproject.PersonDTO
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/people/{person_id} [patch]
func (h *PersonHandler) Update(c *gin.Context) {
	var input ucproject.UpdatePersonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	out, err := h.personUC.Update(c.Request.Context(), c.Param("person_id"), input)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Delete handles DELETE /projects/:project_id/people/:person_id.
// @Summary Delete a person
// @Tags project-people
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param person_id path string true "Person ID"
// @Success 200 {object} MessageResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/people/{person_id} [delete]
func (h *PersonHandler) Delete(c *gin.Context) {
	if err := h.personUC.Delete(c.Request.Context(), c.Param("person_id")); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "person deleted"})
}

// ListCategories handles GET /projects/:project_id/people/:person_id/categories.
// @Summary List category IDs for a person
// @Tags project-people
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param person_id path string true "Person ID"
// @Success 200 {object} object "Wrapper with category_ids array"
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/people/{person_id}/categories [get]
func (h *PersonHandler) ListCategories(c *gin.Context) {
	ids, err := h.categoryUC.List(c.Request.Context(), c.Param("person_id"))
	if err != nil {
		internalError(c, "list person categories", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"category_ids": ids})
}

// ReplaceCategories handles PUT /projects/:project_id/people/:person_id/categories.
// @Summary Replace category IDs for a person
// @Tags project-people
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param person_id path string true "Person ID"
// @Param input body ucproject.ReplaceIDsInput true "Category IDs payload"
// @Success 200 {object} object "Wrapper with category_ids array"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/people/{person_id}/categories [put]
func (h *PersonHandler) ReplaceCategories(c *gin.Context) {
	var input ucproject.ReplaceIDsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	if err := h.categoryUC.Replace(c.Request.Context(), c.Param("person_id"), input.IDs); err != nil {
		internalError(c, "replace person categories", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"category_ids": input.IDs})
}

// ListTags handles GET /projects/:project_id/people/:person_id/tags.
// @Summary List tag IDs for a person
// @Tags project-people
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param person_id path string true "Person ID"
// @Success 200 {object} object "Wrapper with tag_ids array"
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/people/{person_id}/tags [get]
func (h *PersonHandler) ListTags(c *gin.Context) {
	ids, err := h.tagUC.List(c.Request.Context(), c.Param("person_id"))
	if err != nil {
		internalError(c, "list person tags", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"tag_ids": ids})
}

// ReplaceTags handles PUT /projects/:project_id/people/:person_id/tags.
// @Summary Replace tag IDs for a person
// @Tags project-people
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param person_id path string true "Person ID"
// @Param input body ucproject.ReplaceIDsInput true "Tag IDs payload"
// @Success 200 {object} object "Wrapper with tag_ids array"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/people/{person_id}/tags [put]
func (h *PersonHandler) ReplaceTags(c *gin.Context) {
	var input ucproject.ReplaceIDsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	if err := h.tagUC.Replace(c.Request.Context(), c.Param("person_id"), input.IDs); err != nil {
		internalError(c, "replace person tags", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"tag_ids": input.IDs})
}

