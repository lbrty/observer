package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/crypto"
	"github.com/lbrty/observer/internal/domain/user"
	mock_repo "github.com/lbrty/observer/internal/repository/mock"
	ucauth "github.com/lbrty/observer/internal/usecase/auth"
)

func TestRegister_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	mockCredRepo := mock_repo.NewMockCredentialsRepository(ctrl)
	mockSessionRepo := mock_repo.NewMockSessionRepository(ctrl)
	mockMFARepo := mock_repo.NewMockMFARepository(ctrl)
	mockRecoveryRepo := mock_repo.NewMockMFARecoveryCodeRepository(ctrl)
	hasher := crypto.NewArgonHasher()
	tokenGen := newTestTokenGen(t)

	uc := ucauth.NewAuthUseCase(mockUserRepo, mockCredRepo, mockSessionRepo, mockMFARepo, mockRecoveryRepo, hasher, tokenGen, mock_repo.NewMockLoginAttemptStore(ctrl), 168*time.Hour)

	ctx := context.Background()
	input := ucauth.RegisterInput{
		Email:    "test@example.com",
		Password: "securepassword",
	}

	mockUserRepo.EXPECT().
		GetByEmail(ctx, input.Email).
		Return(nil, user.ErrUserNotFound)

	mockUserRepo.EXPECT().
		CreateWithCredentials(ctx, gomock.Any(), gomock.Any()).
		Return(nil)

	out, err := uc.Register(ctx, input)
	require.NoError(t, err)
	assert.NotEmpty(t, out.UserID)
	assert.Contains(t, out.Message, "Registration successful")
}

func TestRegister_AlwaysCreatesGuestRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	mockCredRepo := mock_repo.NewMockCredentialsRepository(ctrl)
	mockSessionRepo := mock_repo.NewMockSessionRepository(ctrl)
	mockMFARepo := mock_repo.NewMockMFARepository(ctrl)
	mockRecoveryRepo := mock_repo.NewMockMFARecoveryCodeRepository(ctrl)
	hasher := crypto.NewArgonHasher()
	tokenGen := newTestTokenGen(t)

	uc := ucauth.NewAuthUseCase(mockUserRepo, mockCredRepo, mockSessionRepo, mockMFARepo, mockRecoveryRepo, hasher, tokenGen, mock_repo.NewMockLoginAttemptStore(ctrl), 168*time.Hour)

	ctx := context.Background()

	mockUserRepo.EXPECT().
		GetByEmail(ctx, "attacker@example.com").
		Return(nil, user.ErrUserNotFound)

	mockUserRepo.EXPECT().
		CreateWithCredentials(ctx, gomock.AssignableToTypeOf(&user.User{}), gomock.Any()).
		DoAndReturn(func(_ context.Context, u *user.User, _ *user.Credentials) error {
			assert.Equal(t, user.RoleGuest, u.Role, "new users must always be created with RoleGuest")
			return nil
		})

	out, err := uc.Register(ctx, ucauth.RegisterInput{
		Email:    "attacker@example.com",
		Password: "Password1!",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, out.Message)
}

func TestRegister_DuplicateEmailReturnsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	mockCredRepo := mock_repo.NewMockCredentialsRepository(ctrl)
	mockSessionRepo := mock_repo.NewMockSessionRepository(ctrl)
	mockMFARepo := mock_repo.NewMockMFARepository(ctrl)
	mockRecoveryRepo := mock_repo.NewMockMFARecoveryCodeRepository(ctrl)
	hasher := crypto.NewArgonHasher()
	tokenGen := newTestTokenGen(t)

	uc := ucauth.NewAuthUseCase(mockUserRepo, mockCredRepo, mockSessionRepo, mockMFARepo, mockRecoveryRepo, hasher, tokenGen, mock_repo.NewMockLoginAttemptStore(ctrl), 168*time.Hour)

	// GetByEmail returns an existing user (no error) — simulates duplicate email.
	mockUserRepo.EXPECT().
		GetByEmail(gomock.Any(), "existing@example.com").
		Return(&user.User{}, nil)

	out, err := uc.Register(context.Background(), ucauth.RegisterInput{
		Email:    "existing@example.com",
		Password: "securepassword",
	})
	require.NoError(t, err, "duplicate email must NOT return an error to the caller")
	assert.NotEmpty(t, out.Message)
}
