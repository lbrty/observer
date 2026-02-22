package crypto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lbrty/observer/internal/crypto"
)

func TestArgonHasher_HashAndVerify(t *testing.T) {
	h := crypto.NewArgonHasher()

	hash, salt, err := h.Hash("mysecretpassword")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEmpty(t, salt)

	err = h.Verify("mysecretpassword", hash, salt)
	assert.NoError(t, err)
}

func TestArgonHasher_WrongPassword(t *testing.T) {
	h := crypto.NewArgonHasher()

	hash, salt, err := h.Hash("correct")
	require.NoError(t, err)

	err = h.Verify("wrong", hash, salt)
	assert.Error(t, err)
}

func TestArgonHasher_Uniqueness(t *testing.T) {
	h := crypto.NewArgonHasher()

	hash1, salt1, _ := h.Hash("samepassword")
	hash2, salt2, _ := h.Hash("samepassword")

	// Different salts each time
	assert.NotEqual(t, salt1, salt2)
	// Different hashes due to different salts
	assert.NotEqual(t, hash1, hash2)
}
