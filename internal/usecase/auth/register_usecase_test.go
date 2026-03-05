package auth_test

import (
	"context"
	"testing"

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
	hasher := crypto.NewArgonHasher()
	tokenGen := newTestTokenGen(t)

	uc := ucauth.NewAuthUseCase(mockUserRepo, mockCredRepo, mockSessionRepo, mockMFARepo, hasher, tokenGen)

	ctx := context.Background()
	input := ucauth.RegisterInput{
		Email:    "test@example.com",
		Password: "securepassword",
		Role:     "consultant",
	}

	mockUserRepo.EXPECT().
		GetByEmail(ctx, input.Email).
		Return(nil, user.ErrUserNotFound)

	mockUserRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(nil)

	mockCredRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(nil)

	out, err := uc.Register(ctx, input)
	require.NoError(t, err)
	assert.NotEmpty(t, out.UserID)
	assert.Contains(t, out.Message, "Registration successful")
}

func TestRegister_EmailExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	mockCredRepo := mock_repo.NewMockCredentialsRepository(ctrl)
	mockSessionRepo := mock_repo.NewMockSessionRepository(ctrl)
	mockMFARepo := mock_repo.NewMockMFARepository(ctrl)
	hasher := crypto.NewArgonHasher()
	tokenGen := newTestTokenGen(t)

	uc := ucauth.NewAuthUseCase(mockUserRepo, mockCredRepo, mockSessionRepo, mockMFARepo, hasher, tokenGen)

	mockUserRepo.EXPECT().
		GetByEmail(gomock.Any(), "taken@example.com").
		Return(&user.User{}, nil)

	_, err := uc.Register(context.Background(), ucauth.RegisterInput{
		Email:    "taken@example.com",
		Password: "securepassword",
		Role:     "consultant",
	})
	assert.ErrorIs(t, err, user.ErrEmailExists)
}

func TestRegister_InvalidRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	mockCredRepo := mock_repo.NewMockCredentialsRepository(ctrl)
	mockSessionRepo := mock_repo.NewMockSessionRepository(ctrl)
	mockMFARepo := mock_repo.NewMockMFARepository(ctrl)
	hasher := crypto.NewArgonHasher()
	tokenGen := newTestTokenGen(t)

	uc := ucauth.NewAuthUseCase(mockUserRepo, mockCredRepo, mockSessionRepo, mockMFARepo, hasher, tokenGen)

	_, err := uc.Register(context.Background(), ucauth.RegisterInput{
		Email:    "test@example.com",
		Password: "securepassword",
		Role:     "superadmin",
	})
	assert.ErrorIs(t, err, user.ErrInvalidRole)
}
