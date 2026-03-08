package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lbrty/observer/internal/domain/audit"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/testutil"
	"github.com/lbrty/observer/internal/ulid"
)

type auditTestFixture struct {
	db      *sqlx.DB
	repo    repository.AuditLogRepository
	userID  string
	cleanup func()
}

func setupAuditTest(t *testing.T) auditTestFixture {
	t.Helper()

	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	dsn, cleanup := testutil.SetupPostgres(t)
	testutil.MigrateUp(t, dsn)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		cleanup()
		t.Fatalf("connect to test db: %v", err)
	}

	userID := ulid.NewString()
	_, err = db.Exec(
		`INSERT INTO users (id, email, phone, role) VALUES ($1, $2, $3, 'admin')`,
		userID, "test-"+userID[:8]+"@example.com", "+99655"+userID[:7],
	)
	if err != nil {
		db.Close()
		cleanup()
		t.Fatalf("insert test user: %v", err)
	}

	for _, pid := range []string{"proj1", "proj2"} {
		_, err = db.Exec(
			`INSERT INTO projects (id, name, owner_id) VALUES ($1, $2, $3)`,
			pid, "Project "+pid, userID,
		)
		if err != nil {
			db.Close()
			cleanup()
			t.Fatalf("insert test project %s: %v", pid, err)
		}
	}

	return auditTestFixture{
		db:     db,
		repo:   repository.NewAuditLogRepository(db),
		userID: userID,
		cleanup: func() {
			db.Close()
			cleanup()
		},
	}
}

func TestAuditLogRepository_Log(t *testing.T) {
	f := setupAuditTest(t)
	defer f.cleanup()

	projID := "proj1"
	err := f.repo.Log(context.Background(), audit.Entry{
		ProjectID:  &projID,
		UserID:     f.userID,
		Action:     "create",
		EntityType: "person",
		Summary:    "created person",
	})
	require.NoError(t, err)
}

func TestAuditLogRepository_List_NoFilters(t *testing.T) {
	f := setupAuditTest(t)
	defer f.cleanup()

	ctx := context.Background()
	projID := "proj1"

	for i := 0; i < 3; i++ {
		require.NoError(t, f.repo.Log(ctx, audit.Entry{
			ProjectID:  &projID,
			UserID:     f.userID,
			Action:     "create",
			EntityType: "person",
			Summary:    "entry",
		}))
	}

	entries, total, err := f.repo.List(ctx, audit.Filter{Page: 1, PerPage: 20})
	require.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, entries, 3)
}

func TestAuditLogRepository_List_ByProjectID(t *testing.T) {
	f := setupAuditTest(t)
	defer f.cleanup()

	ctx := context.Background()
	proj1 := "proj1"
	proj2 := "proj2"

	for i := 0; i < 2; i++ {
		require.NoError(t, f.repo.Log(ctx, audit.Entry{
			ProjectID:  &proj1,
			UserID:     f.userID,
			Action:     "create",
			EntityType: "person",
			Summary:    "entry",
		}))
	}
	require.NoError(t, f.repo.Log(ctx, audit.Entry{
		ProjectID:  &proj2,
		UserID:     f.userID,
		Action:     "update",
		EntityType: "person",
		Summary:    "entry",
	}))

	entries, total, err := f.repo.List(ctx, audit.Filter{ProjectID: &proj1, Page: 1, PerPage: 20})
	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, entries, 2)
	for _, e := range entries {
		assert.Equal(t, &proj1, e.ProjectID)
	}
}

func TestAuditLogRepository_List_ByAction(t *testing.T) {
	f := setupAuditTest(t)
	defer f.cleanup()

	ctx := context.Background()
	projID := "proj1"

	require.NoError(t, f.repo.Log(ctx, audit.Entry{
		ProjectID: &projID, UserID: f.userID, Action: "create", EntityType: "person", Summary: "a",
	}))
	require.NoError(t, f.repo.Log(ctx, audit.Entry{
		ProjectID: &projID, UserID: f.userID, Action: "delete", EntityType: "person", Summary: "b",
	}))
	require.NoError(t, f.repo.Log(ctx, audit.Entry{
		ProjectID: &projID, UserID: f.userID, Action: "create", EntityType: "person", Summary: "c",
	}))

	action := "create"
	entries, total, err := f.repo.List(ctx, audit.Filter{Action: &action, Page: 1, PerPage: 20})
	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, entries, 2)
	for _, e := range entries {
		assert.Equal(t, "create", e.Action)
	}
}

func TestAuditLogRepository_List_DateRange(t *testing.T) {
	f := setupAuditTest(t)
	defer f.cleanup()

	ctx := context.Background()
	projID := "proj1"

	for _, ts := range []string{"2024-01-15 10:00:00+00", "2024-06-15 10:00:00+00", "2024-12-15 10:00:00+00"} {
		id := ulid.NewString()
		_, err := f.db.Exec(
			`INSERT INTO audit_logs (id, project_id, user_id, action, entity_type, summary, ip, user_agent, created_at)
			 VALUES ($1, $2, $3, 'view', 'person', 'entry', '', '', $4)`,
			id, projID, f.userID, ts,
		)
		require.NoError(t, err)
	}

	from := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 9, 30, 23, 59, 59, 0, time.UTC)
	entries, total, err := f.repo.List(ctx, audit.Filter{
		DateFrom: &from,
		DateTo:   &to,
		Page:     1,
		PerPage:  20,
	})
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, entries, 1)
}
