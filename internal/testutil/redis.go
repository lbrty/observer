package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

// SetupRedis spins up a Redis testcontainer and returns the address and a cleanup function.
// The test is skipped automatically if Docker is not available.
func SetupRedis(t *testing.T) (string, func()) {
	t.Helper()

	ctx := context.Background()

	ctr, err := redis.Run(ctx,
		"redis:8-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Skipf("could not start redis container: %v", err)
		return "", func() {}
	}

	addr, err := ctr.Endpoint(ctx, "")
	if err != nil {
		t.Fatalf("get redis endpoint: %v", err)
	}

	return addr, func() {
		if err := ctr.Terminate(ctx); err != nil {
			t.Logf("terminate redis container: %v", err)
		}
	}
}
