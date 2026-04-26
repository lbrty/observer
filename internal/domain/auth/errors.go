package auth

import "errors"

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
	ErrInvalidMFACode  = errors.New("invalid MFA code")
)
