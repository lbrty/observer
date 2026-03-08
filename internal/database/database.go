package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

//go:generate mockgen -destination=mock/database.go -package=mock github.com/lbrty/observer/internal/database DB

// DBTX is satisfied by both *sqlx.DB and *sqlx.Tx, allowing repositories to
// work transparently inside or outside a transaction.
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type txKeyType struct{}

var txKey = txKeyType{}

// ExecCtx returns the active transaction stored in ctx (if any), or falls back
// to the provided *sqlx.DB. Repositories call this instead of r.db directly so
// that WithTx propagation is automatic.
func ExecCtx(ctx context.Context, db *sqlx.DB) DBTX {
	if tx, ok := ctx.Value(txKey).(*sqlx.Tx); ok && tx != nil {
		return tx
	}
	return db
}

// DB is the database interface.
type DB interface {
	Ping(ctx context.Context) error
	Close() error
	GetDB() *sqlx.DB
	// WithTx executes fn inside a single database transaction. The transaction is
	// propagated via ctx so repositories that call ExecCtx pick it up
	// automatically. fn returning a non-nil error causes an automatic rollback.
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type database struct {
	db *sqlx.DB
}

// New opens a new sqlx database connection.
func New(dsn string) (DB, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &database{db: db}, nil
}

func (d *database) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

func (d *database) Close() error {
	return d.db.Close()
}

func (d *database) GetDB() *sqlx.DB {
	return d.db
}

func (d *database) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	ctx = context.WithValue(ctx, txKey, tx)
	if err := fn(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
