package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"

	"github.com/lbrty/observer/internal/domain/user"
)

// userRepo is a PostgreSQL-backed user repository.
type userRepo struct {
	db *sqlx.DB
}

// NewUserRepository creates a UserRepository backed by the given DB.
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepo{db: db}
}

const userColumns = `id, first_name, last_name, email, phone, office_id, role, is_verified, is_active, created_at, updated_at`

func (r *userRepo) Create(ctx context.Context, u *user.User) error {
	const q = `
		INSERT INTO users (id, first_name, last_name, email, phone, office_id, role, is_verified, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, q,
		u.ID.String(), u.FirstName, u.LastName, u.Email, u.Phone, u.OfficeID,
		string(u.Role), u.IsVerified, u.IsActive, u.CreatedAt.UTC(), u.UpdatedAt.UTC(),
	)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id ulid.ULID) (*user.User, error) {
	q := `SELECT ` + userColumns + ` FROM users WHERE id = $1`
	return r.scanUser(r.db.QueryRowContext(ctx, q, id.String()))
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	q := `SELECT ` + userColumns + ` FROM users WHERE email = $1`
	return r.scanUser(r.db.QueryRowContext(ctx, q, email))
}

func (r *userRepo) GetByPhone(ctx context.Context, phone string) (*user.User, error) {
	q := `SELECT ` + userColumns + ` FROM users WHERE phone = $1`
	return r.scanUser(r.db.QueryRowContext(ctx, q, phone))
}

func (r *userRepo) Update(ctx context.Context, u *user.User) error {
	const q = `
		UPDATE users
		SET first_name=$2, last_name=$3, email=$4, phone=$5, office_id=$6,
		    role=$7, is_verified=$8, is_active=$9, updated_at=$10
		WHERE id=$1
	`
	res, err := r.db.ExecContext(ctx, q,
		u.ID.String(), u.FirstName, u.LastName, u.Email, u.Phone, u.OfficeID,
		string(u.Role), u.IsVerified, u.IsActive, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return r.checkRowsAffected(res, user.ErrUserNotFound)
}

func (r *userRepo) UpdateVerified(ctx context.Context, id ulid.ULID, verified bool) error {
	const q = `UPDATE users SET is_verified=$2, updated_at=$3 WHERE id=$1`
	res, err := r.db.ExecContext(ctx, q, id.String(), verified, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("update verified: %w", err)
	}
	return r.checkRowsAffected(res, user.ErrUserNotFound)
}

func (r *userRepo) List(ctx context.Context, filter user.UserListFilter) ([]*user.User, int, error) {
	var where []string
	var args []any
	argN := 0

	nextArg := func(v any) string {
		argN++
		args = append(args, v)
		return fmt.Sprintf("$%d", argN)
	}

	if filter.Search != "" {
		p := nextArg("%" + filter.Search + "%")
		where = append(where, fmt.Sprintf(
			"(first_name ILIKE %s OR last_name ILIKE %s OR email ILIKE %s)", p, p, p,
		))
	}
	if filter.Role != "" {
		where = append(where, fmt.Sprintf("role = %s", nextArg(filter.Role)))
	}
	if filter.IsActive != nil {
		where = append(where, fmt.Sprintf("is_active = %s", nextArg(*filter.IsActive)))
	}

	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " AND ")
	}

	countQ := "SELECT COUNT(*) FROM users " + whereClause
	var total int
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 {
		filter.PerPage = 20
	}
	offset := (filter.Page - 1) * filter.PerPage

	listQ := fmt.Sprintf("SELECT %s FROM users %s ORDER BY created_at DESC LIMIT %s OFFSET %s",
		userColumns, whereClause, nextArg(filter.PerPage), nextArg(offset))

	rows, err := r.db.QueryContext(ctx, listQ, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	users, err := r.scanUsers(rows)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userRepo) scanUser(row *sql.Row) (*user.User, error) {
	var u user.User
	var idStr string
	var role string
	var createdAt, updatedAt time.Time

	err := row.Scan(&idStr, &u.FirstName, &u.LastName, &u.Email, &u.Phone, &u.OfficeID,
		&role, &u.IsVerified, &u.IsActive, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("scan user: %w", err)
	}

	id, err := ulid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}

	u.ID = id
	u.Role = user.Role(role)
	u.CreatedAt = createdAt.UTC()
	u.UpdatedAt = updatedAt.UTC()

	return &u, nil
}

func (r *userRepo) scanUsers(rows *sql.Rows) ([]*user.User, error) {
	var users []*user.User
	for rows.Next() {
		var u user.User
		var idStr string
		var role string
		var createdAt, updatedAt time.Time

		err := rows.Scan(&idStr, &u.FirstName, &u.LastName, &u.Email, &u.Phone, &u.OfficeID,
			&role, &u.IsVerified, &u.IsActive, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan user row: %w", err)
		}

		id, err := ulid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("parse user id: %w", err)
		}

		u.ID = id
		u.Role = user.Role(role)
		u.CreatedAt = createdAt.UTC()
		u.UpdatedAt = updatedAt.UTC()
		users = append(users, &u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate user rows: %w", err)
	}
	return users, nil
}

func (r *userRepo) checkRowsAffected(res sql.Result, notFoundErr error) error {
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return notFoundErr
	}
	return nil
}
