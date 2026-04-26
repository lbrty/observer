package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// PSQL is a PostgreSQL squirrel builder (uses $N placeholders).
var PSQL = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// TimesToUTC converts one or more time pointers to UTC in-place.
func TimesToUTC(times ...*time.Time) {
	for _, t := range times {
		*t = (*t).UTC()
	}
}

// CheckRowsAffected returns notFoundErr if no rows were affected.
func CheckRowsAffected(res sql.Result, notFoundErr error) error {
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return notFoundErr
	}
	return nil
}

// IsUniqueViolation checks if err is a PostgreSQL unique constraint violation (code 23505).
func IsUniqueViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}

// BuildBulkTagQuery builds a SELECT for fetching tag IDs for multiple entities.
func BuildBulkTagQuery(table, fkCol string, entityIDs []string) (string, []any, error) {
	return PSQL.
		Select(fkCol + ", tag_id").
		From(table).
		Where(sq.Eq{fkCol: entityIDs}).
		OrderBy(fkCol + ", tag_id").
		ToSql()
}

// NormalizePagination clamps page and perPage to sensible minimums (1 and 20 respectively)
// and returns the normalized perPage and the SQL OFFSET for use in paginated queries.
func NormalizePagination(page, perPage int) (int, int) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	return perPage, (page - 1) * perPage
}

// ReplaceAllJunction atomically replaces all rows for a parent entity in a junction table.
// It deletes every existing row where parentCol = parentID, then bulk-inserts the new childIDs.
// If childIDs is empty the delete still runs, leaving the parent with no associations.
// The whole operation runs in a single transaction so no partial state is ever visible.
func ReplaceAllJunction(ctx context.Context, db *sqlx.DB, table, parentCol, childCol, parentID string, childIDs []string) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "DELETE FROM "+table+" WHERE "+parentCol+" = $1", parentID); err != nil {
		return fmt.Errorf("delete %s: %w", table, err)
	}

	if len(childIDs) > 0 {
		iq := PSQL.Insert(table).Columns(parentCol, childCol)
		for _, id := range childIDs {
			iq = iq.Values(parentID, id)
		}
		sqlStr, args, err := iq.ToSql()
		if err != nil {
			return fmt.Errorf("build insert %s: %w", table, err)
		}
		if _, err := tx.ExecContext(ctx, sqlStr, args...); err != nil {
			return fmt.Errorf("insert %s: %w", table, err)
		}
	}

	return tx.Commit()
}

// QueryBulkTags executes a bulk tag query and returns a map of entity ID → tag IDs.
func QueryBulkTags(ctx context.Context, db *sqlx.DB, q string, args []any) (map[string][]string, error) {
	rows, err := db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("bulk tag query: %w", err)
	}
	defer rows.Close()

	result := make(map[string][]string)
	for rows.Next() {
		var entityID, tagID string
		if err := rows.Scan(&entityID, &tagID); err != nil {
			return nil, fmt.Errorf("scan bulk tag: %w", err)
		}
		result[entityID] = append(result[entityID], tagID)
	}

	return result, rows.Err()
}
