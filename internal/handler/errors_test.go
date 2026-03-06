package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

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

func TestMapDomainError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		// User errors
		{"user email exists", user.ErrEmailExists, http.StatusConflict, "errors.user.emailExists"},
		{"user phone exists", user.ErrPhoneExists, http.StatusConflict, "errors.user.phoneExists"},
		{"invalid credentials", user.ErrInvalidCredentials, http.StatusUnauthorized, "errors.auth.invalidCredentials"},
		{"user not active", user.ErrUserNotActive, http.StatusForbidden, "errors.user.notActive"},
		{"user not found", user.ErrUserNotFound, http.StatusNotFound, "errors.user.notFound"},
		{"invalid role", user.ErrInvalidRole, http.StatusBadRequest, "errors.user.invalidRole"},

		// Auth/session errors
		{"session not found", domainauth.ErrSessionNotFound, http.StatusUnauthorized, "errors.auth.sessionNotFound"},
		{"session expired", domainauth.ErrSessionExpired, http.StatusUnauthorized, "errors.auth.sessionExpired"},

		// Project errors
		{"project not found", project.ErrProjectNotFound, http.StatusNotFound, "errors.project.notFound"},
		{"project name exists", project.ErrProjectNameExists, http.StatusConflict, "errors.project.nameExists"},
		{"permission not found", project.ErrPermissionNotFound, http.StatusNotFound, "errors.project.permissionNotFound"},
		{"permission exists", project.ErrPermissionExists, http.StatusConflict, "errors.project.permissionExists"},
		{"invalid project role", project.ErrInvalidProjectRole, http.StatusBadRequest, "errors.project.invalidRole"},

		// Reference errors
		{"country not found", reference.ErrCountryNotFound, http.StatusNotFound, "errors.reference.countryNotFound"},
		{"country code exists", reference.ErrCountryCodeExists, http.StatusConflict, "errors.reference.countryCodeExists"},
		{"state not found", reference.ErrStateNotFound, http.StatusNotFound, "errors.reference.stateNotFound"},
		{"place not found", reference.ErrPlaceNotFound, http.StatusNotFound, "errors.reference.placeNotFound"},
		{"office not found", reference.ErrOfficeNotFound, http.StatusNotFound, "errors.reference.officeNotFound"},
		{"category not found", reference.ErrCategoryNotFound, http.StatusNotFound, "errors.reference.categoryNotFound"},
		{"category name exists", reference.ErrCategoryNameExists, http.StatusConflict, "errors.reference.categoryNameExists"},

		// Tag errors
		{"tag not found", tag.ErrTagNotFound, http.StatusNotFound, "errors.tag.notFound"},
		{"tag name exists", tag.ErrTagNameExists, http.StatusConflict, "errors.tag.nameExists"},

		// Person errors
		{"person not found", person.ErrPersonNotFound, http.StatusNotFound, "errors.person.notFound"},
		{"external id exists", person.ErrExternalIDExists, http.StatusConflict, "errors.person.externalIdExists"},
		{"consent constraint", person.ErrConsentConstraint, http.StatusBadRequest, "errors.person.consentConstraint"},
		{"age constraint", person.ErrAgeConstraint, http.StatusBadRequest, "errors.person.ageConstraint"},

		// Support record errors
		{"support record not found", support.ErrRecordNotFound, http.StatusNotFound, "errors.support.notFound"},

		// Migration record errors
		{"migration record not found", migration.ErrRecordNotFound, http.StatusNotFound, "errors.migration.notFound"},

		// Household errors
		{"household not found", household.ErrHouseholdNotFound, http.StatusNotFound, "errors.household.notFound"},
		{"household member not found", household.ErrMemberNotFound, http.StatusNotFound, "errors.household.memberNotFound"},
		{"household member exists", household.ErrMemberExists, http.StatusConflict, "errors.household.memberExists"},

		// Note errors
		{"note not found", note.ErrNoteNotFound, http.StatusNotFound, "errors.note.notFound"},

		// Document errors
		{"document not found", document.ErrDocumentNotFound, http.StatusNotFound, "errors.document.notFound"},

		// Pet errors
		{"pet not found", pet.ErrPetNotFound, http.StatusNotFound, "errors.pet.notFound"},

		// Default fallback
		{"unknown error", errors.New("surprise"), http.StatusInternalServerError, "errors.internal"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, code := MapDomainError(tt.err)
			assert.Equal(t, tt.wantStatus, status)
			assert.Equal(t, tt.wantCode, code)
		})
	}
}
