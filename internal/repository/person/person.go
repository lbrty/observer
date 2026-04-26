package person

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/lbrty/observer/internal/domain/person"
	"github.com/lbrty/observer/internal/repository"
)

type personRepo struct {
	db *sqlx.DB
}

// New creates a PersonRepository.
func New(db *sqlx.DB) repository.PersonRepository {
	return &personRepo{db: db}
}

const personColumns = `id, project_id, consultant_id, office_id, current_place_id, origin_place_id,
	external_id, first_name, last_name, patronymic, email, birth_date, sex, age_group,
	primary_phone, phone_numbers, case_status, consent_given, consent_date, registered_at,
	created_at, updated_at`

// personColumnsTagged is the same column list with explicit "people." table qualifier,
// used when the query joins other tables (e.g. person_tags) to avoid ambiguity.
const personColumnsTagged = `people.id, people.project_id, people.consultant_id, people.office_id, people.current_place_id, people.origin_place_id,
	people.external_id, people.first_name, people.last_name, people.patronymic, people.email, people.birth_date, people.sex, people.age_group,
	people.primary_phone, people.phone_numbers, people.case_status, people.consent_given, people.consent_date, people.registered_at,
	people.created_at, people.updated_at`

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
	repository.TimesToUTC(&p.CreatedAt, &p.UpdatedAt)
	return &p, nil
}

func (r *personRepo) List(ctx context.Context, filter person.PersonListFilter) ([]*person.Person, int, error) {
	cond := sq.And{sq.Eq{"project_id": filter.ProjectID}}

	if filter.ConsultantID != nil {
		cond = append(cond, sq.Eq{"consultant_id": *filter.ConsultantID})
	}
	if filter.OfficeID != nil {
		cond = append(cond, sq.Eq{"office_id": *filter.OfficeID})
	}
	if filter.CaseStatus != nil {
		cond = append(cond, sq.Eq{"case_status": string(*filter.CaseStatus)})
	}
	if filter.Sex != nil {
		cond = append(cond, sq.Eq{"sex": string(*filter.Sex)})
	}
	if filter.AgeGroup != nil {
		cond = append(cond, sq.Expr(
			"(age_group::text = ? OR (age_group IS NULL AND birth_date IS NOT NULL AND infer_age_group(birth_date) = ?))",
			string(*filter.AgeGroup), string(*filter.AgeGroup),
		))
	}
	if filter.CategoryID != nil {
		cond = append(cond, sq.Expr("id IN (SELECT person_id FROM person_categories WHERE category_id = ?)", *filter.CategoryID))
	}
	if filter.RegionID != nil {
		cond = append(cond, sq.Expr("current_place_id IN (SELECT id FROM places WHERE state_id = ?)", *filter.RegionID))
	}
	if filter.HasPets != nil {
		if *filter.HasPets {
			cond = append(cond, sq.Expr("EXISTS (SELECT 1 FROM pets WHERE pets.owner_id = people.id)"))
		} else {
			cond = append(cond, sq.Expr("NOT EXISTS (SELECT 1 FROM pets WHERE pets.owner_id = people.id)"))
		}
	}
	if filter.Search != nil && *filter.Search != "" {
		cond = append(cond, sq.Expr("(first_name % ? OR last_name % ?)", *filter.Search, *filter.Search))
	}

	hasTags := len(filter.TagIDs) > 0
	havingClause := fmt.Sprintf("COUNT(DISTINCT pt.tag_id) = %d", len(filter.TagIDs))

	var countSQL string
	var countArgs []any
	var err error
	if hasTags {
		sub := repository.PSQL.Select("people.id").
			From("people").
			Join("person_tags pt ON pt.person_id = people.id").
			Where(cond).
			Where(sq.Eq{"pt.tag_id": filter.TagIDs}).
			GroupBy("people.id").
			Having(havingClause)
		countSQL, countArgs, err = repository.PSQL.Select("COUNT(*)").FromSelect(sub, "sub").ToSql()
	} else {
		countSQL, countArgs, err = repository.PSQL.Select("COUNT(*)").From("people").Where(cond).ToSql()
	}
	if err != nil {
		return nil, 0, fmt.Errorf("build count query: %w", err)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count people: %w", err)
	}

	perPage, offset := repository.NormalizePagination(filter.Page, filter.PerPage)

	var listSQL string
	var listArgs []any
	if hasTags {
		listSQL, listArgs, err = repository.PSQL.Select(personColumnsTagged).
			From("people").
			Join("person_tags pt ON pt.person_id = people.id").
			Where(cond).
			Where(sq.Eq{"pt.tag_id": filter.TagIDs}).
			GroupBy("people.id").
			Having(havingClause).
			OrderBy("people.created_at DESC").
			Limit(uint64(perPage)).
			Offset(uint64(offset)).
			ToSql()
	} else {
		listSQL, listArgs, err = repository.PSQL.Select(personColumns).
			From("people").
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
	const q = `
		INSERT INTO people (
			id, project_id, consultant_id, office_id, current_place_id, origin_place_id,
			external_id, first_name, last_name, patronymic, email, birth_date, sex, age_group,
			primary_phone, phone_numbers, case_status, consent_given, consent_date, registered_at,
			created_at, updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22)
	`

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
		if repository.IsUniqueViolation(err) {
			return person.ErrExternalIDExists
		}
		return fmt.Errorf("create person: %w", err)
	}
	return nil
}

func (r *personRepo) Update(ctx context.Context, p *person.Person) error {
	const q = `
		UPDATE people
		SET
			consultant_id=$2, office_id=$3, current_place_id=$4, origin_place_id=$5,
			external_id=$6, first_name=$7, last_name=$8, patronymic=$9, email=$10,
			birth_date=$11, sex=$12, age_group=$13, primary_phone=$14, phone_numbers=$15,
			case_status=$16, consent_given=$17, consent_date=$18, registered_at=$19, updated_at=$20
		WHERE id=$1
	`

	p.UpdatedAt = time.Now().UTC()
	res, err := r.db.ExecContext(ctx, q,
		p.ID, p.ConsultantID, p.OfficeID, p.CurrentPlaceID, p.OriginPlaceID,
		p.ExternalID, p.FirstName, p.LastName, p.Patronymic, p.Email,
		p.BirthDate, p.Sex, p.AgeGroup, p.PrimaryPhone, p.PhoneNumbers,
		p.CaseStatus, p.ConsentGiven, p.ConsentDate, p.RegisteredAt, p.UpdatedAt,
	)

	if err != nil {
		if repository.IsUniqueViolation(err) {
			return person.ErrExternalIDExists
		}
		return fmt.Errorf("update person: %w", err)
	}

	return repository.CheckRowsAffected(res, person.ErrPersonNotFound)
}

func (r *personRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM people WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete person: %w", err)
	}
	return repository.CheckRowsAffected(res, person.ErrPersonNotFound)
}
