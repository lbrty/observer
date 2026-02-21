package ulid_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lbrty/observer/internal/ulid"
)

func TestNew_Uniqueness(t *testing.T) {
	const n = 1000
	seen := make(map[string]struct{}, n)
	for i := 0; i < n; i++ {
		id := ulid.NewString()
		_, exists := seen[id]
		assert.False(t, exists, "duplicate ULID: %s", id)
		seen[id] = struct{}{}
	}
}

func TestNew_ThreadSafety(t *testing.T) {
	const goroutines = 50
	const perGoroutine = 100

	ch := make(chan string, goroutines*perGoroutine)
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < perGoroutine; j++ {
				ch <- ulid.NewString()
			}
		}()
	}

	wg.Wait()
	close(ch)

	seen := make(map[string]struct{})
	for id := range ch {
		_, exists := seen[id]
		assert.False(t, exists, "duplicate ULID: %s", id)
		seen[id] = struct{}{}
	}
}

func TestParse_Valid(t *testing.T) {
	original := ulid.New()
	parsed, err := ulid.Parse(original.String())
	require.NoError(t, err)
	assert.Equal(t, original, parsed)
}

func TestParse_Invalid(t *testing.T) {
	_, err := ulid.Parse("not-a-valid-ulid")
	assert.Error(t, err)
}

func TestIsValid(t *testing.T) {
	assert.True(t, ulid.IsValid(ulid.NewString()))
	assert.False(t, ulid.IsValid(""))
	assert.False(t, ulid.IsValid("invalid"))
}
