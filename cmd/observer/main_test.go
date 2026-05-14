package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lbrty/observer/cmd/observer/cmd"
)

func TestRootCmd_Initialized(t *testing.T) {
	rootCmd := cmd.NewRootCmd()
	assert.NotNil(t, rootCmd)
	assert.Equal(t, "observer", rootCmd.Use)
}
