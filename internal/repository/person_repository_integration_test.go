//go:build !short

package repository_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lbrty/observer/internal/domain/person"
	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/repository"
	iulid "github.com/lbrty/observer/internal/ulid"
)

func setupProjectWithOwner(t *testing.T, userRepo repository.UserRepository, projectRepo repository.ProjectRepository) (*user.User, *project.Project) {
	t.Helper()
	ctx := context.Background()
	now := time.Now().UTC()

	owner := &user.User{
		ID:         ulid.Make(),
		FirstName:  "Owner",
		LastName:   "Person",
		Email:      "person-owner-" + iulid.NewString() + "@test.com",
		Phone:      "+" + iulid.NewString()[:11],
		Role:       user.RoleAdmin,
		IsVerified: true,
		IsActive:   true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	require.NoError(t, userRepo.Create(ctx, owner))

	proj := &project.Project{
		ID:      iulid.NewString(),
		Name:    "Person Project " + iulid.NewString(),
		OwnerID: owner.ID.String(),
		Status:  project.ProjectStatusActive,
	}
	require.NoError(t, projectRepo.Create(ctx, proj))

	return owner, proj
}

func makePerson(projectID, firstName string) *person.Person {
	return &person.Person{
		ID:           iulid.NewString(),
		ProjectID:    projectID,
		FirstName:    firstName,
		Sex:          person.SexMale,
		CaseStatus:   person.CaseStatusNew,
		PhoneNumbers: json.RawMessage("[]"),
	}
}

func TestPersonRepo_CRUD(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	personRepo := repository.NewPersonRepository(db)
	ctx := context.Background()

	_, proj := setupProjectWithOwner(t, userRepo, projectRepo)

	p := makePerson(proj.ID, "Aibek")

	// Create
	require.NoError(t, personRepo.Create(ctx, p))

	// Get
	got, err := personRepo.GetByID(ctx, p.ID)
	require.NoError(t, err)
	assert.Equal(t, "Aibek", got.FirstName)
	assert.Equal(t, proj.ID, got.ProjectID)
	assert.Equal(t, person.SexMale, got.Sex)
	assert.Equal(t, person.CaseStatusNew, got.CaseStatus)

	// Update
	p.FirstName = "Aibek Updated"
	p.CaseStatus = person.CaseStatusActive
	require.NoError(t, personRepo.Update(ctx, p))
	got, err = personRepo.GetByID(ctx, p.ID)
	require.NoError(t, err)
	assert.Equal(t, "Aibek Updated", got.FirstName)
	assert.Equal(t, person.CaseStatusActive, got.CaseStatus)

	// Delete
	require.NoError(t, personRepo.Delete(ctx, p.ID))
	_, err = personRepo.GetByID(ctx, p.ID)
	assert.ErrorIs(t, err, person.ErrPersonNotFound)
}

func TestPersonRepo_List_WithFilter(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	personRepo := repository.NewPersonRepository(db)
	ctx := context.Background()

	_, proj := setupProjectWithOwner(t, userRepo, projectRepo)

	for _, name := range []string{"Alice", "Bob", "Charlie"} {
		p := makePerson(proj.ID, name)
		require.NoError(t, personRepo.Create(ctx, p))
	}

	people, total, err := personRepo.List(ctx, person.PersonListFilter{
		ProjectID: proj.ID,
		Page:      1,
		PerPage:   10,
	})
	require.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, people, 3)
}
