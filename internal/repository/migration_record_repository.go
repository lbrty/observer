package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lbrty/observer/internal/domain/migration"
)

type migrationRecordRepo struct {
	db *sqlx.DB
}

// NewMigrationRecordRepository creates a MigrationRecordRepository.
func NewMigrationRecordRepository(db *sqlx.DB) MigrationRecordRepository {
	return &migrationRecordRepo{db: db}
}

const migrationCols = `id, person_id, from_place_id, destination_place_id, migration_date,
	movement_reason, housing_at_destination, notes, created_at, updated_at`

func scanMigration(row interface{ Scan(dest ...any) error }) (*migration.Record, error) {
	var r migration.Record
	var updatedAt sql.NullTime
	err := row.Scan(
		&r.ID, &r.PersonID, &r.FromPlaceID, &r.DestinationPlaceID, &r.MigrationDate,
		&r.MovementReason, &r.HousingAtDestination, &r.Notes, &r.CreatedAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}
	r.CreatedAt = r.CreatedAt.UTC()
	if updatedAt.Valid {
		r.UpdatedAt = updatedAt.Time.UTC()
	}
	return &r, nil
}

func (r *migrationRecordRepo) ListByPerson(ctx context.Context, personID string) ([]*migration.Record, error) {
	q := "SELECT " + migrationCols + " FROM migration_records WHERE person_id = $1 ORDER BY created_at DESC"
	rows, err := r.db.QueryContext(ctx, q, personID)
	if err != nil {
		return nil, fmt.Errorf("list migration records: %w", err)
	}
	defer rows.Close()

	var out []*migration.Record
	for rows.Next() {
		rec, err := scanMigration(rows)
		if err != nil {
			return nil, fmt.Errorf("scan migration record: %w", err)
		}
		out = append(out, rec)
	}
	return out, rows.Err()
}

func (r *migrationRecordRepo) GetByID(ctx context.Context, id string) (*migration.Record, error) {
	q := "SELECT " + migrationCols + " FROM migration_records WHERE id = $1"
	rec, err := scanMigration(r.db.QueryRowContext(ctx, q, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, migration.ErrRecordNotFound
		}
		return nil, fmt.Errorf("get migration record: %w", err)
	}
	return rec, nil
}

func (r *migrationRecordRepo) Create(ctx context.Context, rec *migration.Record) error {
	const q = `INSERT INTO migration_records (
		id, person_id, from_place_id, destination_place_id, migration_date,
		movement_reason, housing_at_destination, notes, created_at
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	rec.CreatedAt = time.Now().UTC()
	_, err := r.db.ExecContext(ctx, q,
		rec.ID, rec.PersonID, rec.FromPlaceID, rec.DestinationPlaceID, rec.MigrationDate,
		rec.MovementReason, rec.HousingAtDestination, rec.Notes, rec.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create migration record: %w", err)
	}
	return nil
}

func (r *migrationRecordRepo) Update(ctx context.Context, rec *migration.Record) error {
	const q = `UPDATE migration_records SET
		from_place_id=$2, destination_place_id=$3, migration_date=$4,
		movement_reason=$5, housing_at_destination=$6, notes=$7, updated_at=$8
	WHERE id=$1`
	rec.UpdatedAt = time.Now().UTC()
	res, err := r.db.ExecContext(ctx, q,
		rec.ID, rec.FromPlaceID, rec.DestinationPlaceID, rec.MigrationDate,
		rec.MovementReason, rec.HousingAtDestination, rec.Notes, rec.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update migration record: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return migration.ErrRecordNotFound
	}
	return nil
}
