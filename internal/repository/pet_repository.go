package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lbrty/observer/internal/domain/pet"
)

type petRepo struct {
	db *sqlx.DB
}

// NewPetRepository creates a PetRepository.
func NewPetRepository(db *sqlx.DB) PetRepository {
	return &petRepo{db: db}
}

func (r *petRepo) List(ctx context.Context, projectID string, page, perPage int) ([]*pet.Pet, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM pets WHERE project_id = $1`, projectID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count pets: %w", err)
	}

	offset := (page - 1) * perPage
	const q = `SELECT id, project_id, owner_id, name, status, registration_id, notes, created_at, updated_at
		FROM pets WHERE project_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, q, projectID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list pets: %w", err)
	}
	defer rows.Close()

	var out []*pet.Pet
	for rows.Next() {
		var p pet.Pet
		if err := rows.Scan(&p.ID, &p.ProjectID, &p.OwnerID, &p.Name, &p.Status, &p.RegistrationID, &p.Notes, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan pet: %w", err)
		}
		p.CreatedAt = p.CreatedAt.UTC()
		p.UpdatedAt = p.UpdatedAt.UTC()
		out = append(out, &p)
	}
	return out, total, rows.Err()
}

func (r *petRepo) GetByID(ctx context.Context, id string) (*pet.Pet, error) {
	const q = `SELECT id, project_id, owner_id, name, status, registration_id, notes, created_at, updated_at FROM pets WHERE id = $1`
	var p pet.Pet
	err := r.db.QueryRowContext(ctx, q, id).Scan(&p.ID, &p.ProjectID, &p.OwnerID, &p.Name, &p.Status, &p.RegistrationID, &p.Notes, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pet.ErrPetNotFound
		}
		return nil, fmt.Errorf("get pet: %w", err)
	}
	p.CreatedAt = p.CreatedAt.UTC()
	p.UpdatedAt = p.UpdatedAt.UTC()
	return &p, nil
}

func (r *petRepo) Create(ctx context.Context, p *pet.Pet) error {
	const q = `INSERT INTO pets (id, project_id, owner_id, name, status, registration_id, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now
	_, err := r.db.ExecContext(ctx, q, p.ID, p.ProjectID, p.OwnerID, p.Name, p.Status, p.RegistrationID, p.Notes, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create pet: %w", err)
	}
	return nil
}

func (r *petRepo) Update(ctx context.Context, p *pet.Pet) error {
	const q = `UPDATE pets SET owner_id=$2, name=$3, status=$4, registration_id=$5, notes=$6, updated_at=$7 WHERE id=$1`
	p.UpdatedAt = time.Now().UTC()
	res, err := r.db.ExecContext(ctx, q, p.ID, p.OwnerID, p.Name, p.Status, p.RegistrationID, p.Notes, p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update pet: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return pet.ErrPetNotFound
	}
	return nil
}

func (r *petRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM pets WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete pet: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return pet.ErrPetNotFound
	}
	return nil
}
