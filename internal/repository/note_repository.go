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

const noteCols = `id, person_id, author_id, body, created_at, updated_at`

func scanNote(row interface{ Scan(dest ...any) error }) (*note.Note, error) {
	var n note.Note
	var updatedAt sql.NullTime
	err := row.Scan(&n.ID, &n.PersonID, &n.AuthorID, &n.Body, &n.CreatedAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	TimesToUTC(&n.CreatedAt)
	if updatedAt.Valid {
		t := updatedAt.Time.UTC()
		n.UpdatedAt = t
	}
	return &n, nil
}

func (r *personNoteRepo) List(ctx context.Context, personID string) ([]*note.Note, error) {
	q := "SELECT " + noteCols + " FROM person_notes WHERE person_id = $1 ORDER BY created_at DESC"
	rows, err := r.db.QueryContext(ctx, q, personID)
	if err != nil {
		return nil, fmt.Errorf("list person notes: %w", err)
	}
	defer rows.Close()

	var out []*note.Note
	for rows.Next() {
		n, err := scanNote(rows)
		if err != nil {
			return nil, fmt.Errorf("scan note: %w", err)
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

func (r *personNoteRepo) GetByID(ctx context.Context, id string) (*note.Note, error) {
	q := "SELECT " + noteCols + " FROM person_notes WHERE id = $1"
	n, err := scanNote(r.db.QueryRowContext(ctx, q, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, note.ErrNoteNotFound
		}
		return nil, fmt.Errorf("get note: %w", err)
	}
	return n, nil
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

func (r *personNoteRepo) Update(ctx context.Context, n *note.Note) error {
	const q = `UPDATE person_notes SET body=$2, updated_at=$3 WHERE id=$1`
	n.UpdatedAt = time.Now().UTC()
	res, err := r.db.ExecContext(ctx, q, n.ID, n.Body, n.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update note: %w", err)
	}
	return CheckRowsAffected(res, note.ErrNoteNotFound)
}

func (r *personNoteRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM person_notes WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete note: %w", err)
	}
	return CheckRowsAffected(res, note.ErrNoteNotFound)
}
