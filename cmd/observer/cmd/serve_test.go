package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServeCmd_Initialized(t *testing.T) {
	assert.NotNil(t, ServeCmd)
	assert.Equal(t, "serve", ServeCmd.Use)
}

func TestServeCmd_Flags(t *testing.T) {
	hostFlag := ServeCmd.Flags().Lookup("host")
	assert.NotNil(t, hostFlag)

	portFlag := ServeCmd.Flags().Lookup("port")
	assert.NotNil(t, portFlag)
}
