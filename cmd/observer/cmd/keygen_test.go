package cmd

import (
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeygenCmd_Initialized(t *testing.T) {
	assert.NotNil(t, KeygenCmd)
	assert.Equal(t, "keygen", KeygenCmd.Use)
}

func TestKeygenCmd_MinimumBits(t *testing.T) {
	KeygenCmd.Flags().Set("bits", "2048")
	KeygenCmd.Flags().Set("output", t.TempDir())

	err := KeygenCmd.RunE(KeygenCmd, nil)
	assert.Error(t, err, "should fail with bits < 4096")

	// restore default
	KeygenCmd.Flags().Set("bits", "4096")
}

func TestKeygenCmd_GeneratesKeys(t *testing.T) {
	tmpDir := t.TempDir()
	KeygenCmd.Flags().Set("bits", "4096")
	KeygenCmd.Flags().Set("output", tmpDir)

	err := KeygenCmd.RunE(KeygenCmd, nil)
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
