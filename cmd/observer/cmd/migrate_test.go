package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrateCmd_Initialized(t *testing.T) {
	assert.NotNil(t, MigrateCmd)
	assert.Equal(t, "migrate", MigrateCmd.Use)
}

func TestMigrateCmd_SubcommandsRegistered(t *testing.T) {
	names := make(map[string]bool)
	for _, sub := range MigrateCmd.Commands() {
		names[sub.Use] = true
	}

	assert.True(t, names["up"])
	assert.True(t, names["version"])

	// create requires exactly 1 argument
	assert.NotNil(t, migrateCreateCmd.Args)
}
