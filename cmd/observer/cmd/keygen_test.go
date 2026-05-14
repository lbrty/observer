package cmd

import (
	"encoding/pem"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeygenCmd_Initialized(t *testing.T) {
	keygenCmd := NewKeygenCmd()
	assert.NotNil(t, keygenCmd)
	assert.Equal(t, "keygen", keygenCmd.Use)
}

func TestKeygenCmd_MinimumBits(t *testing.T) {
	keygenCmd := NewKeygenCmd()
	keygenCmd.SetOut(io.Discard)
	keygenCmd.Flags().Set("bits", "2048")
	keygenCmd.Flags().Set("output", t.TempDir())

	err := keygenCmd.RunE(keygenCmd, nil)
	assert.Error(t, err, "should fail with bits < 4096")
}

func TestKeygenCmd_GeneratesKeys(t *testing.T) {
	keygenCmd := NewKeygenCmd()
	tmpDir := t.TempDir()
	keygenCmd.SetOut(io.Discard)
	keygenCmd.Flags().Set("bits", "4096")
	keygenCmd.Flags().Set("output", tmpDir)

	err := keygenCmd.RunE(keygenCmd, nil)
	require.NoError(t, err)

	privPath := filepath.Join(tmpDir, "private_key.pem")
	pubPath := filepath.Join(tmpDir, "public_key.pem")

	assert.FileExists(t, privPath)
	assert.FileExists(t, pubPath)

	privData, err := os.ReadFile(privPath)
	require.NoError(t, err)
	block, _ := pem.Decode(privData)
	assert.NotNil(t, block, "private key should be valid PEM")
	assert.Equal(t, "RSA PRIVATE KEY", block.Type)

	pubData, err := os.ReadFile(pubPath)
	require.NoError(t, err)
	pubBlock, _ := pem.Decode(pubData)
	assert.NotNil(t, pubBlock, "public key should be valid PEM")
	assert.Equal(t, "PUBLIC KEY", pubBlock.Type)
}
