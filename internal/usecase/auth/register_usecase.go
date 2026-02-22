package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/lbrty/observer/internal/crypto"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
)

// RegisterUseCase handles user registration.
type RegisterUseCase struct {
	userRepo repository.UserRepository
	credRepo repository.CredentialsRepository
	hasher   crypto.PasswordHasher
}

// NewRegisterUseCase creates a RegisterUseCase.
func NewRegisterUseCase(
	userRepo repository.UserRepository,
	credRepo repository.CredentialsRepository,
	hasher crypto.PasswordHasher,
) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo: userRepo,
		credRepo: credRepo,
		hasher:   hasher,
	}
}

// Execute registers a new user.
func (uc *RegisterUseCase) Execute(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
	if _, err := user.ValidateRole(input.Role); err != nil {
		return nil, err
	}

	if _, err := uc.userRepo.GetByEmail(ctx, input.Email); err == nil {
		return nil, user.ErrEmailExists
	}

	if _, err := uc.userRepo.GetByPhone(ctx, input.Phone); err == nil {
		return nil, user.ErrPhoneExists
	}

	hash, salt, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	userID := ulid.New()
	now := time.Now().UTC()

	newUser := &user.User{
		ID:         userID,
		Email:      input.Email,
		Phone:      input.Phone,
		Role:       user.Role(input.Role),
		IsVerified: false,
		IsActive:   true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := uc.userRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	cred := &user.Credentials{
		UserID:       userID,
		PasswordHash: hash,
		Salt:         salt,
		UpdatedAt:    now,
	}

	if err := uc.credRepo.Create(ctx, cred); err != nil {
		return nil, fmt.Errorf("create credentials: %w", err)
	}

	return &RegisterOutput{
		UserID:  userID.String(),
		Message: "Registration successful. Please check your email to verify your account.",
	}, nil
}
