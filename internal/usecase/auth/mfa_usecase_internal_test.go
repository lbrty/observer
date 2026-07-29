package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/crypto"
	cryptomock "github.com/lbrty/observer/internal/crypto/mock"
	domainauth "github.com/lbrty/observer/internal/domain/auth"
	"github.com/lbrty/observer/internal/domain/user"
	repomock "github.com/lbrty/observer/internal/repository/mock"
)

func TestGenerateRecoveryCodesStoresDeterministicHashes(t *testing.T) {
	userID := ulid.Make()

	plain, records, err := generateRecoveryCodes(2, userID)
	require.NoError(t, err)
	require.Len(t, plain, 2)
	require.Len(t, records, 2)
	require.NotEqual(t, plain[0], plain[1])

	for i, code := range plain {
		require.Len(t, code, 32)
		hash := sha256.Sum256([]byte(code))
		require.Equal(t, hex.EncodeToString(hash[:]), records[i].CodeHash)
		require.Equal(t, userID, records[i].UserID)
	}
}

func TestVerifyMFAConsumesRecoveryCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepo := repomock.NewMockUserRepository(ctrl)
	sessionRepo := repomock.NewMockSessionRepository(ctrl)
	mfaRepo := repomock.NewMockMFARepository(ctrl)
	recoveryRepo := repomock.NewMockMFARecoveryCodeRepository(ctrl)
	tokenGen := cryptomock.NewMockTokenGenerator(ctrl)
	userID := ulid.Make()
	code := "0123456789abcdef0123456789abcdef"
	hash := sha256.Sum256([]byte(code))

	tokenGen.EXPECT().ValidateMFAToken("mfa-token").Return(&crypto.Claims{
		RegisteredClaims: jwt.RegisteredClaims{Subject: userID.String()},
		Type:             "mfa_pending",
	}, nil)
	mfaRepo.EXPECT().GetByUserID(gomock.Any(), userID).Return(&user.MFAConfig{
		UserID:    userID,
		Secret:    "invalid-totp-secret",
		IsEnabled: true,
	}, nil)
	recoveryRepo.EXPECT().ConsumeUnused(gomock.Any(), userID, hex.EncodeToString(hash[:])).Return(nil)
	userRepo.EXPECT().GetByID(gomock.Any(), userID).Return(&user.User{
		ID:       userID,
		Role:     user.RoleGuest,
		IsActive: true,
	}, nil)
	tokenGen.EXPECT().GenerateAccessToken(userID, string(user.RoleGuest)).
		Return("access-token", time.Now().Add(time.Minute), nil)
	tokenGen.EXPECT().GenerateRefreshToken().Return("refresh-token", nil)
	sessionRepo.EXPECT().Create(gomock.Any(), gomock.AssignableToTypeOf(&domainauth.Session{})).Return(nil)

	uc := NewAuthUseCase(userRepo, nil, sessionRepo, mfaRepo, recoveryRepo, nil, tokenGen, nil, time.Hour)
	output, err := uc.VerifyMFA(context.Background(), VerifyMFAInput{
		MFAToken: "mfa-token",
		TOTPCode: code,
	}, "agent", "192.0.2.1")

	require.NoError(t, err)
	require.Equal(t, "access-token", output.Tokens.AccessToken)
}
