package admin_test

import (
	"context"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/user"
	mock_repo "github.com/lbrty/observer/internal/repository/mock"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
)

func ptr[T any](v T) *T { return &v }

func TestUpdateUserUseCase_PartialUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockUserRepository(ctrl)
	uc := ucadmin.NewUpdateUserUseCase(mockRepo)

	ctx := context.Background()
	uid := ulid.MustNew(ulid.Now(), nil)

	existing := &user.User{
		ID:         uid,
		FirstName:  "Alice",
		LastName:   "Smith",
		Email:      "alice@example.com",
		Phone:      "+49555000111",
		Role:       user.RoleStaff,
		IsActive:   true,
		IsVerified: true,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	mockRepo.EXPECT().GetByID(ctx, uid).Return(existing, nil)
	mockRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, u *user.User) error {
		assert.Equal(t, "Bob", u.FirstName)
		assert.Equal(t, "Smith", u.LastName)          // unchanged
		assert.Equal(t, "alice@example.com", u.Email) // unchanged
		return nil
	})

	input := ucadmin.UpdateUserInput{
		FirstName: ptr("Bob"),
	}

	out, err := uc.Execute(ctx, uid, input)
	require.NoError(t, err)
	assert.Equal(t, "Bob", out.FirstName)
	assert.Equal(t, "Smith", out.LastName)
}

func TestUpdateUserUseCase_RoleChange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockUserRepository(ctrl)
	uc := ucadmin.NewUpdateUserUseCase(mockRepo)

	ctx := context.Background()
	uid := ulid.MustNew(ulid.Now(), nil)

	existing := &user.User{
		ID:        uid,
		Email:     "user@example.com",
		Phone:     "+49555000222",
		Role:      user.RoleGuest,
		IsActive:  true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	mockRepo.EXPECT().GetByID(ctx, uid).Return(existing, nil)
	mockRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)

	out, err := uc.Execute(ctx, uid, ucadmin.UpdateUserInput{
		Role: ptr("admin"),
	})
	require.NoError(t, err)
	assert.Equal(t, "admin", out.Role)
}

func TestUpdateUserUseCase_InvalidRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockUserRepository(ctrl)
	uc := ucadmin.NewUpdateUserUseCase(mockRepo)

	ctx := context.Background()
	uid := ulid.MustNew(ulid.Now(), nil)

	existing := &user.User{
		ID:        uid,
		Email:     "user@example.com",
		Phone:     "+49555000333",
		Role:      user.RoleStaff,
		IsActive:  true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	mockRepo.EXPECT().GetByID(ctx, uid).Return(existing, nil)

	_, err := uc.Execute(ctx, uid, ucadmin.UpdateUserInput{
		Role: ptr("superadmin"),
	})
	assert.ErrorIs(t, err, user.ErrInvalidRole)
}

func TestUpdateUserUseCase_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockUserRepository(ctrl)
	uc := ucadmin.NewUpdateUserUseCase(mockRepo)

	uid := ulid.MustNew(ulid.Now(), nil)
	mockRepo.EXPECT().GetByID(gomock.Any(), uid).Return(nil, user.ErrUserNotFound)

	_, err := uc.Execute(context.Background(), uid, ucadmin.UpdateUserInput{
		FirstName: ptr("New"),
	})
	assert.ErrorIs(t, err, user.ErrUserNotFound)
}

func TestUpdateUserUseCase_Deactivate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockUserRepository(ctrl)
	uc := ucadmin.NewUpdateUserUseCase(mockRepo)

	ctx := context.Background()
	uid := ulid.MustNew(ulid.Now(), nil)

	existing := &user.User{
		ID:        uid,
		Email:     "user@example.com",
		Phone:     "+49555000444",
		Role:      user.RoleStaff,
		IsActive:  true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	mockRepo.EXPECT().GetByID(ctx, uid).Return(existing, nil)
	mockRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, u *user.User) error {
		assert.False(t, u.IsActive)
		return nil
	})

	out, err := uc.Execute(ctx, uid, ucadmin.UpdateUserInput{
		IsActive: ptr(false),
	})
	require.NoError(t, err)
	assert.False(t, out.IsActive)
}
