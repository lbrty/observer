package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/lbrty/observer/internal/domain/reference"
)

// countryRepo is a PostgreSQL-backed country repository.
type countryRepo struct {
	db *sqlx.DB
}

// NewCountryRepository creates a CountryRepository.
func NewCountryRepository(db *sqlx.DB) CountryRepository {
	return &countryRepo{db: db}
}

func (r *countryRepo) List(ctx context.Context) ([]*reference.Country, error) {
	const q = `SELECT id, name, code, created_at, updated_at FROM countries ORDER BY name`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list countries: %w", err)
	}
	defer rows.Close()

	var out []*reference.Country
	for rows.Next() {
		var c reference.Country
		if err := rows.Scan(&c.ID, &c.Name, &c.Code, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan country: %w", err)
		}
		c.CreatedAt = c.CreatedAt.UTC()
		c.UpdatedAt = c.UpdatedAt.UTC()
		out = append(out, &c)
	}
	return out, rows.Err()
}

func (r *countryRepo) GetByID(ctx context.Context, id string) (*reference.Country, error) {
	const q = `SELECT id, name, code, created_at, updated_at FROM countries WHERE id = $1`
	var c reference.Country
	err := r.db.QueryRowContext(ctx, q, id).Scan(&c.ID, &c.Name, &c.Code, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, reference.ErrCountryNotFound
		}
		return nil, fmt.Errorf("get country: %w", err)
	}
	c.CreatedAt = c.CreatedAt.UTC()
	c.UpdatedAt = c.UpdatedAt.UTC()
	return &c, nil
}

func (r *countryRepo) Create(ctx context.Context, c *reference.Country) error {
	const q = `INSERT INTO countries (id, name, code, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	now := time.Now().UTC()
	c.CreatedAt = now
	c.UpdatedAt = now
	_, err := r.db.ExecContext(ctx, q, c.ID, c.Name, c.Code, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return reference.ErrCountryCodeExists
		}
		return fmt.Errorf("create country: %w", err)
	}
	return nil
}

func (r *countryRepo) Update(ctx context.Context, c *reference.Country) error {
	const q = `UPDATE countries SET name=$2, code=$3, updated_at=$4 WHERE id=$1`
	c.UpdatedAt = time.Now().UTC()
	res, err := r.db.ExecContext(ctx, q, c.ID, c.Name, c.Code, c.UpdatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return reference.ErrCountryCodeExists
		}
		return fmt.Errorf("update country: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return reference.ErrCountryNotFound
	}
	return nil
}

func (r *countryRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM countries WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete country: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return reference.ErrCountryNotFound
	}
	return nil
}

// stateRepo is a PostgreSQL-backed state repository.
type stateRepo struct {
	db *sqlx.DB
}

// NewStateRepository creates a StateRepository.
func NewStateRepository(db *sqlx.DB) StateRepository {
	return &stateRepo{db: db}
}

func (r *stateRepo) ListAll(ctx context.Context) ([]*reference.State, error) {
	const q = `
		SELECT id, country_id, name, code, conflict_zone, created_at, updated_at
		FROM states ORDER BY name
	`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list all states: %w", err)
	}
	defer rows.Close()

	var out []*reference.State
	for rows.Next() {
		var s reference.State
		if err := rows.Scan(&s.ID, &s.CountryID, &s.Name, &s.Code, &s.ConflictZone, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan state: %w", err)
		}
		s.CreatedAt = s.CreatedAt.UTC()
		s.UpdatedAt = s.UpdatedAt.UTC()
		out = append(out, &s)
	}
	return out, rows.Err()
}

func (r *stateRepo) List(ctx context.Context, countryID string) ([]*reference.State, error) {
	const q = `
		SELECT id, country_id, name, code, conflict_zone, created_at, updated_at
		FROM states WHERE country_id = $1 ORDER BY name
	`
	rows, err := r.db.QueryContext(ctx, q, countryID)
	if err != nil {
		return nil, fmt.Errorf("list states: %w", err)
	}
	defer rows.Close()

	var out []*reference.State
	for rows.Next() {
		var s reference.State
		if err := rows.Scan(&s.ID, &s.CountryID, &s.Name, &s.Code, &s.ConflictZone, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan state: %w", err)
		}
		s.CreatedAt = s.CreatedAt.UTC()
		s.UpdatedAt = s.UpdatedAt.UTC()
		out = append(out, &s)
	}
	return out, rows.Err()
}

func (r *stateRepo) GetByID(ctx context.Context, id string) (*reference.State, error) {
	const q = `SELECT id, country_id, name, code, conflict_zone, created_at, updated_at FROM states WHERE id = $1`
	var s reference.State
	err := r.db.QueryRowContext(ctx, q, id).Scan(&s.ID, &s.CountryID, &s.Name, &s.Code, &s.ConflictZone, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, reference.ErrStateNotFound
		}
		return nil, fmt.Errorf("get state: %w", err)
	}
	s.CreatedAt = s.CreatedAt.UTC()
	s.UpdatedAt = s.UpdatedAt.UTC()
	return &s, nil
}

func (r *stateRepo) Create(ctx context.Context, s *reference.State) error {
	const q = `
		INSERT INTO states (id, country_id, name, code, conflict_zone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	now := time.Now().UTC()
	s.CreatedAt = now
	s.UpdatedAt = now
	_, err := r.db.ExecContext(ctx, q, s.ID, s.CountryID, s.Name, s.Code, s.ConflictZone, s.CreatedAt, s.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create state: %w", err)
	}
	return nil
}

func (r *stateRepo) Update(ctx context.Context, s *reference.State) error {
	const q = `UPDATE states SET name=$2, code=$3, conflict_zone=$4, updated_at=$5 WHERE id=$1`
	s.UpdatedAt = time.Now().UTC()
	res, err := r.db.ExecContext(ctx, q, s.ID, s.Name, s.Code, s.ConflictZone, s.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update state: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return reference.ErrStateNotFound
	}
	return nil
}

func (r *stateRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM states WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete state: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return reference.ErrStateNotFound
	}
	return nil
}

// placeRepo is a PostgreSQL-backed place repository.
type placeRepo struct {
	db *sqlx.DB
}

// NewPlaceRepository creates a PlaceRepository.
func NewPlaceRepository(db *sqlx.DB) PlaceRepository {
	return &placeRepo{db: db}
}

func (r *placeRepo) ListAll(ctx context.Context) ([]*reference.Place, error) {
	const q = `
		SELECT id, state_id, name, lat, lon, created_at, updated_at
		FROM places ORDER BY name
	`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list all places: %w", err)
	}
	defer rows.Close()

	var out []*reference.Place
	for rows.Next() {
		var p reference.Place
		if err := rows.Scan(&p.ID, &p.StateID, &p.Name, &p.Lat, &p.Lon, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan place: %w", err)
		}
		p.CreatedAt = p.CreatedAt.UTC()
		p.UpdatedAt = p.UpdatedAt.UTC()
		out = append(out, &p)
	}
	return out, rows.Err()
}

func (r *placeRepo) List(ctx context.Context, stateID string) ([]*reference.Place, error) {
	const q = `
		SELECT id, state_id, name, lat, lon, created_at, updated_at
		FROM places WHERE state_id = $1 ORDER BY name
	`
	rows, err := r.db.QueryContext(ctx, q, stateID)
	if err != nil {
		return nil, fmt.Errorf("list places: %w", err)
	}
	defer rows.Close()

	var out []*reference.Place
	for rows.Next() {
		var p reference.Place
		if err := rows.Scan(&p.ID, &p.StateID, &p.Name, &p.Lat, &p.Lon, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan place: %w", err)
		}
		p.CreatedAt = p.CreatedAt.UTC()
		p.UpdatedAt = p.UpdatedAt.UTC()
		out = append(out, &p)
	}
	return out, rows.Err()
}

func (r *placeRepo) GetByID(ctx context.Context, id string) (*reference.Place, error) {
	const q = `SELECT id, state_id, name, lat, lon, created_at, updated_at FROM places WHERE id = $1`
	var p reference.Place
	err := r.db.QueryRowContext(ctx, q, id).Scan(&p.ID, &p.StateID, &p.Name, &p.Lat, &p.Lon, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, reference.ErrPlaceNotFound
		}
		return nil, fmt.Errorf("get place: %w", err)
	}
	p.CreatedAt = p.CreatedAt.UTC()
	p.UpdatedAt = p.UpdatedAt.UTC()
	return &p, nil
}

func (r *placeRepo) Create(ctx context.Context, p *reference.Place) error {
	const q = `
		INSERT INTO places (id, state_id, name, lat, lon, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now
	_, err := r.db.ExecContext(ctx, q, p.ID, p.StateID, p.Name, p.Lat, p.Lon, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create place: %w", err)
	}
	return nil
}

func (r *placeRepo) Update(ctx context.Context, p *reference.Place) error {
	const q = `UPDATE places SET name=$2, lat=$3, lon=$4, updated_at=$5 WHERE id=$1`
	p.UpdatedAt = time.Now().UTC()
	res, err := r.db.ExecContext(ctx, q, p.ID, p.Name, p.Lat, p.Lon, p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update place: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return reference.ErrPlaceNotFound
	}
	return nil
}

func (r *placeRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM places WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete place: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return reference.ErrPlaceNotFound
	}
	return nil
}

// officeRepo is a PostgreSQL-backed office repository.
type officeRepo struct {
	db *sqlx.DB
}

// NewOfficeRepository creates an OfficeRepository.
func NewOfficeRepository(db *sqlx.DB) OfficeRepository {
	return &officeRepo{db: db}
}

func (r *officeRepo) List(ctx context.Context) ([]*reference.Office, error) {
	const q = `SELECT id, name, place_id, created_at, updated_at FROM offices ORDER BY name`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list offices: %w", err)
	}
	defer rows.Close()

	var out []*reference.Office
	for rows.Next() {
		var o reference.Office
		if err := rows.Scan(&o.ID, &o.Name, &o.PlaceID, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan office: %w", err)
		}
		o.CreatedAt = o.CreatedAt.UTC()
		o.UpdatedAt = o.UpdatedAt.UTC()
		out = append(out, &o)
	}
	return out, rows.Err()
}

func (r *officeRepo) GetByID(ctx context.Context, id string) (*reference.Office, error) {
	const q = `SELECT id, name, place_id, created_at, updated_at FROM offices WHERE id = $1`
	var o reference.Office
	err := r.db.QueryRowContext(ctx, q, id).Scan(&o.ID, &o.Name, &o.PlaceID, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, reference.ErrOfficeNotFound
		}
		return nil, fmt.Errorf("get office: %w", err)
	}
	o.CreatedAt = o.CreatedAt.UTC()
	o.UpdatedAt = o.UpdatedAt.UTC()
	return &o, nil
}

func (r *officeRepo) Create(ctx context.Context, o *reference.Office) error {
	const q = `INSERT INTO offices (id, name, place_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	now := time.Now().UTC()
	o.CreatedAt = now
	o.UpdatedAt = now
	_, err := r.db.ExecContext(ctx, q, o.ID, o.Name, o.PlaceID, o.CreatedAt, o.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create office: %w", err)
	}
	return nil
}

func (r *officeRepo) Update(ctx context.Context, o *reference.Office) error {
	const q = `UPDATE offices SET name=$2, place_id=$3, updated_at=$4 WHERE id=$1`
	o.UpdatedAt = time.Now().UTC()
	res, err := r.db.ExecContext(ctx, q, o.ID, o.Name, o.PlaceID, o.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update office: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return reference.ErrOfficeNotFound
	}
	return nil
}

func (r *officeRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM offices WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete office: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return reference.ErrOfficeNotFound
	}
	return nil
}

// categoryRepo is a PostgreSQL-backed category repository.
type categoryRepo struct {
	db *sqlx.DB
}

// NewCategoryRepository creates a CategoryRepository.
func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) List(ctx context.Context) ([]*reference.Category, error) {
	const q = `SELECT id, name, description, created_at, updated_at FROM categories ORDER BY name`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	defer rows.Close()

	var out []*reference.Category
	for rows.Next() {
		var c reference.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}
		c.CreatedAt = c.CreatedAt.UTC()
		c.UpdatedAt = c.UpdatedAt.UTC()
		out = append(out, &c)
	}
	return out, rows.Err()
}

func (r *categoryRepo) GetByID(ctx context.Context, id string) (*reference.Category, error) {
	const q = `SELECT id, name, description, created_at, updated_at FROM categories WHERE id = $1`
	var c reference.Category
	err := r.db.QueryRowContext(ctx, q, id).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, reference.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("get category: %w", err)
	}
	c.CreatedAt = c.CreatedAt.UTC()
	c.UpdatedAt = c.UpdatedAt.UTC()
	return &c, nil
}

func (r *categoryRepo) Create(ctx context.Context, c *reference.Category) error {
	const q = `INSERT INTO categories (id, name, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	now := time.Now().UTC()
	c.CreatedAt = now
	c.UpdatedAt = now
	_, err := r.db.ExecContext(ctx, q, c.ID, c.Name, c.Description, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return reference.ErrCategoryNameExists
		}
		return fmt.Errorf("create category: %w", err)
	}
	return nil
}

func (r *categoryRepo) Update(ctx context.Context, c *reference.Category) error {
	const q = `UPDATE categories SET name=$2, description=$3, updated_at=$4 WHERE id=$1`
	c.UpdatedAt = time.Now().UTC()
	res, err := r.db.ExecContext(ctx, q, c.ID, c.Name, c.Description, c.UpdatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return reference.ErrCategoryNameExists
		}
		return fmt.Errorf("update category: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return reference.ErrCategoryNotFound
	}
	return nil
}

func (r *categoryRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM categories WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return reference.ErrCategoryNotFound
	}
	return nil
}
