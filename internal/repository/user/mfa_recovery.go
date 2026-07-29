package user

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"

	domainuser "github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/repository"
)

type mfaRecoveryCodeRepo struct {
	db *sqlx.DB
}

// NewMFARecovery creates an MFARecoveryCodeRepository backed by the given DB.
func NewMFARecovery(db *sqlx.DB) repository.MFARecoveryCodeRepository {
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

func (r *mfaRecoveryCodeRepo) ConsumeUnused(ctx context.Context, userID ulid.ULID, codeHash string) error {
	const q = `
		UPDATE mfa_recovery_codes
		SET used_at = NOW()
		WHERE id = (
			SELECT id
			FROM mfa_recovery_codes
			WHERE user_id = $1 AND code_hash = $2 AND used_at IS NULL
			LIMIT 1
			FOR UPDATE
		)
	`
	result, err := r.db.ExecContext(ctx, q, userID.String(), codeHash)
	if err != nil {
		return fmt.Errorf("consume recovery code: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("recovery code rows affected: %w", err)
	}
	if affected == 0 {
		return domainuser.ErrMFARecoveryCodeNotFound
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
