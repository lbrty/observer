package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/lbrty/observer/internal/crypto"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
)

// CreateUserUseCase handles admin user creation.
type CreateUserUseCase struct {
	userRepo repository.UserRepository
	credRepo repository.CredentialsRepository
	hasher   crypto.PasswordHasher
}

// NewCreateUserUseCase creates a CreateUserUseCase.
func NewCreateUserUseCase(
	userRepo repository.UserRepository,
	credRepo repository.CredentialsRepository,
	hasher crypto.PasswordHasher,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepo: userRepo,
		credRepo: credRepo,
		hasher:   hasher,
	}
}

// Execute creates a new user with credentials.
func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*UserDTO, error) {
	if _, err := user.ValidateRole(input.Role); err != nil {
		return nil, err
	}

	if _, err := uc.userRepo.GetByEmail(ctx, input.Email); err == nil {
		return nil, user.ErrEmailExists
	}

	if input.Phone != "" {
		if _, err := uc.userRepo.GetByPhone(ctx, input.Phone); err == nil {
			return nil, user.ErrPhoneExists
		}
	}

	hash, salt, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	userID := ulid.New()
	now := time.Now().UTC()

	newUser := &user.User{
		ID:         userID,
		FirstName:  input.FirstName,
		LastName:   input.LastName,
		Email:      input.Email,
		Phone:      input.Phone,
		OfficeID:   input.OfficeID,
		Role:       user.Role(input.Role),
		IsVerified: input.IsVerified,
		IsActive:   input.IsActive,
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

	return &UserDTO{
		ID:         userID.String(),
		FirstName:  newUser.FirstName,
		LastName:   newUser.LastName,
		Email:      newUser.Email,
		Phone:      newUser.Phone,
		OfficeID:   newUser.OfficeID,
		Role:       string(newUser.Role),
		IsVerified: newUser.IsVerified,
		IsActive:   newUser.IsActive,
		CreatedAt:  newUser.CreatedAt,
		UpdatedAt:  newUser.UpdatedAt,
	}, nil
}
