package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
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

func scanPet(row interface{ Scan(dest ...any) error }) (*pet.Pet, error) {
	var p pet.Pet
	if err := row.Scan(&p.ID, &p.ProjectID, &p.OwnerID, &p.Name, &p.Status, &p.RegistrationID, &p.Notes, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, err
	}
	TimesToUTC(&p.CreatedAt, &p.UpdatedAt)
	return &p, nil
}

func (r *petRepo) List(ctx context.Context, filter pet.PetListFilter) ([]*pet.Pet, int, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	perPage := filter.PerPage
	if perPage < 1 {
		perPage = 20
	}

	cond := sq.And{sq.Eq{"project_id": filter.ProjectID}}

	if filter.Status != nil {
		cond = append(cond, sq.Eq{"status": *filter.Status})
	}

	hasTags := len(filter.TagIDs) > 0
	havingClause := fmt.Sprintf("COUNT(DISTINCT pt.tag_id) = %d", len(filter.TagIDs))

	var countSQL string
	var countArgs []any
	var err error
	if hasTags {
		sub := psql.Select("pets.id").
			From("pets").
			Join("pet_tags pt ON pt.pet_id = pets.id").
			Where(cond).
			Where(sq.Eq{"pt.tag_id": filter.TagIDs}).
			GroupBy("pets.id").
			Having(havingClause)
		countSQL, countArgs, err = psql.Select("COUNT(*)").FromSelect(sub, "sub").ToSql()
	} else {
		countSQL, countArgs, err = psql.Select("COUNT(*)").From("pets").Where(cond).ToSql()
	}

	if err != nil {
		return nil, 0, fmt.Errorf("build count query: %w", err)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count pets: %w", err)
	}

	offset := (page - 1) * perPage

	const petColumns = `id, project_id, owner_id, name, status, registration_id, notes, created_at, updated_at`
	const petColumnsTagged = `pets.id, pets.project_id, pets.owner_id, pets.name, pets.status, pets.registration_id, pets.notes, pets.created_at, pets.updated_at`

	var listSQL string
	var listArgs []any
	if hasTags {
		listSQL, listArgs, err = psql.Select(petColumnsTagged).
			From("pets").
			Join("pet_tags pt ON pt.pet_id = pets.id").
			Where(cond).
			Where(sq.Eq{"pt.tag_id": filter.TagIDs}).
			GroupBy("pets.id").
			Having(havingClause).
			OrderBy("pets.created_at DESC").
			Limit(uint64(perPage)).
			Offset(uint64(offset)).
			ToSql()
	} else {
		listSQL, listArgs, err = psql.Select(petColumns).
			From("pets").
			Where(cond).
			OrderBy("created_at DESC").
			Limit(uint64(perPage)).
			Offset(uint64(offset)).
			ToSql()
	}

	if err != nil {
		return nil, 0, fmt.Errorf("build list query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, listSQL, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("list pets: %w", err)
	}
	defer rows.Close()

	var out []*pet.Pet
	for rows.Next() {
		p, err := scanPet(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan pet: %w", err)
		}
		out = append(out, p)
	}

	return out, total, rows.Err()
}

func (r *petRepo) GetByID(ctx context.Context, id string) (*pet.Pet, error) {
	const q = `SELECT id, project_id, owner_id, name, status, registration_id, notes, created_at, updated_at FROM pets WHERE id = $1`
	p, err := scanPet(r.db.QueryRowContext(ctx, q, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pet.ErrPetNotFound
		}
		return nil, fmt.Errorf("get pet: %w", err)
	}

	return p, nil
}

func (r *petRepo) Create(ctx context.Context, p *pet.Pet) error {
	const q = `
		INSERT INTO pets (id, project_id, owner_id, name, status, registration_id, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

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

	return CheckRowsAffected(res, pet.ErrPetNotFound)
}

func (r *petRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM pets WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete pet: %w", err)
	}

	return CheckRowsAffected(res, pet.ErrPetNotFound)
}

type petTagRepo struct {
	db *sqlx.DB
}

// NewPetTagRepository creates a PetTagRepository.
func NewPetTagRepository(db *sqlx.DB) PetTagRepository {
	return &petTagRepo{db: db}
}

func (r *petTagRepo) List(ctx context.Context, petID string) ([]string, error) {
	const q = `SELECT tag_id FROM pet_tags WHERE pet_id = $1 ORDER BY tag_id`
	rows, err := r.db.QueryContext(ctx, q, petID)
	if err != nil {
		return nil, fmt.Errorf("list pet tags: %w", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan tag id: %w", err)
		}
		ids = append(ids, id)
	}

	return ids, rows.Err()
}

func (r *petTagRepo) ListBulk(ctx context.Context, entityIDs []string) (map[string][]string, error) {
	if len(entityIDs) == 0 {
		return map[string][]string{}, nil
	}

	q, args, err := buildBulkTagQuery("pet_tags", "pet_id", entityIDs)
	if err != nil {
		return nil, fmt.Errorf("build list bulk: %w", err)
	}
	return queryBulkTags(ctx, r.db, q, args)
}

func (r *petTagRepo) ReplaceAll(ctx context.Context, petID string, tagIDs []string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM pet_tags WHERE pet_id = $1`, petID); err != nil {
		return fmt.Errorf("delete pet tags: %w", err)
	}

	if len(tagIDs) > 0 {
		iq := psql.Insert("pet_tags").Columns("pet_id", "tag_id")
		for _, tagID := range tagIDs {
			iq = iq.Values(petID, tagID)
		}
		sqlStr, args, err := iq.ToSql()
		if err != nil {
			return fmt.Errorf("build insert tags: %w", err)
		}
		if _, err := tx.ExecContext(ctx, sqlStr, args...); err != nil {
			return fmt.Errorf("insert pet tags: %w", err)
		}
	}

	return tx.Commit()
}
