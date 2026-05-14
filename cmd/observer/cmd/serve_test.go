package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServeCmd_Initialized(t *testing.T) {
	serveCmd := NewServeCmd()
	assert.NotNil(t, serveCmd)
	assert.Equal(t, "serve", serveCmd.Use)
}

func TestServeCmd_Flags(t *testing.T) {
	serveCmd := NewServeCmd()
	hostFlag := serveCmd.Flags().Lookup("host")
	assert.NotNil(t, hostFlag)

	portFlag := serveCmd.Flags().Lookup("port")
	assert.NotNil(t, portFlag)
}
