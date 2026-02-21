package database

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

//go:generate mockgen -destination=mock/database.go -package=mock github.com/lbrty/observer/internal/database DB

// DB is the database interface.
type DB interface {
	Ping(ctx context.Context) error
	Close() error
	GetDB() *sqlx.DB
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
