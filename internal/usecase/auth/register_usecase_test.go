package auth_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/crypto"
	"github.com/lbrty/observer/internal/domain/user"
	mock_user "github.com/lbrty/observer/internal/domain/user/mock"
	ucauth "github.com/lbrty/observer/internal/usecase/auth"
)

func TestRegisterUseCase_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_user.NewMockUserRepository(ctrl)
	mockCredRepo := mock_user.NewMockCredentialsRepository(ctrl)
	hasher := crypto.NewArgonHasher()

	uc := ucauth.NewRegisterUseCase(mockUserRepo, mockCredRepo, hasher)

	ctx := context.Background()
	input := ucauth.RegisterInput{
		Email:    "test@example.com",
		Phone:    "+49555000111",
		Password: "securepassword",
		Role:     "consultant",
	}

	mockUserRepo.EXPECT().
		GetByEmail(ctx, input.Email).
		Return(nil, user.ErrUserNotFound)

	mockUserRepo.EXPECT().
		GetByPhone(ctx, input.Phone).
		Return(nil, user.ErrUserNotFound)

	mockUserRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(nil)

	mockCredRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(nil)

	out, err := uc.Execute(ctx, input)
	require.NoError(t, err)
	assert.NotEmpty(t, out.UserID)
	assert.Contains(t, out.Message, "Registration successful")
}

func TestRegisterUseCase_EmailExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_user.NewMockUserRepository(ctrl)
	mockCredRepo := mock_user.NewMockCredentialsRepository(ctrl)
	hasher := crypto.NewArgonHasher()

	uc := ucauth.NewRegisterUseCase(mockUserRepo, mockCredRepo, hasher)

	mockUserRepo.EXPECT().
		GetByEmail(gomock.Any(), "taken@example.com").
		Return(&user.User{}, nil)

	_, err := uc.Execute(context.Background(), ucauth.RegisterInput{
		Email:    "taken@example.com",
		Phone:    "+49555000222",
		Password: "securepassword",
		Role:     "consultant",
	})
	assert.ErrorIs(t, err, user.ErrEmailExists)
}

func TestRegisterUseCase_PhoneExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_user.NewMockUserRepository(ctrl)
	mockCredRepo := mock_user.NewMockCredentialsRepository(ctrl)
	hasher := crypto.NewArgonHasher()

	uc := ucauth.NewRegisterUseCase(mockUserRepo, mockCredRepo, hasher)

	mockUserRepo.EXPECT().
		GetByEmail(gomock.Any(), gomock.Any()).
		Return(nil, user.ErrUserNotFound)

	mockUserRepo.EXPECT().
		GetByPhone(gomock.Any(), "+49555000333").
		Return(&user.User{}, nil)

	_, err := uc.Execute(context.Background(), ucauth.RegisterInput{
		Email:    "free@example.com",
		Phone:    "+49555000333",
		Password: "securepassword",
		Role:     "consultant",
	})
	assert.ErrorIs(t, err, user.ErrPhoneExists)
}

func TestRegisterUseCase_InvalidRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_user.NewMockUserRepository(ctrl)
	mockCredRepo := mock_user.NewMockCredentialsRepository(ctrl)
	hasher := crypto.NewArgonHasher()

	uc := ucauth.NewRegisterUseCase(mockUserRepo, mockCredRepo, hasher)

	_, err := uc.Execute(context.Background(), ucauth.RegisterInput{
		Email:    "test@example.com",
		Phone:    "+49555000444",
		Password: "securepassword",
		Role:     "superadmin",
	})
	assert.ErrorIs(t, err, user.ErrInvalidRole)
}
