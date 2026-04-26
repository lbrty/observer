package support

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/lbrty/observer/internal/domain/support"
	"github.com/lbrty/observer/internal/repository"
)

type supportRecordRepo struct {
	db *sqlx.DB
}

// New creates a SupportRecordRepository.
func New(db *sqlx.DB) repository.SupportRecordRepository {
	return &supportRecordRepo{db: db}
}

const supportCols = `id, person_id, project_id, consultant_id, recorded_by, office_id,
	referred_to_office, type, sphere, referral_status, provided_at, notes, created_at, updated_at`

const supportListCols = `sr.id, sr.person_id, sr.project_id, sr.consultant_id, sr.recorded_by, sr.office_id,
	sr.referred_to_office, sr.type, sr.sphere, sr.referral_status, sr.provided_at, sr.notes,
	sr.created_at, sr.updated_at, p.first_name, p.last_name`

func scanSupport(row interface{ Scan(dest ...any) error }) (*support.Record, error) {
	var r support.Record
	err := row.Scan(
		&r.ID, &r.PersonID, &r.ProjectID, &r.ConsultantID, &r.RecordedBy, &r.OfficeID,
		&r.ReferredToOffice, &r.Type, &r.Sphere, &r.ReferralStatus, &r.ProvidedAt, &r.Notes,
		&r.CreatedAt, &r.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	repository.TimesToUTC(&r.CreatedAt, &r.UpdatedAt)
	return &r, nil
}

func scanSupportWithPerson(row interface{ Scan(dest ...any) error }) (*support.Record, error) {
	var r support.Record
	err := row.Scan(
		&r.ID, &r.PersonID, &r.ProjectID, &r.ConsultantID, &r.RecordedBy, &r.OfficeID,
		&r.ReferredToOffice, &r.Type, &r.Sphere, &r.ReferralStatus, &r.ProvidedAt, &r.Notes,
		&r.CreatedAt, &r.UpdatedAt, &r.PersonFirstName, &r.PersonLastName,
	)

	if err != nil {
		return nil, err
	}

	repository.TimesToUTC(&r.CreatedAt, &r.UpdatedAt)
	return &r, nil
}

func (r *supportRecordRepo) List(ctx context.Context, filter support.RecordListFilter) ([]*support.Record, int, error) {
	cond := sq.And{sq.Eq{"sr.project_id": filter.ProjectID}}
	if filter.PersonID != nil {
		cond = append(cond, sq.Eq{"sr.person_id": *filter.PersonID})
	}

	if filter.ConsultantID != nil {
		cond = append(cond, sq.Eq{"sr.consultant_id": *filter.ConsultantID})
	}

	if filter.OfficeID != nil {
		cond = append(cond, sq.Eq{"sr.office_id": *filter.OfficeID})
	}

	if filter.Type != nil {
		cond = append(cond, sq.Eq{"sr.type": string(*filter.Type)})
	}

	if filter.Sphere != nil {
		cond = append(cond, sq.Eq{"sr.sphere": string(*filter.Sphere)})
	}

	if filter.ReferralStatus != nil {
		cond = append(cond, sq.Eq{"sr.referral_status": string(*filter.ReferralStatus)})
	}

	if filter.DateFrom != nil {
		cond = append(cond, sq.GtOrEq{"sr.provided_at": *filter.DateFrom})
	}

	if filter.DateTo != nil {
		cond = append(cond, sq.LtOrEq{"sr.provided_at": *filter.DateTo})
	}

	countSQL, countArgs, err := repository.PSQL.Select("COUNT(*)").From("support_records sr").Where(cond).ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build count support records query: %w", err)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count support records: %w", err)
	}

	perPage, off := repository.NormalizePagination(filter.Page, filter.PerPage)
	offset := uint64(off)

	listSQL, listArgs, err := repository.PSQL.
		Select(supportListCols).
		From("support_records sr").
		LeftJoin("people p ON p.id = sr.person_id").
		Where(cond).
		OrderBy("sr.created_at DESC").
		Limit(uint64(perPage)).
		Offset(offset).
		ToSql()

	if err != nil {
		return nil, 0, fmt.Errorf("build list support records query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, listSQL, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("list support records: %w", err)
	}
	defer rows.Close()

	var out []*support.Record
	for rows.Next() {
		rec, err := scanSupportWithPerson(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan support record: %w", err)
		}
		out = append(out, rec)
	}

	return out, total, rows.Err()
}

func (r *supportRecordRepo) GetByID(ctx context.Context, id string) (*support.Record, error) {
	q := "SELECT " + supportCols + " FROM support_records WHERE id = $1"
	rec, err := scanSupport(r.db.QueryRowContext(ctx, q, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, support.ErrRecordNotFound
		}
		return nil, fmt.Errorf("get support record: %w", err)
	}

	return rec, nil
}

func (r *supportRecordRepo) Create(ctx context.Context, rec *support.Record) error {
	const q = `
		INSERT INTO support_records (
			id, person_id, project_id, consultant_id, recorded_by, office_id,
			referred_to_office, type, sphere, referral_status, provided_at, notes,
			created_at, updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
	`
	now := time.Now().UTC()
	rec.CreatedAt = now
	rec.UpdatedAt = now
	_, err := r.db.ExecContext(ctx, q,
		rec.ID, rec.PersonID, rec.ProjectID, rec.ConsultantID, rec.RecordedBy, rec.OfficeID,
		rec.ReferredToOffice, rec.Type, rec.Sphere, rec.ReferralStatus, rec.ProvidedAt, rec.Notes,
		rec.CreatedAt, rec.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("create support record: %w", err)
	}

	return nil
}

func (r *supportRecordRepo) Update(ctx context.Context, rec *support.Record) error {
	const q = `
		UPDATE support_records
		SET
			consultant_id=$2, office_id=$3, referred_to_office=$4, type=$5, sphere=$6,
			referral_status=$7, provided_at=$8, notes=$9, updated_at=$10
		WHERE id=$1
	`

	rec.UpdatedAt = time.Now().UTC()
	res, err := r.db.ExecContext(ctx, q,
		rec.ID, rec.ConsultantID, rec.OfficeID, rec.ReferredToOffice, rec.Type, rec.Sphere,
		rec.ReferralStatus, rec.ProvidedAt, rec.Notes, rec.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("update support record: %w", err)
	}

	return repository.CheckRowsAffected(res, support.ErrRecordNotFound)
}

func (r *supportRecordRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM support_records WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete support record: %w", err)
	}

	return repository.CheckRowsAffected(res, support.ErrRecordNotFound)
}
