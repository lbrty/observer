package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"

	domainuser "github.com/lbrty/observer/internal/domain/user"
)

type mfaRecoveryCodeRepo struct {
	db *sqlx.DB
}

// NewMFARecoveryCodeRepository creates an MFARecoveryCodeRepository backed by the given DB.
func NewMFARecoveryCodeRepository(db *sqlx.DB) MFARecoveryCodeRepository {
	return &mfaRecoveryCodeRepo{db: db}
}

func (r *mfaRecoveryCodeRepo) CreateBatch(ctx context.Context, codes []*domainuser.MFARecoveryCode) error {
	const q = `
		INSERT INTO mfa_recovery_codes (id, user_id, code_hash, created_at)
		VALUES ($1, $2, $3, $4)
	`
	for _, c := range codes {
		if _, err := r.db.ExecContext(ctx, q, c.ID.String(), c.UserID.String(), c.CodeHash, c.CreatedAt); err != nil {
			return fmt.Errorf("create recovery code: %w", err)
		}
	}
	return nil
}

func (r *mfaRecoveryCodeRepo) FindUnused(ctx context.Context, userID ulid.ULID, codeHash string) (*domainuser.MFARecoveryCode, error) {
	const q = `
		SELECT id, user_id, code_hash, used_at, created_at
		FROM mfa_recovery_codes
		WHERE user_id = $1 AND code_hash = $2 AND used_at IS NULL
		LIMIT 1
	`
	var row struct {
		ID        string     `db:"id"`
		UserID    string     `db:"user_id"`
		CodeHash  string     `db:"code_hash"`
		UsedAt    *time.Time `db:"used_at"`
		CreatedAt time.Time  `db:"created_at"`
	}
	if err := r.db.GetContext(ctx, &row, q, userID.String(), codeHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainuser.ErrMFARecoveryCodeNotFound
		}
		return nil, fmt.Errorf("find recovery code: %w", err)
	}
	uid, err := ulid.Parse(row.UserID)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}
	parsedID, err := ulid.Parse(row.ID)
	if err != nil {
		return nil, fmt.Errorf("parse code id: %w", err)
	}
	return &domainuser.MFARecoveryCode{
		ID:        parsedID,
		UserID:    uid,
		CodeHash:  row.CodeHash,
		UsedAt:    row.UsedAt,
		CreatedAt: row.CreatedAt,
	}, nil
}

func (r *mfaRecoveryCodeRepo) MarkUsed(ctx context.Context, id ulid.ULID) error {
	const q = `UPDATE mfa_recovery_codes SET used_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, q, id.String())
	if err != nil {
		return fmt.Errorf("mark recovery code used: %w", err)
	}
	return nil
}

func (r *mfaRecoveryCodeRepo) DeleteByUserID(ctx context.Context, userID ulid.ULID) error {
	const q = `DELETE FROM mfa_recovery_codes WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, q, userID.String())
	if err != nil {
		return fmt.Errorf("delete recovery codes: %w", err)
	}
	return nil
}
