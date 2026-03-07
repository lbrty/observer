package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
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

func scanPet(row interface{ Scan(dest ...any) error }) (*pet.Pet, error) {
	var p pet.Pet
	if err := row.Scan(&p.ID, &p.ProjectID, &p.OwnerID, &p.Name, &p.Status, &p.RegistrationID, &p.Notes, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, err
	}
	TimesToUTC(&p.CreatedAt, &p.UpdatedAt)
	return &p, nil
}

func (r *petRepo) List(ctx context.Context, projectID string, status string, tagIDs []string, page, perPage int) ([]*pet.Pet, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}

	where := []string{"project_id = $1"}
	args := []any{projectID}
	ix := 1

	if status != "" {
		ix++
		where = append(where, "status = $"+strconv.Itoa(ix))
		args = append(args, status)
	}

	var tagJoin string
	if len(tagIDs) > 0 {
		placeholders := make([]string, len(tagIDs))
		for i, tagID := range tagIDs {
			ix++
			placeholders[i] = "$" + strconv.Itoa(ix)
			args = append(args, tagID)
		}
		tagJoin = " JOIN pet_tags pt ON pt.pet_id = pets.id AND pt.tag_id IN (" + strings.Join(placeholders, ",") + ")"
		where = append(where, "1=1 GROUP BY pets.id HAVING COUNT(DISTINCT pt.tag_id) = "+strconv.Itoa(len(tagIDs)))
	}

	whereClause := "WHERE " + strings.Join(where, " AND ")

	var countQ string
	if tagJoin != "" {
		countQ = "SELECT COUNT(*) FROM (SELECT pets.id FROM pets" + tagJoin + " " + whereClause + ") sub"
	} else {
		countQ = "SELECT COUNT(*) FROM pets " + whereClause
	}
	var total int
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count pets: %w", err)
	}

	offset := (page - 1) * perPage
	ix++
	limitParam := "$" + strconv.Itoa(ix)
	ix++
	offsetParam := "$" + strconv.Itoa(ix)
	args = append(args, perPage, offset)

	var q string
	if tagJoin != "" {
		q = fmt.Sprintf(`SELECT pets.id, pets.project_id, pets.owner_id, pets.name, pets.status, pets.registration_id, pets.notes, pets.created_at, pets.updated_at
			FROM pets%s %s ORDER BY pets.created_at DESC LIMIT %s OFFSET %s`, tagJoin, whereClause, limitParam, offsetParam)
	} else {
		q = fmt.Sprintf(`SELECT id, project_id, owner_id, name, status, registration_id, notes, created_at, updated_at
			FROM pets %s ORDER BY created_at DESC LIMIT %s OFFSET %s`, whereClause, limitParam, offsetParam)
	}
	rows, err := r.db.QueryContext(ctx, q, args...)
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
	q, args := buildBulkTagQuery("pet_tags", "pet_id", entityIDs)
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
		var sb strings.Builder
		sb.WriteString("INSERT INTO pet_tags (pet_id, tag_id) VALUES ")
		args := make([]any, 0, len(tagIDs)*2)
		for i, tagID := range tagIDs {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
			args = append(args, petID, tagID)
		}
		if _, err := tx.ExecContext(ctx, sb.String(), args...); err != nil {
			return fmt.Errorf("insert pet tags: %w", err)
		}
	}

	return tx.Commit()
}
