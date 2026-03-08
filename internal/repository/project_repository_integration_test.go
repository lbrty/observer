//go:build !short

package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/repository"
	iulid "github.com/lbrty/observer/internal/ulid"
)

func makeProject(name, ownerID string) *project.Project {
	return &project.Project{
		ID:      iulid.NewString(),
		Name:    name,
		OwnerID: ownerID,
		Status:  project.ProjectStatusActive,
	}
}

func createOwnerUser(t *testing.T, ctx context.Context, userRepo repository.UserRepository) *user.User {
	t.Helper()
	now := time.Now().UTC()
	u := &user.User{
		ID:         ulid.Make(),
		FirstName:  "Owner",
		LastName:   "User",
		Email:      "owner-" + iulid.NewString() + "@test.com",
		Phone:      "+" + iulid.NewString()[:11],
		Role:       user.RoleAdmin,
		IsVerified: true,
		IsActive:   true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	require.NoError(t, userRepo.Create(ctx, u))
	return u
}

func TestProjectRepo_CRUD(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	ctx := context.Background()

	owner := createOwnerUser(t, ctx, userRepo)

	p := makeProject("Test Project", owner.ID.String())
	require.NoError(t, projectRepo.Create(ctx, p))

	// Get
	got, err := projectRepo.GetByID(ctx, p.ID)
	require.NoError(t, err)
	assert.Equal(t, "Test Project", got.Name)
	assert.Equal(t, owner.ID.String(), got.OwnerID)
	assert.Equal(t, project.ProjectStatusActive, got.Status)

	// Update
	p.Name = "Updated Project"
	p.Status = project.ProjectStatusArchived
	require.NoError(t, projectRepo.Update(ctx, p))
	got, err = projectRepo.GetByID(ctx, p.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Project", got.Name)
	assert.Equal(t, project.ProjectStatusArchived, got.Status)

	// List
	list, total, err := projectRepo.List(ctx, project.ProjectListFilter{Page: 1, PerPage: 10})
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, list, 1)
}

func TestProjectRepo_ListWithFilter(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	ctx := context.Background()

	owner := createOwnerUser(t, ctx, userRepo)

	p1 := makeProject("Active Project", owner.ID.String())
	p1.Status = project.ProjectStatusActive
	require.NoError(t, projectRepo.Create(ctx, p1))

	p2 := makeProject("Archived Project", owner.ID.String())
	p2.Status = project.ProjectStatusArchived
	require.NoError(t, projectRepo.Create(ctx, p2))

	// Filter by status
	active := project.ProjectStatusActive
	list, total, err := projectRepo.List(ctx, project.ProjectListFilter{
		Status:  &active,
		Page:    1,
		PerPage: 10,
	})
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, list, 1)
	assert.Equal(t, "Active Project", list[0].Name)

	// Filter by owner
	ownerID := owner.ID.String()
	list, total, err = projectRepo.List(ctx, project.ProjectListFilter{
		OwnerID: &ownerID,
		Page:    1,
		PerPage: 10,
	})
	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, list, 2)
}

func TestPermissionRepo_CRUD(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	permRepo := repository.NewProjectPermissionRepository(db)
	ctx := context.Background()

	owner := createOwnerUser(t, ctx, userRepo)
	proj := makeProject("Perm Project", owner.ID.String())
	require.NoError(t, projectRepo.Create(ctx, proj))

	// Create a second user for the permission
	staffUser := createOwnerUser(t, ctx, userRepo)

	perm := &project.ProjectPermission{
		ID:               iulid.NewString(),
		ProjectID:        proj.ID,
		UserID:           staffUser.ID.String(),
		Role:             project.ProjectRoleManager,
		CanViewContact:   true,
		CanViewPersonal:  false,
		CanViewDocuments: true,
	}

	// Create
	require.NoError(t, permRepo.Create(ctx, perm))

	// Get
	got, err := permRepo.GetByID(ctx, perm.ID)
	require.NoError(t, err)
	assert.Equal(t, proj.ID, got.ProjectID)
	assert.Equal(t, staffUser.ID.String(), got.UserID)
	assert.Equal(t, project.ProjectRoleManager, got.Role)
	assert.True(t, got.CanViewContact)
	assert.False(t, got.CanViewPersonal)
	assert.True(t, got.CanViewDocuments)

	// List
	list, err := permRepo.List(ctx, proj.ID)
	require.NoError(t, err)
	assert.Len(t, list, 1)

	// Update
	perm.Role = project.ProjectRoleViewer
	perm.CanViewPersonal = true
	require.NoError(t, permRepo.Update(ctx, perm))
	got, err = permRepo.GetByID(ctx, perm.ID)
	require.NoError(t, err)
	assert.Equal(t, project.ProjectRoleViewer, got.Role)
	assert.True(t, got.CanViewPersonal)

	// Delete
	require.NoError(t, permRepo.Delete(ctx, perm.ID))
	_, err = permRepo.GetByID(ctx, perm.ID)
	assert.ErrorIs(t, err, project.ErrPermissionNotFound)
}
