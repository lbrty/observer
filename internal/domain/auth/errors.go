package auth

import "errors"

var (
	ErrSessionNotFound          = errors.New("session not found")
	ErrSessionExpired           = errors.New("session expired")
	ErrInvalidMFACode           = errors.New("invalid MFA code")
	ErrAccountLocked            = errors.New("account locked, contact administrator")
	ErrAccountTemporarilyLocked = errors.New("account temporarily locked, please try again later")
	ErrRateLimiterUnavailable   = errors.New("rate limiter unavailable, try again later")
)
