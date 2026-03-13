package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lbrty/observer/internal/config"
	"github.com/lbrty/observer/internal/storage"
)

func TestNewS3Storage_EmptyBucket(t *testing.T) {
	_, err := storage.NewS3Storage(config.StorageConfig{
		S3Endpoint:  "http://localhost:9000",
		S3Region:    "us-east-1",
		S3AccessKey: "key",
		S3SecretKey: "secret",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "S3_BUCKET")
}
