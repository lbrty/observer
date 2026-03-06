package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

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
