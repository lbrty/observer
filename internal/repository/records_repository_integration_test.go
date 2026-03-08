//go:build !short

package repository_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lbrty/observer/internal/domain/household"
	"github.com/lbrty/observer/internal/domain/person"
	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/domain/support"
	"github.com/lbrty/observer/internal/domain/tag"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/repository"
	iulid "github.com/lbrty/observer/internal/ulid"
)

// setupProjectEnv creates a user, project, and returns them along with the DB repos.
func setupProjectEnv(t *testing.T, db *sqlx.DB) (repository.UserRepository, repository.ProjectRepository, *user.User, *project.Project) {
	t.Helper()
	ctx := context.Background()
	now := time.Now().UTC()

	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)

	owner := &user.User{
		ID:         ulid.Make(),
		FirstName:  "Env",
		LastName:   "Owner",
		Email:      "env-" + iulid.NewString() + "@test.com",
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
		Name:    "Env Project " + iulid.NewString(),
		OwnerID: owner.ID.String(),
		Status:  project.ProjectStatusActive,
	}
	require.NoError(t, projectRepo.Create(ctx, proj))

	return userRepo, projectRepo, owner, proj
}

func TestTagRepo_CRUD(t *testing.T) {
	db := setupTestDB(t)
	_, _, _, proj := setupProjectEnv(t, db)
	tagRepo := repository.NewTagRepository(db)
	ctx := context.Background()

	tg := &tag.Tag{
		ID:        iulid.NewString(),
		ProjectID: proj.ID,
		Name:      "Urgent",
		Color:     "#ff0000",
	}

	// Create
	require.NoError(t, tagRepo.Create(ctx, tg))

	// Get
	got, err := tagRepo.GetByID(ctx, tg.ID)
	require.NoError(t, err)
	assert.Equal(t, "Urgent", got.Name)
	assert.Equal(t, "#ff0000", got.Color)
	assert.Equal(t, proj.ID, got.ProjectID)

	// Update
	tg.Name = "High Priority"
	tg.Color = "#ff8800"
	require.NoError(t, tagRepo.Update(ctx, tg))
	got, err = tagRepo.GetByID(ctx, tg.ID)
	require.NoError(t, err)
	assert.Equal(t, "High Priority", got.Name)
	assert.Equal(t, "#ff8800", got.Color)

	// List
	list, err := tagRepo.List(ctx, proj.ID)
	require.NoError(t, err)
	assert.Len(t, list, 1)

	// Delete
	require.NoError(t, tagRepo.Delete(ctx, tg.ID))
	_, err = tagRepo.GetByID(ctx, tg.ID)
	assert.ErrorIs(t, err, tag.ErrTagNotFound)
}

func TestTagRepo_DuplicateName(t *testing.T) {
	db := setupTestDB(t)
	_, _, _, proj := setupProjectEnv(t, db)
	tagRepo := repository.NewTagRepository(db)
	ctx := context.Background()

	t1 := &tag.Tag{ID: iulid.NewString(), ProjectID: proj.ID, Name: "DupTag", Color: "#000"}
	require.NoError(t, tagRepo.Create(ctx, t1))

	t2 := &tag.Tag{ID: iulid.NewString(), ProjectID: proj.ID, Name: "DupTag", Color: "#111"}
	err := tagRepo.Create(ctx, t2)
	assert.ErrorIs(t, err, tag.ErrTagNameExists)
}

func TestHouseholdRepo_CRUD(t *testing.T) {
	db := setupTestDB(t)
	_, _, _, proj := setupProjectEnv(t, db)
	hhRepo := repository.NewHouseholdRepository(db)
	ctx := context.Background()

	hh := &household.Household{
		ID:        iulid.NewString(),
		ProjectID: proj.ID,
	}

	// Create
	require.NoError(t, hhRepo.Create(ctx, hh))

	// Get
	got, err := hhRepo.GetByID(ctx, hh.ID)
	require.NoError(t, err)
	assert.Equal(t, hh.ID, got.ID)
	assert.Equal(t, proj.ID, got.ProjectID)

	// List
	list, total, err := hhRepo.List(ctx, proj.ID, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, list, 1)

	// Delete
	require.NoError(t, hhRepo.Delete(ctx, hh.ID))
	_, err = hhRepo.GetByID(ctx, hh.ID)
	assert.ErrorIs(t, err, household.ErrHouseholdNotFound)
}

func TestSupportRecordRepo_CRUD(t *testing.T) {
	db := setupTestDB(t)
	_, _, _, proj := setupProjectEnv(t, db)
	personRepo := repository.NewPersonRepository(db)
	supportRepo := repository.NewSupportRecordRepository(db)
	ctx := context.Background()

	p := &person.Person{
		ID:           iulid.NewString(),
		ProjectID:    proj.ID,
		FirstName:    "Support Test",
		Sex:          person.SexFemale,
		CaseStatus:   person.CaseStatusNew,
		PhoneNumbers: json.RawMessage("[]"),
	}
	require.NoError(t, personRepo.Create(ctx, p))

	sphere := support.SphereHousingAssistance
	rec := &support.Record{
		ID:        iulid.NewString(),
		PersonID:  p.ID,
		ProjectID: proj.ID,
		Type:      support.SupportTypeLegal,
		Sphere:    &sphere,
	}

	// Create
	require.NoError(t, supportRepo.Create(ctx, rec))

	// Get
	got, err := supportRepo.GetByID(ctx, rec.ID)
	require.NoError(t, err)
	assert.Equal(t, rec.ID, got.ID)
	assert.Equal(t, p.ID, got.PersonID)
	assert.Equal(t, support.SupportTypeLegal, got.Type)
	require.NotNil(t, got.Sphere)
	assert.Equal(t, support.SphereHousingAssistance, *got.Sphere)

	// List
	list, total, err := supportRepo.List(ctx, support.RecordListFilter{
		ProjectID: proj.ID,
		Page:      1,
		PerPage:   10,
	})
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, list, 1)

	// Delete
	require.NoError(t, supportRepo.Delete(ctx, rec.ID))
	_, err = supportRepo.GetByID(ctx, rec.ID)
	assert.ErrorIs(t, err, support.ErrRecordNotFound)
}
