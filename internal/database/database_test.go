package database_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lbrty/observer/internal/database"
	"github.com/lbrty/observer/internal/testutil"
)

func TestDatabase_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	dsn, cleanup := testutil.SetupPostgres(t)
	defer cleanup()

	db, err := database.New(dsn)
	require.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	err = db.Ping(ctx)
	assert.NoError(t, err)

	assert.NotNil(t, db.GetDB())
}
