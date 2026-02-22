package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/crypto"
	mock_auth "github.com/lbrty/observer/internal/domain/auth/mock"
	"github.com/lbrty/observer/internal/domain/user"
	mock_user "github.com/lbrty/observer/internal/domain/user/mock"
	"github.com/lbrty/observer/internal/ulid"
	ucauth "github.com/lbrty/observer/internal/usecase/auth"
)

func setupLoginUseCase(t *testing.T) (
	*ucauth.LoginUseCase,
	*mock_user.MockUserRepository,
	*mock_user.MockCredentialsRepository,
	*mock_auth.MockSessionRepository,
	*mock_user.MockMFARepository,
	crypto.PasswordHasher,
) {
	t.Helper()
	ctrl := gomock.NewController(t)

	mockUserRepo := mock_user.NewMockUserRepository(ctrl)
	mockCredRepo := mock_user.NewMockCredentialsRepository(ctrl)
	mockSessionRepo := mock_auth.NewMockSessionRepository(ctrl)
	mockMFARepo := mock_user.NewMockMFARepository(ctrl)
	hasher := crypto.NewArgonHasher()
	tokenGen := newTestTokenGen(t)

	uc := ucauth.NewLoginUseCase(
		mockUserRepo, mockCredRepo, mockSessionRepo, mockMFARepo, hasher, tokenGen,
	)
	return uc, mockUserRepo, mockCredRepo, mockSessionRepo, mockMFARepo, hasher
}

func newTestTokenGen(t *testing.T) crypto.TokenGenerator {
	t.Helper()
	tmpDir := t.TempDir()
	privPath, pubPath := generateTestKeys(t, tmpDir)
	keys, err := crypto.LoadRSAKeys(privPath, pubPath)
	require.NoError(t, err)
	return crypto.NewRSATokenGenerator(keys, 0, 0, 0, "test")
}

func TestLoginUseCase_Success(t *testing.T) {
	uc, mockUserRepo, mockCredRepo, mockSessionRepo, mockMFARepo, hasher := setupLoginUseCase(t)

	ctx := context.Background()
	password := "securepassword"
	hash, salt, err := hasher.Hash(password)
	require.NoError(t, err)

	uid := ulid.New()
	u := &user.User{ID: uid, Email: "test@example.com", Role: user.RoleConsultant, IsActive: true}
	cred := &user.Credentials{UserID: uid, PasswordHash: hash, Salt: salt}

	mockUserRepo.EXPECT().GetByEmail(ctx, u.Email).Return(u, nil)
	mockCredRepo.EXPECT().GetByUserID(ctx, uid).Return(cred, nil)
	mockMFARepo.EXPECT().GetByUserID(ctx, uid).Return(nil, errors.New("not found"))
	mockSessionRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	out, err := uc.Execute(ctx, ucauth.LoginInput{Email: u.Email, Password: password}, "agent", "1.2.3.4")
	require.NoError(t, err)
	assert.False(t, out.RequiresMFA)
	assert.NotNil(t, out.Tokens)
	assert.NotEmpty(t, out.Tokens.AccessToken)
	assert.NotEmpty(t, out.Tokens.RefreshToken)
}

func TestLoginUseCase_InvalidCredentials(t *testing.T) {
	uc, mockUserRepo, _, _, _, _ := setupLoginUseCase(t)

	mockUserRepo.EXPECT().
		GetByEmail(gomock.Any(), "bad@example.com").
		Return(nil, user.ErrUserNotFound)

	_, err := uc.Execute(context.Background(), ucauth.LoginInput{
		Email: "bad@example.com", Password: "pass",
	}, "", "")
	assert.ErrorIs(t, err, user.ErrInvalidCredentials)
}

func TestLoginUseCase_InactiveUser(t *testing.T) {
	uc, mockUserRepo, _, _, _, _ := setupLoginUseCase(t)

	uid := ulid.New()
	u := &user.User{ID: uid, Email: "inactive@example.com", IsActive: false}

	mockUserRepo.EXPECT().GetByEmail(gomock.Any(), u.Email).Return(u, nil)

	_, err := uc.Execute(context.Background(), ucauth.LoginInput{
		Email: u.Email, Password: "pass",
	}, "", "")
	assert.ErrorIs(t, err, user.ErrUserNotActive)
}
