package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lbrty/observer/internal/config"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear relevant env vars
	for _, key := range []string{"SERVER_HOST", "SERVER_PORT", "LOG_LEVEL", "JWT_ISSUER"} {
		os.Unsetenv(key)
	}

	cfg, err := config.Load()
	require.NoError(t, err)

	assert.Equal(t, config.DefaultServerHost, cfg.Server.Host)
	assert.Equal(t, config.DefaultServerPort, cfg.Server.Port)
	assert.Equal(t, config.DefaultServerReadTimeout, cfg.Server.ReadTimeout)
	assert.Equal(t, config.DefaultLogLevel, cfg.Log.Level)
	assert.Equal(t, "observer", cfg.JWT.Issuer)
	assert.Equal(t, 15*time.Minute, cfg.JWT.AccessTTL)
	assert.Equal(t, 168*time.Hour, cfg.JWT.RefreshTTL)
}

func TestLoad_EnvOverride(t *testing.T) {
	t.Setenv("SERVER_HOST", "0.0.0.0")
	t.Setenv("SERVER_PORT", "8080")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("JWT_ISSUER", "myapp")

	cfg, err := config.Load()
	require.NoError(t, err)

	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "debug", cfg.Log.Level)
	assert.Equal(t, "myapp", cfg.JWT.Issuer)
}
