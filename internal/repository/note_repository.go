package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lbrty/observer/internal/domain/note"
)

type personNoteRepo struct {
	db *sqlx.DB
}

// NewPersonNoteRepository creates a PersonNoteRepository.
func NewPersonNoteRepository(db *sqlx.DB) PersonNoteRepository {
	return &personNoteRepo{db: db}
}

func (r *personNoteRepo) List(ctx context.Context, personID string) ([]*note.Note, error) {
	const q = `SELECT id, person_id, author_id, body, created_at FROM person_notes WHERE person_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, q, personID)
	if err != nil {
		return nil, fmt.Errorf("list person notes: %w", err)
	}
	defer rows.Close()

	var out []*note.Note
	for rows.Next() {
		var n note.Note
		if err := rows.Scan(&n.ID, &n.PersonID, &n.AuthorID, &n.Body, &n.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan note: %w", err)
		}
		n.CreatedAt = n.CreatedAt.UTC()
		out = append(out, &n)
	}
	return out, rows.Err()
}

func (r *personNoteRepo) GetByID(ctx context.Context, id string) (*note.Note, error) {
	const q = `SELECT id, person_id, author_id, body, created_at FROM person_notes WHERE id = $1`
	var n note.Note
	err := r.db.QueryRowContext(ctx, q, id).Scan(&n.ID, &n.PersonID, &n.AuthorID, &n.Body, &n.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, note.ErrNoteNotFound
		}
		return nil, fmt.Errorf("get note: %w", err)
	}
	n.CreatedAt = n.CreatedAt.UTC()
	return &n, nil
}

func (r *personNoteRepo) Create(ctx context.Context, n *note.Note) error {
	const q = `INSERT INTO person_notes (id, person_id, author_id, body, created_at) VALUES ($1, $2, $3, $4, $5)`
	n.CreatedAt = time.Now().UTC()
	_, err := r.db.ExecContext(ctx, q, n.ID, n.PersonID, n.AuthorID, n.Body, n.CreatedAt)
	if err != nil {
		return fmt.Errorf("create note: %w", err)
	}
	return nil
}

func (r *personNoteRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM person_notes WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete note: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return note.ErrNoteNotFound
	}
	return nil
}
