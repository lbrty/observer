package search

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/handler"
	"github.com/lbrty/observer/internal/middleware"
	ucsearch "github.com/lbrty/observer/internal/usecase/search"
)

// SearchHandler exposes the global cross-project search endpoint.
type SearchHandler struct {
	uc *ucsearch.SearchUseCase
}

// NewSearchHandler creates a SearchHandler.
func NewSearchHandler(uc *ucsearch.SearchUseCase) *SearchHandler {
	return &SearchHandler{uc: uc}
}

// Search handles GET /api/search?q=...&limit=...
func (h *SearchHandler) Search(c *gin.Context) {
	userID, ok := middleware.UserIDFrom(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, handler.ErrJSON("errors.auth.missingAuthorization", "unauthorized"))
		return
	}

	role, _ := middleware.UserRoleFrom(c)

	q := strings.TrimSpace(c.Query("q"))
	if len(q) < 2 {
		c.JSON(http.StatusBadRequest, handler.ErrJSON("errors.search.queryTooShort", "query must be at least 2 characters"))
		return
	}

	limit := 5
	if ls := c.Query("limit"); ls != "" {
		n, err := strconv.Atoi(ls)
		if err == nil {
			limit = n
		}
	}
	if limit < 5 {
		limit = 5
	}
	if limit > 50 {
		limit = 50
	}

	out, err := h.uc.Execute(c.Request.Context(), userID.String(), role, q, limit)
	if err != nil {
		handler.InternalError(c, "search", err)
		return
	}

	c.JSON(http.StatusOK, out)
}
