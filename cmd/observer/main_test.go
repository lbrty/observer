package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCmd_Initialized(t *testing.T) {
	assert.NotNil(t, rootCmd)
	assert.Equal(t, "observer", rootCmd.Use)
}
