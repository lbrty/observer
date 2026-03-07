package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

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

// buildBulkTagQuery builds a SELECT for fetching tag IDs for multiple entities.
func buildBulkTagQuery(table, fkCol string, entityIDs []string) (string, []any) {
	args := make([]any, len(entityIDs))
	params := make([]string, len(entityIDs))
	for i, id := range entityIDs {
		args[i] = id
		params[i] = fmt.Sprintf("$%d", i+1)
	}
	q := fmt.Sprintf("SELECT %s, tag_id FROM %s WHERE %s IN (%s) ORDER BY %s, tag_id",
		fkCol, table, fkCol, joinStrings(params, ", "), fkCol)
	return q, args
}

// queryBulkTags executes a bulk tag query and returns a map of entity ID → tag IDs.
func queryBulkTags(ctx context.Context, db *sqlx.DB, q string, args []any) (map[string][]string, error) {
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

func joinStrings(ss []string, sep string) string {
	if len(ss) == 0 {
		return ""
	}
	out := ss[0]
	for _, s := range ss[1:] {
		out += sep + s
	}
	return out
}
