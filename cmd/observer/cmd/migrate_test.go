package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrateCmd_Initialized(t *testing.T) {
	migrateCmd := NewMigrateCmd()
	assert.NotNil(t, migrateCmd)
	assert.Equal(t, "migrate", migrateCmd.Use)
}

func TestMigrateCmd_SubcommandsRegistered(t *testing.T) {
	migrateCmd := NewMigrateCmd()
	names := make(map[string]bool)
	for _, sub := range migrateCmd.Commands() {
		names[sub.Use] = true
	}

	assert.True(t, names["up"])
	assert.True(t, names["version"])
	assert.True(t, names["create [name]"])
}
