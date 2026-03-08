//go:build !short

package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lbrty/observer/internal/domain/auth"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/repository"
)

func TestCredentialsRepo_CreateAndGet(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	credRepo := repository.NewCredentialsRepository(db)
	ctx := context.Background()

	u := makeUser("cred-create@test.com", "+12000000001")
	require.NoError(t, userRepo.Create(ctx, u))

	cred := &user.Credentials{
		UserID:       u.ID,
		PasswordHash: "hashed_password_123",
		Salt:         "salt_abc",
		UpdatedAt:    time.Now().UTC(),
	}
	require.NoError(t, credRepo.Create(ctx, cred))

	got, err := credRepo.GetByUserID(ctx, u.ID)
	require.NoError(t, err)
	assert.Equal(t, u.ID, got.UserID)
	assert.Equal(t, "hashed_password_123", got.PasswordHash)
	assert.Equal(t, "salt_abc", got.Salt)
}

func TestCredentialsRepo_Update(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	credRepo := repository.NewCredentialsRepository(db)
	ctx := context.Background()

	u := makeUser("cred-update@test.com", "+12000000002")
	require.NoError(t, userRepo.Create(ctx, u))

	cred := &user.Credentials{
		UserID:       u.ID,
		PasswordHash: "old_hash",
		Salt:         "old_salt",
		UpdatedAt:    time.Now().UTC(),
	}
	require.NoError(t, credRepo.Create(ctx, cred))

	cred.PasswordHash = "new_hash"
	cred.Salt = "new_salt"
	require.NoError(t, credRepo.Update(ctx, cred))

	got, err := credRepo.GetByUserID(ctx, u.ID)
	require.NoError(t, err)
	assert.Equal(t, "new_hash", got.PasswordHash)
	assert.Equal(t, "new_salt", got.Salt)
}

func TestSessionRepo_CreateAndGet(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	ctx := context.Background()

	u := makeUser("session-create@test.com", "+12000000003")
	require.NoError(t, userRepo.Create(ctx, u))

	sess := &auth.Session{
		ID:           ulid.Make(),
		UserID:       u.ID,
		RefreshToken: "refresh_token_abc",
		UserAgent:    "TestAgent/1.0",
		IP:           "127.0.0.1",
		ExpiresAt:    time.Now().Add(24 * time.Hour).UTC(),
		CreatedAt:    time.Now().UTC(),
	}
	require.NoError(t, sessionRepo.Create(ctx, sess))

	got, err := sessionRepo.GetByRefreshToken(ctx, "refresh_token_abc")
	require.NoError(t, err)
	assert.Equal(t, sess.ID, got.ID)
	assert.Equal(t, sess.UserID, got.UserID)
	assert.Equal(t, "refresh_token_abc", got.RefreshToken)
	assert.Equal(t, "TestAgent/1.0", got.UserAgent)
	assert.Equal(t, "127.0.0.1", got.IP)
}

func TestSessionRepo_Delete(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	ctx := context.Background()

	u := makeUser("session-delete@test.com", "+12000000004")
	require.NoError(t, userRepo.Create(ctx, u))

	sess := &auth.Session{
		ID:           ulid.Make(),
		UserID:       u.ID,
		RefreshToken: "refresh_token_del",
		UserAgent:    "TestAgent/1.0",
		IP:           "127.0.0.1",
		ExpiresAt:    time.Now().Add(24 * time.Hour).UTC(),
		CreatedAt:    time.Now().UTC(),
	}
	require.NoError(t, sessionRepo.Create(ctx, sess))

	require.NoError(t, sessionRepo.Delete(ctx, sess.ID))

	_, err := sessionRepo.GetByRefreshToken(ctx, "refresh_token_del")
	assert.ErrorIs(t, err, auth.ErrSessionNotFound)
}
