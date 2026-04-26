package auth

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/lbrty/observer/internal/repository"
)

type redisLoginAttemptStore struct {
	client   *redis.Client
	userRepo repository.UserRepository
}

// NewLoginAttemptStore creates a Redis-backed LoginAttemptStore.
func NewLoginAttemptStore(client *redis.Client, userRepo repository.UserRepository) repository.LoginAttemptStore {
	return &redisLoginAttemptStore{client: client, userRepo: userRepo}
}

const (
	lockoutPrefix   = "login_lockout:"
	attemptPrefix   = "login_attempts:"
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
		if err := s.client.Expire(ctx, attemptKey, 24*time.Hour).Err(); err != nil {
			slog.Warn("set attempt counter ttl", slog.Any("err", err))
		}
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
		// Also persist to PostgreSQL so lockout survives Redis restarts.
		if err := s.userRepo.LockPermanently(ctx, email); err != nil {
			slog.Error("persist permanent lock to db", slog.Any("err", err))
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

	switch ttl {
	case -2:
		// Key doesn't exist — not locked.
		return 0, nil
	case -1:
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
