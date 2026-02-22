package crypto_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
	"time"

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

	privFile, _ := os.Create(privPath)
	pem.Encode(privFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	privFile.Close()

	pubBytes, _ := x509.MarshalPKIXPublicKey(&key.PublicKey)
	pubFile, _ := os.Create(pubPath)
	pem.Encode(pubFile, &pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})
	pubFile.Close()

	keys, err := crypto.LoadRSAKeys(privPath, pubPath)
	require.NoError(t, err)
	return keys
}

func TestRSATokenGenerator_AccessToken(t *testing.T) {
	keys := setupRSAKeys(t)
	gen := crypto.NewRSATokenGenerator(keys, 15*time.Minute, 168*time.Hour, 5*time.Minute, "observer")

	uid := ulid.New()
	token, expiresAt, err := gen.GenerateAccessToken(uid, "consultant")
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.True(t, expiresAt.After(time.Now()))

	claims, err := gen.ValidateAccessToken(token)
	require.NoError(t, err)
	assert.Equal(t, uid.String(), claims.UserID)
	assert.Equal(t, "consultant", claims.Role)
	assert.Equal(t, "access", claims.Type)
}

func TestRSATokenGenerator_MFAToken(t *testing.T) {
	keys := setupRSAKeys(t)
	gen := crypto.NewRSATokenGenerator(keys, 15*time.Minute, 168*time.Hour, 5*time.Minute, "observer")

	uid := ulid.New()
	token, err := gen.GenerateMFAToken(uid)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := gen.ValidateMFAToken(token)
	require.NoError(t, err)
	assert.Equal(t, uid.String(), claims.UserID)
	assert.Equal(t, "mfa_pending", claims.Type)
}

func TestRSATokenGenerator_TypeMismatch(t *testing.T) {
	keys := setupRSAKeys(t)
	gen := crypto.NewRSATokenGenerator(keys, 15*time.Minute, 168*time.Hour, 5*time.Minute, "observer")

	uid := ulid.New()
	accessToken, _, err := gen.GenerateAccessToken(uid, "consultant")
	require.NoError(t, err)

	// Access token should be rejected as MFA token
	_, err = gen.ValidateMFAToken(accessToken)
	assert.Error(t, err)
}

func TestRSATokenGenerator_RefreshToken(t *testing.T) {
	keys := setupRSAKeys(t)
	gen := crypto.NewRSATokenGenerator(keys, 15*time.Minute, 168*time.Hour, 5*time.Minute, "observer")

	token1, err := gen.GenerateRefreshToken()
	require.NoError(t, err)

	token2, err := gen.GenerateRefreshToken()
	require.NoError(t, err)

	assert.NotEmpty(t, token1)
	assert.NotEmpty(t, token2)
	assert.NotEqual(t, token1, token2)
}
