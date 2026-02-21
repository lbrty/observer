package user

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailExists        = errors.New("email already exists")
	ErrPhoneExists        = errors.New("phone already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotActive      = errors.New("user account is not active")
	ErrInvalidRole        = errors.New("invalid role")
)
