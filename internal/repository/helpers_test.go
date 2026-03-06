package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestTimesToUTC(t *testing.T) {
	loc, _ := time.LoadLocation("America/New_York")
	a := time.Now().In(loc)
	b := time.Now().In(loc)
	TimesToUTC(&a, &b)
	assert.Equal(t, time.UTC, a.Location())
	assert.Equal(t, time.UTC, b.Location())
}

type mockResult struct{ rows int64 }

func (m mockResult) LastInsertId() (int64, error) { return 0, nil }
func (m mockResult) RowsAffected() (int64, error) { return m.rows, nil }

var errNotFound = errors.New("not found")

func TestCheckRowsAffected_Found(t *testing.T) {
	err := CheckRowsAffected(mockResult{rows: 1}, errNotFound)
	assert.NoError(t, err)
}

func TestCheckRowsAffected_NotFound(t *testing.T) {
	err := CheckRowsAffected(mockResult{rows: 0}, errNotFound)
	assert.ErrorIs(t, err, errNotFound)
}

func TestIsUniqueViolation_True(t *testing.T) {
	pqErr := &pq.Error{Code: "23505"}
	assert.True(t, IsUniqueViolation(pqErr))
}

func TestIsUniqueViolation_False(t *testing.T) {
	assert.False(t, IsUniqueViolation(errors.New("other")))
}
