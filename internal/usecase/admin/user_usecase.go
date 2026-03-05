package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/lbrty/observer/internal/crypto"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/usecase"
	iulid "github.com/lbrty/observer/internal/ulid"
)

// ResetPasswordInput holds the new password for an admin reset.
type ResetPasswordInput struct {
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// UserUseCase handles admin user management operations.
type UserUseCase struct {
	userRepo repository.UserRepository
	credRepo repository.CredentialsRepository
	hasher   crypto.PasswordHasher
}

// NewUserUseCase creates a UserUseCase.
func NewUserUseCase(
	userRepo repository.UserRepository,
	credRepo repository.CredentialsRepository,
	hasher crypto.PasswordHasher,
) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
		credRepo: credRepo,
		hasher:   hasher,
	}
}

// List returns a paginated list of users.
func (uc *UserUseCase) List(ctx context.Context, input ListUsersInput) (*ListUsersOutput, error) {
	input.Page, input.PerPage = usecase.ClampPagination(input.Page, input.PerPage)

	filter := user.UserListFilter{
		Page:     input.Page,
		PerPage:  input.PerPage,
		Search:   input.Search,
		Role:     input.Role,
		IsActive: input.IsActive,
	}

	users, total, err := uc.userRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	dtos := make([]UserDTO, len(users))
	for i, u := range users {
		dtos[i] = userToDTO(u)
	}

	return &ListUsersOutput{
		Users:   dtos,
		Total:   total,
		Page:    input.Page,
		PerPage: input.PerPage,
	}, nil
}

// Get returns a user by ID.
func (uc *UserUseCase) Get(ctx context.Context, id ulid.ULID) (*UserDTO, error) {
	u, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	dto := userToDTO(u)
	return &dto, nil
}

// Create creates a new user with credentials.
func (uc *UserUseCase) Create(ctx context.Context, input CreateUserInput) (*UserDTO, error) {
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

	userID := iulid.New()
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

// Update applies a partial update to the user and returns the updated DTO.
func (uc *UserUseCase) Update(ctx context.Context, id ulid.ULID, input UpdateUserInput) (*UserDTO, error) {
	u, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user for update: %w", err)
	}

	if input.FirstName != nil {
		u.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		u.LastName = *input.LastName
	}
	if input.Email != nil {
		u.Email = *input.Email
	}
	if input.Phone != nil {
		u.Phone = *input.Phone
	}
	if input.OfficeID != nil {
		u.OfficeID = input.OfficeID
	}
	if input.Role != nil {
		role, err := user.ValidateRole(*input.Role)
		if err != nil {
			return nil, err
		}
		u.Role = role
	}
	if input.IsActive != nil {
		u.IsActive = *input.IsActive
	}
	if input.IsVerified != nil {
		u.IsVerified = *input.IsVerified
	}

	if err := uc.userRepo.Update(ctx, u); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	dto := userToDTO(u)
	return &dto, nil
}

// ResetPassword resets the password for the given user.
func (uc *UserUseCase) ResetPassword(ctx context.Context, userID ulid.ULID, input ResetPasswordInput) error {
	cred, err := uc.credRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get credentials: %w", err)
	}

	hash, salt, err := uc.hasher.Hash(input.NewPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	cred.PasswordHash = hash
	cred.Salt = salt
	if err := uc.credRepo.Update(ctx, cred); err != nil {
		return fmt.Errorf("update credentials: %w", err)
	}

	return nil
}

func userToDTO(u *user.User) UserDTO {
	return UserDTO{
		ID:         u.ID.String(),
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		Email:      u.Email,
		Phone:      u.Phone,
		OfficeID:   u.OfficeID,
		Role:       string(u.Role),
		IsVerified: u.IsVerified,
		IsActive:   u.IsActive,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}
