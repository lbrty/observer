package auth

import (
	"context"

	"github.com/lbrty/observer/internal/repository"
)

// LogoutUseCase invalidates a user session.
type LogoutUseCase struct {
	sessionRepo repository.SessionRepository
}

// NewLogoutUseCase creates a LogoutUseCase.
func NewLogoutUseCase(sessionRepo repository.SessionRepository) *LogoutUseCase {
	return &LogoutUseCase{sessionRepo: sessionRepo}
}

// Execute deletes the session associated with the given refresh token.
func (uc *LogoutUseCase) Execute(ctx context.Context, refreshToken string) error {
	return uc.sessionRepo.DeleteByRefreshToken(ctx, refreshToken)
}
