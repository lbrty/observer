package crypto_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lbrty/observer/internal/crypto"
	"github.com/lbrty/observer/internal/ulid"
)

func setupRSAKeys(t *testing.T) *crypto.RSAKeys {
	t.Helper()
	dir := t.TempDir()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	privPath := filepath.Join(dir, "private.pem")
	pubPath := filepath.Join(dir, "public.pem")

	privFile, err := os.Create(privPath)
	require.NoError(t, err)
	require.NoError(t, pem.Encode(privFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}))
	privFile.Close()

	pubBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	require.NoError(t, err)
	pubFile, err := os.Create(pubPath)
	require.NoError(t, err)
	require.NoError(t, pem.Encode(pubFile, &pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes}))
	pubFile.Close()

	keys, err := crypto.LoadRSAKeys(privPath, pubPath)
	require.NoError(t, err)
	return keys
}

func TestRSATokenGenerator_AccessToken(t *testing.T) {
	keys := setupRSAKeys(t)
	gen := crypto.NewRSATokenGenerator(keys, 15*time.Minute, 5*time.Minute, "observer")

	uid := ulid.New()
	token, expiresAt, err := gen.GenerateAccessToken(uid, "consultant")
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.True(t, expiresAt.After(time.Now()))

	claims, err := gen.ValidateAccessToken(token)
	require.NoError(t, err)
	assert.Equal(t, uid.String(), claims.Subject)
	assert.Equal(t, "consultant", claims.Role)
	assert.Equal(t, "access", claims.Type)
}

func TestRSATokenGenerator_MFAToken(t *testing.T) {
	keys := setupRSAKeys(t)
	gen := crypto.NewRSATokenGenerator(keys, 15*time.Minute, 5*time.Minute, "observer")

	uid := ulid.New()
	token, err := gen.GenerateMFAToken(uid)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := gen.ValidateMFAToken(token)
	require.NoError(t, err)
	assert.Equal(t, uid.String(), claims.Subject)
	assert.Equal(t, "mfa_pending", claims.Type)
}

func TestRSATokenGenerator_TypeMismatch(t *testing.T) {
	keys := setupRSAKeys(t)
	gen := crypto.NewRSATokenGenerator(keys, 15*time.Minute, 5*time.Minute, "observer")

	uid := ulid.New()
	accessToken, _, err := gen.GenerateAccessToken(uid, "consultant")
	require.NoError(t, err)

	// Access token should be rejected as MFA token
	_, err = gen.ValidateMFAToken(accessToken)
	assert.Error(t, err)
}

func TestRSATokenGenerator_RejectsWrongIssuer(t *testing.T) {
	keys := setupRSAKeys(t)
	gen := crypto.NewRSATokenGenerator(keys, 15*time.Minute, 5*time.Minute, "observer")
	uid := ulid.New()
	claims := crypto.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "other-service",
			Subject:   uid.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		Role: "admin",
		Type: "access",
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(keys.PrivateKey)
	require.NoError(t, err)

	_, err = gen.ValidateAccessToken(token)
	require.Error(t, err)
}

func TestRSATokenGenerator_RejectsOtherRSAAlgorithm(t *testing.T) {
	keys := setupRSAKeys(t)
	gen := crypto.NewRSATokenGenerator(keys, 15*time.Minute, 5*time.Minute, "observer")
	uid := ulid.New()
	claims := crypto.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "observer",
			Subject:   uid.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		Role: "admin",
		Type: "access",
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS512, claims).SignedString(keys.PrivateKey)
	require.NoError(t, err)

	_, err = gen.ValidateAccessToken(token)
	require.Error(t, err)
}

func TestRSATokenGenerator_RefreshToken(t *testing.T) {
	keys := setupRSAKeys(t)
	gen := crypto.NewRSATokenGenerator(keys, 15*time.Minute, 5*time.Minute, "observer")

	token1, err := gen.GenerateRefreshToken()
	require.NoError(t, err)

	token2, err := gen.GenerateRefreshToken()
	require.NoError(t, err)

	assert.NotEmpty(t, token1)
	assert.NotEmpty(t, token2)
	assert.NotEqual(t, token1, token2)
}

func TestGenerateRefreshToken_Entropy(t *testing.T) {
	keys := setupRSAKeys(t)
	gen := crypto.NewRSATokenGenerator(keys, time.Hour, 5*time.Minute, "test")

	tok1, err := gen.GenerateRefreshToken()
	require.NoError(t, err)

	tok2, err := gen.GenerateRefreshToken()
	require.NoError(t, err)

	// Must be hex-encoded 32 bytes = 64 chars
	assert.Len(t, tok1, 64, "refresh token must be 64 hex chars (32 random bytes)")
	assert.Len(t, tok2, 64)
	assert.NotEqual(t, tok1, tok2)

	_, err = hex.DecodeString(tok1)
	assert.NoError(t, err, "refresh token must be valid hex")
}
