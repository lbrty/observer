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

	"github.com/lbrty/observer/internal/domain/support"
)

type supportRecordRepo struct {
	db *sqlx.DB
}

// NewSupportRecordRepository creates a SupportRecordRepository.
func NewSupportRecordRepository(db *sqlx.DB) SupportRecordRepository {
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
	TimesToUTC(&r.CreatedAt, &r.UpdatedAt)
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
	TimesToUTC(&r.CreatedAt, &r.UpdatedAt)
	return &r, nil
}

func (r *supportRecordRepo) List(ctx context.Context, filter support.RecordListFilter) ([]*support.Record, int, error) {
	var (
		where []string
		args  []any
		ix    int
	)

	ix++
	where = append(where, "project_id = $"+strconv.Itoa(ix))
	args = append(args, filter.ProjectID)

	if filter.PersonID != nil {
		ix++
		where = append(where, "person_id = $"+strconv.Itoa(ix))
		args = append(args, *filter.PersonID)
	}
	if filter.ConsultantID != nil {
		ix++
		where = append(where, "consultant_id = $"+strconv.Itoa(ix))
		args = append(args, *filter.ConsultantID)
	}
	if filter.OfficeID != nil {
		ix++
		where = append(where, "office_id = $"+strconv.Itoa(ix))
		args = append(args, *filter.OfficeID)
	}
	if filter.Type != nil {
		ix++
		where = append(where, "type = $"+strconv.Itoa(ix))
		args = append(args, string(*filter.Type))
	}
	if filter.Sphere != nil {
		ix++
		where = append(where, "sphere = $"+strconv.Itoa(ix))
		args = append(args, string(*filter.Sphere))
	}
	if filter.ReferralStatus != nil {
		ix++
		where = append(where, "referral_status = $"+strconv.Itoa(ix))
		args = append(args, string(*filter.ReferralStatus))
	}
	if filter.DateFrom != nil {
		ix++
		where = append(where, "provided_at >= $"+strconv.Itoa(ix))
		args = append(args, *filter.DateFrom)
	}
	if filter.DateTo != nil {
		ix++
		where = append(where, "provided_at <= $"+strconv.Itoa(ix))
		args = append(args, *filter.DateTo)
	}

	// Prefix all filter columns with table alias for the JOIN query.
	for i, w := range where {
		where[i] = "sr." + w
	}
	whereClause := "WHERE " + strings.Join(where, " AND ")

	const joinClause = "FROM support_records sr LEFT JOIN people p ON p.id = sr.person_id "

	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) "+joinClause+whereClause, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count support records: %w", err)
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	perPage := filter.PerPage
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	ix++
	args = append(args, perPage)
	limitP := "$" + strconv.Itoa(ix)
	ix++
	args = append(args, offset)
	offsetP := "$" + strconv.Itoa(ix)

	q := "SELECT " + supportListCols + " " + joinClause +
		whereClause + " ORDER BY sr.created_at DESC LIMIT " + limitP + " OFFSET " + offsetP

	rows, err := r.db.QueryContext(ctx, q, args...)
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
	const q = `INSERT INTO support_records (
		id, person_id, project_id, consultant_id, recorded_by, office_id,
		referred_to_office, type, sphere, referral_status, provided_at, notes,
		created_at, updated_at
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`
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
	const q = `UPDATE support_records SET
		consultant_id=$2, office_id=$3, referred_to_office=$4, type=$5, sphere=$6,
		referral_status=$7, provided_at=$8, notes=$9, updated_at=$10
	WHERE id=$1`
	rec.UpdatedAt = time.Now().UTC()
	res, err := r.db.ExecContext(ctx, q,
		rec.ID, rec.ConsultantID, rec.OfficeID, rec.ReferredToOffice, rec.Type, rec.Sphere,
		rec.ReferralStatus, rec.ProvidedAt, rec.Notes, rec.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update support record: %w", err)
	}
	return CheckRowsAffected(res, support.ErrRecordNotFound)
}

func (r *supportRecordRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM support_records WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete support record: %w", err)
	}
	return CheckRowsAffected(res, support.ErrRecordNotFound)
}
