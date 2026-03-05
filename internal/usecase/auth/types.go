package auth

import "time"

// RegisterInput holds data for user registration.
type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"required"`
}

// RegisterOutput is the response after a successful registration.
type RegisterOutput struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

// LoginInput holds credentials for login.
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginOutput is the response after a successful login.
type LoginOutput struct {
	RequiresMFA bool       `json:"requires_mfa"`
	MFAToken    string     `json:"mfa_token,omitempty"`
	Tokens      *TokenPair `json:"tokens,omitempty"`
	User        *UserDTO   `json:"user,omitempty"`
}

// TokenPair holds an access + refresh token pair.
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// RefreshTokenInput carries the refresh token for renewal.
type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LogoutInput carries the refresh token to invalidate.
type LogoutInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// UserDTO is a serializable user representation.
type UserDTO struct {
	ID         string    `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Role       string    `json:"role"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
}

// UpdateProfileInput holds fields a user can update on their own profile.
type UpdateProfileInput struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Phone     *string `json:"phone"`
}

// ChangePasswordInput holds data for a password change.
type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}
