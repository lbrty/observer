package auth

import (
	"context"

	domainauth "github.com/lbrty/observer/internal/domain/auth"
)

// LogoutUseCase invalidates a user session.
type LogoutUseCase struct {
	sessionRepo domainauth.SessionRepository
}

// NewLogoutUseCase creates a LogoutUseCase.
func NewLogoutUseCase(sessionRepo domainauth.SessionRepository) *LogoutUseCase {
	return &LogoutUseCase{sessionRepo: sessionRepo}
}

// Execute deletes the session associated with the given refresh token.
func (uc *LogoutUseCase) Execute(ctx context.Context, refreshToken string) error {
	return uc.sessionRepo.DeleteByRefreshToken(ctx, refreshToken)
}
