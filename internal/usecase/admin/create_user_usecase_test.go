package admin_test

import (
	"context"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	cryptomock "github.com/lbrty/observer/internal/crypto/mock"
	"github.com/lbrty/observer/internal/domain/user"
	mock_repo "github.com/lbrty/observer/internal/repository/mock"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
)

func TestUserUseCase_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	mockCredRepo := mock_repo.NewMockCredentialsRepository(ctrl)
	mockHasher := cryptomock.NewMockPasswordHasher(ctrl)
	uc := ucadmin.NewUserUseCase(mockUserRepo, mockCredRepo, mockHasher)

	ctx := context.Background()

	// Email not taken
	mockUserRepo.EXPECT().GetByEmail(ctx, "alice@example.com").Return(nil, user.ErrUserNotFound)
	// No phone to check
	mockHasher.EXPECT().Hash("securepass").Return("hashed", "salt", nil)
	mockUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
	mockCredRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	out, err := uc.Create(ctx, ucadmin.CreateUserInput{
		FirstName: "Alice",
		LastName:  "Smith",
		Email:     "alice@example.com",
		Password:  "securepass",
		Role:      "admin",
		IsActive:  true,
	})
	require.NoError(t, err)
	assert.Equal(t, "Alice", out.FirstName)
	assert.Equal(t, "Smith", out.LastName)
	assert.Equal(t, "alice@example.com", out.Email)
	assert.Equal(t, "admin", out.Role)
	assert.True(t, out.IsActive)
	assert.NotEmpty(t, out.ID)
}

func TestUserUseCase_Create_InvalidRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	mockCredRepo := mock_repo.NewMockCredentialsRepository(ctrl)
	mockHasher := cryptomock.NewMockPasswordHasher(ctrl)
	uc := ucadmin.NewUserUseCase(mockUserRepo, mockCredRepo, mockHasher)

	_, err := uc.Create(context.Background(), ucadmin.CreateUserInput{
		FirstName: "Bob",
		Email:     "bob@example.com",
		Password:  "securepass",
		Role:      "superadmin",
	})
	assert.ErrorIs(t, err, user.ErrInvalidRole)
}

func TestUserUseCase_Create_DuplicateEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	mockCredRepo := mock_repo.NewMockCredentialsRepository(ctrl)
	mockHasher := cryptomock.NewMockPasswordHasher(ctrl)
	uc := ucadmin.NewUserUseCase(mockUserRepo, mockCredRepo, mockHasher)

	ctx := context.Background()

	// Email already taken (GetByEmail returns a user, no error)
	mockUserRepo.EXPECT().GetByEmail(ctx, "alice@example.com").Return(&user.User{
		ID: ulid.Make(), Email: "alice@example.com",
	}, nil)

	_, err := uc.Create(ctx, ucadmin.CreateUserInput{
		FirstName: "Alice",
		Email:     "alice@example.com",
		Password:  "securepass",
		Role:      "staff",
	})
	assert.ErrorIs(t, err, user.ErrEmailExists)
}

func TestUserUseCase_Create_DuplicatePhone(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	mockCredRepo := mock_repo.NewMockCredentialsRepository(ctrl)
	mockHasher := cryptomock.NewMockPasswordHasher(ctrl)
	uc := ucadmin.NewUserUseCase(mockUserRepo, mockCredRepo, mockHasher)

	ctx := context.Background()

	// Email not taken
	mockUserRepo.EXPECT().GetByEmail(ctx, "bob@example.com").Return(nil, user.ErrUserNotFound)
	// Phone already taken
	mockUserRepo.EXPECT().GetByPhone(ctx, "+380991234567").Return(&user.User{
		ID: ulid.Make(), Phone: "+380991234567",
	}, nil)

	_, err := uc.Create(ctx, ucadmin.CreateUserInput{
		FirstName: "Bob",
		Email:     "bob@example.com",
		Phone:     "+380991234567",
		Password:  "securepass",
		Role:      "staff",
	})
	assert.ErrorIs(t, err, user.ErrPhoneExists)
}

func TestUserUseCase_ResetPassword_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	mockCredRepo := mock_repo.NewMockCredentialsRepository(ctrl)
	mockHasher := cryptomock.NewMockPasswordHasher(ctrl)
	uc := ucadmin.NewUserUseCase(mockUserRepo, mockCredRepo, mockHasher)

	ctx := context.Background()
	userID := ulid.Make()

	mockCredRepo.EXPECT().GetByUserID(ctx, userID).Return(&user.Credentials{
		UserID:       userID,
		PasswordHash: "oldhash",
		Salt:         "oldsalt",
	}, nil)
	mockHasher.EXPECT().Hash("newpassword123").Return("newhash", "newsalt", nil)
	mockCredRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, c *user.Credentials) error {
		assert.Equal(t, "newhash", c.PasswordHash)
		assert.Equal(t, "newsalt", c.Salt)
		return nil
	})

	err := uc.ResetPassword(ctx, userID, ucadmin.ResetPasswordInput{
		NewPassword: "newpassword123",
	})
	require.NoError(t, err)
}

func TestUserUseCase_ResetPassword_CredNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	mockCredRepo := mock_repo.NewMockCredentialsRepository(ctrl)
	mockHasher := cryptomock.NewMockPasswordHasher(ctrl)
	uc := ucadmin.NewUserUseCase(mockUserRepo, mockCredRepo, mockHasher)

	ctx := context.Background()
	userID := ulid.Make()

	mockCredRepo.EXPECT().GetByUserID(ctx, userID).Return(nil, user.ErrUserNotFound)

	err := uc.ResetPassword(ctx, userID, ucadmin.ResetPasswordInput{
		NewPassword: "newpassword123",
	})
	assert.ErrorIs(t, err, user.ErrUserNotFound)
}
