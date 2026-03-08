//go:build !short

package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/repository"
)

func makeUser(email, phone string) *user.User {
	now := time.Now().UTC()
	return &user.User{
		ID:         ulid.Make(),
		FirstName:  "Test",
		LastName:   "User",
		Email:      email,
		Phone:      phone,
		Role:       user.RoleStaff,
		IsVerified: false,
		IsActive:   true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func TestUserRepo_CreateAndGetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	u := makeUser("create-get@test.com", "+10000000001")
	require.NoError(t, repo.Create(ctx, u))

	got, err := repo.GetByID(ctx, u.ID)
	require.NoError(t, err)
	assert.Equal(t, u.ID, got.ID)
	assert.Equal(t, u.FirstName, got.FirstName)
	assert.Equal(t, u.LastName, got.LastName)
	assert.Equal(t, u.Email, got.Email)
	assert.Equal(t, u.Phone, got.Phone)
	assert.Equal(t, u.Role, got.Role)
	assert.Equal(t, u.IsVerified, got.IsVerified)
	assert.Equal(t, u.IsActive, got.IsActive)
}

func TestUserRepo_GetByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	u := makeUser("byemail@test.com", "+10000000002")
	require.NoError(t, repo.Create(ctx, u))

	got, err := repo.GetByEmail(ctx, "byemail@test.com")
	require.NoError(t, err)
	assert.Equal(t, u.ID, got.ID)
	assert.Equal(t, u.Email, got.Email)
}

func TestUserRepo_GetByEmail_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	_, err := repo.GetByEmail(ctx, "nonexistent@test.com")
	assert.ErrorIs(t, err, user.ErrUserNotFound)
}

func TestUserRepo_DuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	u1 := makeUser("dup@test.com", "+10000000003")
	require.NoError(t, repo.Create(ctx, u1))

	u2 := makeUser("dup@test.com", "+10000000004")
	err := repo.Create(ctx, u2)
	assert.Error(t, err)
}

func TestUserRepo_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	u := makeUser("update@test.com", "+10000000005")
	require.NoError(t, repo.Create(ctx, u))

	u.FirstName = "Updated"
	u.LastName = "Name"
	u.IsVerified = true
	require.NoError(t, repo.Update(ctx, u))

	got, err := repo.GetByID(ctx, u.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated", got.FirstName)
	assert.Equal(t, "Name", got.LastName)
	assert.True(t, got.IsVerified)
}

func TestUserRepo_List(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	for i, email := range []string{"list1@test.com", "list2@test.com", "list3@test.com"} {
		u := makeUser(email, "+1000000010"+string(rune('0'+i)))
		require.NoError(t, repo.Create(ctx, u))
	}

	users, total, err := repo.List(ctx, user.UserListFilter{Page: 1, PerPage: 10})
	require.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, users, 3)
}
