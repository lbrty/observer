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
	"github.com/lib/pq"

	"github.com/lbrty/observer/internal/domain/person"
)

type personRepo struct {
	db *sqlx.DB
}

// NewPersonRepository creates a PersonRepository.
func NewPersonRepository(db *sqlx.DB) PersonRepository {
	return &personRepo{db: db}
}

const personColumns = `id, project_id, consultant_id, office_id, current_place_id, origin_place_id,
	external_id, first_name, last_name, patronymic, email, birth_date, sex, age_group,
	primary_phone, phone_numbers, case_status, consent_given, consent_date, registered_at,
	created_at, updated_at`

func scanPerson(row interface{ Scan(dest ...any) error }) (*person.Person, error) {
	var p person.Person
	err := row.Scan(
		&p.ID, &p.ProjectID, &p.ConsultantID, &p.OfficeID, &p.CurrentPlaceID, &p.OriginPlaceID,
		&p.ExternalID, &p.FirstName, &p.LastName, &p.Patronymic, &p.Email, &p.BirthDate, &p.Sex, &p.AgeGroup,
		&p.PrimaryPhone, &p.PhoneNumbers, &p.CaseStatus, &p.ConsentGiven, &p.ConsentDate, &p.RegisteredAt,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	p.CreatedAt = p.CreatedAt.UTC()
	p.UpdatedAt = p.UpdatedAt.UTC()
	return &p, nil
}

func (r *personRepo) List(ctx context.Context, filter person.PersonListFilter) ([]*person.Person, int, error) {
	var (
		where []string
		args  []any
		ix    int
	)

	ix++
	where = append(where, "project_id = $"+strconv.Itoa(ix))
	args = append(args, filter.ProjectID)

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
	if filter.CaseStatus != nil {
		ix++
		where = append(where, "case_status = $"+strconv.Itoa(ix))
		args = append(args, string(*filter.CaseStatus))
	}
	if filter.Search != nil && *filter.Search != "" {
		ix++
		where = append(where, "(first_name % $"+strconv.Itoa(ix)+" OR last_name % $"+strconv.Itoa(ix)+")")
		args = append(args, *filter.Search)
	}

	whereClause := "WHERE " + strings.Join(where, " AND ")

	countQ := "SELECT COUNT(*) FROM people " + whereClause
	var total int
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count people: %w", err)
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
	limitParam := "$" + strconv.Itoa(ix)
	ix++
	args = append(args, offset)
	offsetParam := "$" + strconv.Itoa(ix)

	q := "SELECT " + personColumns + " FROM people " +
		whereClause + " ORDER BY created_at DESC LIMIT " + limitParam + " OFFSET " + offsetParam

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list people: %w", err)
	}
	defer rows.Close()

	var out []*person.Person
	for rows.Next() {
		p, err := scanPerson(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan person: %w", err)
		}
		out = append(out, p)
	}
	return out, total, rows.Err()
}

func (r *personRepo) GetByID(ctx context.Context, id string) (*person.Person, error) {
	q := "SELECT " + personColumns + " FROM people WHERE id = $1"
	p, err := scanPerson(r.db.QueryRowContext(ctx, q, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, person.ErrPersonNotFound
		}
		return nil, fmt.Errorf("get person: %w", err)
	}
	return p, nil
}

func (r *personRepo) Create(ctx context.Context, p *person.Person) error {
	const q = `INSERT INTO people (
		id, project_id, consultant_id, office_id, current_place_id, origin_place_id,
		external_id, first_name, last_name, patronymic, email, birth_date, sex, age_group,
		primary_phone, phone_numbers, case_status, consent_given, consent_date, registered_at,
		created_at, updated_at
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22)`
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now
	_, err := r.db.ExecContext(ctx, q,
		p.ID, p.ProjectID, p.ConsultantID, p.OfficeID, p.CurrentPlaceID, p.OriginPlaceID,
		p.ExternalID, p.FirstName, p.LastName, p.Patronymic, p.Email, p.BirthDate, p.Sex, p.AgeGroup,
		p.PrimaryPhone, p.PhoneNumbers, p.CaseStatus, p.ConsentGiven, p.ConsentDate, p.RegisteredAt,
		p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return person.ErrExternalIDExists
		}
		return fmt.Errorf("create person: %w", err)
	}
	return nil
}

func (r *personRepo) Update(ctx context.Context, p *person.Person) error {
	const q = `UPDATE people SET
		consultant_id=$2, office_id=$3, current_place_id=$4, origin_place_id=$5,
		external_id=$6, first_name=$7, last_name=$8, patronymic=$9, email=$10,
		birth_date=$11, sex=$12, age_group=$13, primary_phone=$14, phone_numbers=$15,
		case_status=$16, consent_given=$17, consent_date=$18, registered_at=$19, updated_at=$20
	WHERE id=$1`
	p.UpdatedAt = time.Now().UTC()
	res, err := r.db.ExecContext(ctx, q,
		p.ID, p.ConsultantID, p.OfficeID, p.CurrentPlaceID, p.OriginPlaceID,
		p.ExternalID, p.FirstName, p.LastName, p.Patronymic, p.Email,
		p.BirthDate, p.Sex, p.AgeGroup, p.PrimaryPhone, p.PhoneNumbers,
		p.CaseStatus, p.ConsentGiven, p.ConsentDate, p.RegisteredAt, p.UpdatedAt,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return person.ErrExternalIDExists
		}
		return fmt.Errorf("update person: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return person.ErrPersonNotFound
	}
	return nil
}

func (r *personRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM people WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete person: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return person.ErrPersonNotFound
	}
	return nil
}
