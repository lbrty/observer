package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// LoginAttemptStore tracks failed login attempts in Redis with escalating lockouts.
type LoginAttemptStore interface {
	// RecordFailure increments the failure count and applies a lockout if threshold is reached.
	// Returns the lockout duration if the account is now locked, or 0 if not yet locked.
	RecordFailure(ctx context.Context, email string) (time.Duration, error)
	// IsLocked checks whether the given email is currently locked out.
	// Returns the remaining lockout duration, or 0 if not locked.
	IsLocked(ctx context.Context, email string) (time.Duration, error)
	// ClearAttempts removes all failure tracking for the given email (on successful login or admin unlock).
	ClearAttempts(ctx context.Context, email string) error
}

type redisLoginAttemptStore struct {
	client *redis.Client
}

// NewLoginAttemptStore creates a Redis-backed LoginAttemptStore.
func NewLoginAttemptStore(client *redis.Client) LoginAttemptStore {
	return &redisLoginAttemptStore{client: client}
}

const (
	lockoutPrefix  = "login_lockout:"
	attemptPrefix  = "login_attempts:"
	maxFreeAttempts = 5
)

// Escalating lockout durations indexed by lockout tier (0-based).
var lockoutTiers = []time.Duration{
	5 * time.Minute,
	15 * time.Minute,
	1 * time.Hour,
	0, // 0 means permanent — no TTL, requires admin unlock
}

func (s *redisLoginAttemptStore) RecordFailure(ctx context.Context, email string) (time.Duration, error) {
	attemptKey := attemptPrefix + email
	lockKey := lockoutPrefix + email

	count, err := s.client.Incr(ctx, attemptKey).Result()
	if err != nil {
		return 0, fmt.Errorf("incr login attempts: %w", err)
	}

	// Keep the attempt counter around for 24h so it resets naturally if no further failures.
	if count == 1 {
		s.client.Expire(ctx, attemptKey, 24*time.Hour)
	}

	if count < int64(maxFreeAttempts) {
		return 0, nil
	}

	// Determine lockout tier based on how many times we've hit the threshold.
	tierIndex := int(count)/maxFreeAttempts - 1
	if tierIndex >= len(lockoutTiers) {
		tierIndex = len(lockoutTiers) - 1
	}

	lockDuration := lockoutTiers[tierIndex]

	if lockDuration == 0 {
		// Permanent lock — no expiry.
		if err := s.client.Set(ctx, lockKey, "permanent", 0).Err(); err != nil {
			return 0, fmt.Errorf("set permanent lock: %w", err)
		}
		return -1, nil // -1 signals permanent
	}

	if err := s.client.Set(ctx, lockKey, "locked", lockDuration).Err(); err != nil {
		return 0, fmt.Errorf("set lockout: %w", err)
	}

	return lockDuration, nil
}

func (s *redisLoginAttemptStore) IsLocked(ctx context.Context, email string) (time.Duration, error) {
	lockKey := lockoutPrefix + email

	ttl, err := s.client.TTL(ctx, lockKey).Result()
	if err != nil {
		return 0, fmt.Errorf("check lockout ttl: %w", err)
	}

	switch {
	case ttl == -2:
		// Key doesn't exist — not locked.
		return 0, nil
	case ttl == -1:
		// Key exists with no expiry — permanent lock.
		return -1, nil
	default:
		return ttl, nil
	}
}

func (s *redisLoginAttemptStore) ClearAttempts(ctx context.Context, email string) error {
	pipe := s.client.Pipeline()
	pipe.Del(ctx, attemptPrefix+email)
	pipe.Del(ctx, lockoutPrefix+email)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("clear login attempts: %w", err)
	}
	return nil
}
