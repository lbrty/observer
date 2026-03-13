//go:build !short

package storage_test

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/lbrty/observer/internal/config"
	"github.com/lbrty/observer/internal/storage"
)

const (
	minioAccessKey = "minioadmin"
	minioSecretKey = "minioadmin"
	testBucket     = "test-bucket"
)

func setupMinIO(t *testing.T) string {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "minio/minio:latest",
			ExposedPorts: []string{"9000/tcp"},
			Env: map[string]string{
				"MINIO_ROOT_USER":     minioAccessKey,
				"MINIO_ROOT_PASSWORD": minioSecretKey,
			},
			Cmd:        []string{"server", "/data"},
			WaitingFor: wait.ForHTTP("/minio/health/live").WithPort("9000/tcp"),
		},
		Started: true,
	})
	require.NoError(t, err)
	t.Cleanup(func() { container.Terminate(ctx) }) //nolint:errcheck

	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "9000/tcp")
	require.NoError(t, err)

	endpoint := fmt.Sprintf("http://%s:%s", host, port.Port())

	// Create test bucket using the SDK directly.
	cfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion("us-east-1"),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(minioAccessKey, minioSecretKey, ""),
		),
	)
	require.NoError(t, err)

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})
	_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(testBucket)})
	require.NoError(t, err)

	return endpoint
}

func newTestS3Storage(t *testing.T, endpoint string) *storage.S3Storage {
	t.Helper()
	s, err := storage.NewS3Storage(config.StorageConfig{
		S3Endpoint:  endpoint,
		S3Region:    "us-east-1",
		S3Bucket:    testBucket,
		S3AccessKey: minioAccessKey,
		S3SecretKey: minioSecretKey,
	})
	require.NoError(t, err)
	return s
}

func TestS3Storage_SaveOpenDelete(t *testing.T) {
	endpoint := setupMinIO(t)
	s := newTestS3Storage(t, endpoint)
	ctx := context.Background()

	const path = "projects/abc/documents/file.txt"
	const content = "hello observer"

	require.NoError(t, s.Save(ctx, path, strings.NewReader(content)))

	rc, err := s.Open(ctx, path)
	require.NoError(t, err)
	defer rc.Close()

	got, err := io.ReadAll(rc)
	require.NoError(t, err)
	assert.Equal(t, content, string(got))

	require.NoError(t, s.Delete(ctx, path))

	_, err = s.Open(ctx, path)
	assert.Error(t, err)
}

func TestS3Storage_DeleteMissingKey(t *testing.T) {
	endpoint := setupMinIO(t)
	s := newTestS3Storage(t, endpoint)

	// DeleteObject on S3 is idempotent — missing key must not return an error.
	err := s.Delete(context.Background(), "nonexistent/path/file.txt")
	assert.NoError(t, err)
}
