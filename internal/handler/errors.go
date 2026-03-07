package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	domainauth "github.com/lbrty/observer/internal/domain/auth"
	"github.com/lbrty/observer/internal/domain/document"
	"github.com/lbrty/observer/internal/domain/household"
	"github.com/lbrty/observer/internal/domain/migration"
	"github.com/lbrty/observer/internal/domain/note"
	"github.com/lbrty/observer/internal/domain/person"
	"github.com/lbrty/observer/internal/domain/pet"
	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/domain/reference"
	"github.com/lbrty/observer/internal/domain/support"
	"github.com/lbrty/observer/internal/domain/tag"
	"github.com/lbrty/observer/internal/domain/user"
)

// errJSON builds a JSON error response with both a human-readable message and
// a machine-readable code that maps to a frontend i18n key.
func errJSON(code, msg string) gin.H {
	return gin.H{"error": msg, "code": code}
}

// internalError logs the error with context and returns a 500 JSON response.
func internalError(c *gin.Context, msg string, err error) {
	slog.Error(msg, slog.String("error", err.Error()), slog.String("path", c.Request.URL.Path), slog.String("method", c.Request.Method))
	c.JSON(http.StatusInternalServerError, errJSON("errors.internal", "internal server error"))
}

// MapDomainError maps a domain error to HTTP status and i18n code.
func MapDomainError(err error) (int, string) {
	switch {
	// User errors
	case errors.Is(err, user.ErrEmailExists):
		return http.StatusConflict, "errors.user.emailExists"
	case errors.Is(err, user.ErrPhoneExists):
		return http.StatusConflict, "errors.user.phoneExists"
	case errors.Is(err, user.ErrInvalidCredentials):
		return http.StatusUnauthorized, "errors.auth.invalidCredentials"
	case errors.Is(err, user.ErrUserNotActive):
		return http.StatusForbidden, "errors.user.notActive"
	case errors.Is(err, user.ErrUserNotFound):
		return http.StatusNotFound, "errors.user.notFound"
	case errors.Is(err, user.ErrInvalidRole):
		return http.StatusBadRequest, "errors.user.invalidRole"

	// Auth/session errors
	case errors.Is(err, domainauth.ErrSessionNotFound):
		return http.StatusUnauthorized, "errors.auth.sessionNotFound"
	case errors.Is(err, domainauth.ErrSessionExpired):
		return http.StatusUnauthorized, "errors.auth.sessionExpired"

	// Project errors
	case errors.Is(err, project.ErrProjectNotFound):
		return http.StatusNotFound, "errors.project.notFound"
	case errors.Is(err, project.ErrProjectNameExists):
		return http.StatusConflict, "errors.project.nameExists"
	case errors.Is(err, project.ErrPermissionNotFound):
		return http.StatusNotFound, "errors.project.permissionNotFound"
	case errors.Is(err, project.ErrPermissionExists):
		return http.StatusConflict, "errors.project.permissionExists"
	case errors.Is(err, project.ErrInvalidProjectRole):
		return http.StatusBadRequest, "errors.project.invalidRole"

	// Reference errors
	case errors.Is(err, reference.ErrCountryNotFound):
		return http.StatusNotFound, "errors.reference.countryNotFound"
	case errors.Is(err, reference.ErrCountryCodeExists):
		return http.StatusConflict, "errors.reference.countryCodeExists"
	case errors.Is(err, reference.ErrStateNotFound):
		return http.StatusNotFound, "errors.reference.stateNotFound"
	case errors.Is(err, reference.ErrPlaceNotFound):
		return http.StatusNotFound, "errors.reference.placeNotFound"
	case errors.Is(err, reference.ErrOfficeNotFound):
		return http.StatusNotFound, "errors.reference.officeNotFound"
	case errors.Is(err, reference.ErrCategoryNotFound):
		return http.StatusNotFound, "errors.reference.categoryNotFound"
	case errors.Is(err, reference.ErrCategoryNameExists):
		return http.StatusConflict, "errors.reference.categoryNameExists"

	// Tag errors
	case errors.Is(err, tag.ErrTagNotFound):
		return http.StatusNotFound, "errors.tag.notFound"
	case errors.Is(err, tag.ErrTagNameExists):
		return http.StatusConflict, "errors.tag.nameExists"

	// Person errors
	case errors.Is(err, person.ErrPersonNotFound):
		return http.StatusNotFound, "errors.person.notFound"
	case errors.Is(err, person.ErrExternalIDExists):
		return http.StatusConflict, "errors.person.externalIdExists"
	case errors.Is(err, person.ErrConsentConstraint):
		return http.StatusBadRequest, "errors.person.consentConstraint"
	case errors.Is(err, person.ErrAgeConstraint):
		return http.StatusBadRequest, "errors.person.ageConstraint"

	// Support record errors
	case errors.Is(err, support.ErrRecordNotFound):
		return http.StatusNotFound, "errors.support.notFound"

	// Migration record errors
	case errors.Is(err, migration.ErrRecordNotFound):
		return http.StatusNotFound, "errors.migration.notFound"

	// Household errors
	case errors.Is(err, household.ErrHouseholdNotFound):
		return http.StatusNotFound, "errors.household.notFound"
	case errors.Is(err, household.ErrMemberNotFound):
		return http.StatusNotFound, "errors.household.memberNotFound"
	case errors.Is(err, household.ErrMemberExists):
		return http.StatusConflict, "errors.household.memberExists"

	// Note errors
	case errors.Is(err, note.ErrNoteNotFound):
		return http.StatusNotFound, "errors.note.notFound"

	// Document errors
	case errors.Is(err, document.ErrDocumentNotFound):
		return http.StatusNotFound, "errors.document.notFound"
	case errors.Is(err, document.ErrNotImage):
		return http.StatusBadRequest, "errors.document.notImage"

	// Pet errors
	case errors.Is(err, pet.ErrPetNotFound):
		return http.StatusNotFound, "errors.pet.notFound"

	default:
		return http.StatusInternalServerError, "errors.internal"
	}
}

// HandleError writes a JSON error response for a domain error.
func HandleError(c *gin.Context, err error) {
	status, code := MapDomainError(err)
	if status == http.StatusInternalServerError {
		slog.ErrorContext(c.Request.Context(), "unexpected error",
			"error", err, "path", c.Request.URL.Path)
	}
	c.JSON(status, errJSON(code, err.Error()))
}
