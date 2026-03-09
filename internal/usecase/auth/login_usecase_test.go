package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/crypto"
	domainauth "github.com/lbrty/observer/internal/domain/auth"
	"github.com/lbrty/observer/internal/domain/user"
	mock_repo "github.com/lbrty/observer/internal/repository/mock"
	"github.com/lbrty/observer/internal/ulid"
	ucauth "github.com/lbrty/observer/internal/usecase/auth"
)

func setupAuthUseCase(t *testing.T) (
	*ucauth.AuthUseCase,
	*mock_repo.MockUserRepository,
	*mock_repo.MockCredentialsRepository,
	*mock_repo.MockSessionRepository,
	*mock_repo.MockMFARepository,
	crypto.PasswordHasher,
) {
	t.Helper()
	ctrl := gomock.NewController(t)

	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	mockCredRepo := mock_repo.NewMockCredentialsRepository(ctrl)
	mockSessionRepo := mock_repo.NewMockSessionRepository(ctrl)
	mockMFARepo := mock_repo.NewMockMFARepository(ctrl)
	hasher := crypto.NewArgonHasher()
	tokenGen := newTestTokenGen(t)

	uc := ucauth.NewAuthUseCase(
		mockUserRepo, mockCredRepo, mockSessionRepo, mockMFARepo, hasher, tokenGen,
	)
	return uc, mockUserRepo, mockCredRepo, mockSessionRepo, mockMFARepo, hasher
}

// setupAuthUseCaseWithTokenGen is like setupAuthUseCase but also returns the TokenGenerator
// so tests can produce valid MFA tokens for use case calls.
func setupAuthUseCaseWithTokenGen(t *testing.T) (
	*ucauth.AuthUseCase,
	*mock_repo.MockUserRepository,
	*mock_repo.MockCredentialsRepository,
	*mock_repo.MockSessionRepository,
	*mock_repo.MockMFARepository,
	crypto.PasswordHasher,
	crypto.TokenGenerator,
) {
	t.Helper()
	ctrl := gomock.NewController(t)

	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	mockCredRepo := mock_repo.NewMockCredentialsRepository(ctrl)
	mockSessionRepo := mock_repo.NewMockSessionRepository(ctrl)
	mockMFARepo := mock_repo.NewMockMFARepository(ctrl)
	hasher := crypto.NewArgonHasher()
	tokenGen := newTestTokenGen(t)

	uc := ucauth.NewAuthUseCase(
		mockUserRepo, mockCredRepo, mockSessionRepo, mockMFARepo, hasher, tokenGen,
	)
	return uc, mockUserRepo, mockCredRepo, mockSessionRepo, mockMFARepo, hasher, tokenGen
}

func newTestTokenGen(t *testing.T) crypto.TokenGenerator {
	t.Helper()
	tmpDir := t.TempDir()
	privPath, pubPath := generateTestKeys(t, tmpDir)
	keys, err := crypto.LoadRSAKeys(privPath, pubPath)
	require.NoError(t, err)
	return crypto.NewRSATokenGenerator(keys, 0, 0, 5*time.Minute, "test")
}

func TestLogin_Success(t *testing.T) {
	uc, mockUserRepo, mockCredRepo, mockSessionRepo, mockMFARepo, hasher := setupAuthUseCase(t)

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

	out, err := uc.Login(ctx, ucauth.LoginInput{Email: u.Email, Password: password}, "agent", "1.2.3.4")
	require.NoError(t, err)
	assert.False(t, out.RequiresMFA)
	assert.NotNil(t, out.Tokens)
	assert.NotEmpty(t, out.Tokens.AccessToken)
	assert.NotEmpty(t, out.Tokens.RefreshToken)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	uc, mockUserRepo, _, _, _, _ := setupAuthUseCase(t)

	mockUserRepo.EXPECT().
		GetByEmail(gomock.Any(), "bad@example.com").
		Return(nil, user.ErrUserNotFound)

	_, err := uc.Login(context.Background(), ucauth.LoginInput{
		Email: "bad@example.com", Password: "pass",
	}, "", "")
	assert.ErrorIs(t, err, user.ErrInvalidCredentials)
}

func TestLogin_InactiveUser(t *testing.T) {
	uc, mockUserRepo, _, _, _, _ := setupAuthUseCase(t)

	uid := ulid.New()
	u := &user.User{ID: uid, Email: "inactive@example.com", IsActive: false}

	mockUserRepo.EXPECT().GetByEmail(gomock.Any(), u.Email).Return(u, nil)

	_, err := uc.Login(context.Background(), ucauth.LoginInput{
		Email: u.Email, Password: "pass",
	}, "", "")
	assert.ErrorIs(t, err, user.ErrUserNotActive)
}

func TestVerifyMFA_BlocksDeactivatedUser(t *testing.T) {
	uc, mockUserRepo, _, _, mockMFARepo, _, tokenGen := setupAuthUseCaseWithTokenGen(t)

	ctx := context.Background()
	uid := ulid.New()

	mfaToken, err := tokenGen.GenerateMFAToken(uid)
	require.NoError(t, err)

	key, err := totp.Generate(totp.GenerateOpts{Issuer: "test", AccountName: "test@example.com"})
	require.NoError(t, err)
	code, err := totp.GenerateCode(key.Secret(), time.Now())
	require.NoError(t, err)

	mockMFARepo.EXPECT().GetByUserID(ctx, uid).Return(&user.MFAConfig{
		UserID:    uid,
		Secret:    key.Secret(),
		IsEnabled: true,
	}, nil)

	deactivatedAt := time.Now()
	mockUserRepo.EXPECT().GetByID(ctx, uid).Return(&user.User{
		ID:            uid,
		IsActive:      false,
		DeactivatedAt: &deactivatedAt,
	}, nil)

	_, err = uc.VerifyMFA(ctx, ucauth.VerifyMFAInput{
		MFAToken: mfaToken,
		TOTPCode: code,
	}, "ua", "127.0.0.1")
	require.Error(t, err)
}

func TestRefreshToken_BlocksDeactivatedUser(t *testing.T) {
	uc, mockUserRepo, _, mockSessionRepo, _, _ := setupAuthUseCase(t)

	ctx := context.Background()
	uid := ulid.New()
	sessionID := ulid.New()

	session := &domainauth.Session{
		ID:           sessionID,
		UserID:       uid,
		RefreshToken: "test-refresh-token",
		ExpiresAt:    time.Now().Add(time.Hour),
		CreatedAt:    time.Now(),
	}

	mockSessionRepo.EXPECT().
		GetByRefreshToken(ctx, "test-refresh-token").
		Return(session, nil)

	mockSessionRepo.EXPECT().
		Delete(ctx, sessionID).
		Return(nil)

	deactivatedAt := time.Now()
	mockUserRepo.EXPECT().GetByID(ctx, uid).Return(&user.User{
		ID:            uid,
		IsActive:      false,
		DeactivatedAt: &deactivatedAt,
	}, nil)

	_, err := uc.RefreshToken(ctx, ucauth.RefreshTokenInput{
		RefreshToken: "test-refresh-token",
	})
	require.Error(t, err)
}
