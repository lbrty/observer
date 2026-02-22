package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/lbrty/observer/internal/crypto"
	domainauth "github.com/lbrty/observer/internal/domain/auth"
	"github.com/lbrty/observer/internal/ulid"
)

// RefreshTokenUseCase issues a new token pair from a valid refresh token.
type RefreshTokenUseCase struct {
	sessionRepo domainauth.SessionRepository
	tokenGen    crypto.TokenGenerator
}

// NewRefreshTokenUseCase creates a RefreshTokenUseCase.
func NewRefreshTokenUseCase(sessionRepo domainauth.SessionRepository, tokenGen crypto.TokenGenerator) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{sessionRepo: sessionRepo, tokenGen: tokenGen}
}

// Execute rotates the refresh token and issues a new access token.
func (uc *RefreshTokenUseCase) Execute(ctx context.Context, input RefreshTokenInput) (*TokenPair, error) {
	session, err := uc.sessionRepo.GetByRefreshToken(ctx, input.RefreshToken)
	if err != nil {
		return nil, domainauth.ErrSessionNotFound
	}

	if time.Now().UTC().After(session.ExpiresAt) {
		_ = uc.sessionRepo.Delete(ctx, session.ID)
		return nil, domainauth.ErrSessionExpired
	}

	if err := uc.sessionRepo.Delete(ctx, session.ID); err != nil {
		return nil, fmt.Errorf("delete old session: %w", err)
	}

	accessToken, expiresAt, err := uc.tokenGen.GenerateAccessToken(session.UserID, "")
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	newRefreshToken, err := uc.tokenGen.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	newSession := &domainauth.Session{
		ID:           ulid.New(),
		UserID:       session.UserID,
		RefreshToken: newRefreshToken,
		UserAgent:    session.UserAgent,
		IP:           session.IP,
		ExpiresAt:    time.Now().UTC().Add(7 * 24 * time.Hour),
		CreatedAt:    time.Now().UTC(),
	}

	if err := uc.sessionRepo.Create(ctx, newSession); err != nil {
		return nil, fmt.Errorf("create new session: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}
