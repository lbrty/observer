package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/lbrty/observer/internal/config"
	"github.com/lbrty/observer/internal/crypto"
	"github.com/lbrty/observer/internal/database"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
)

// CreateAdminCmd creates an admin user from the command line.
var CreateAdminCmd = &cobra.Command{
	Use:   "create-admin",
	Short: "Create an admin user",
	Long: `Create a platform administrator account.

Connects to the database, hashes the password with Argon2id, and inserts
the user with admin role, verified and active. Requires DATABASE_DSN
to be set. Rejects duplicate emails and phone numbers.`,
	Example: `  # Create an admin with required fields
  observer create-admin --email admin@example.com --password "s3cure-p4ss"

  # With optional profile fields
  observer create-admin \
    --email admin@example.com \
    --password "s3cure-p4ss" \
    --first-name Admin \
    --last-name User \
    --phone "+1234567890"`,
	RunE: runCreateAdmin,
}

func init() {
	CreateAdminCmd.Flags().String("email", "", "Admin email (required)")
	CreateAdminCmd.Flags().String("password", "", "Admin password (required, min 8 chars)")
	CreateAdminCmd.Flags().String("first-name", "", "First name")
	CreateAdminCmd.Flags().String("last-name", "", "Last name")
	CreateAdminCmd.Flags().String("phone", "", "Phone number")

	_ = CreateAdminCmd.MarkFlagRequired("email")
	_ = CreateAdminCmd.MarkFlagRequired("password")
}

func runCreateAdmin(cmd *cobra.Command, _ []string) error {
	email, _ := cmd.Flags().GetString("email")
	password, _ := cmd.Flags().GetString("password")
	firstName, _ := cmd.Flags().GetString("first-name")
	lastName, _ := cmd.Flags().GetString("last-name")
	phone, _ := cmd.Flags().GetString("phone")

	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	db, err := database.New(cfg.Database.DSN)
	if err != nil {
		return err
	}
	defer db.Close()

	sqlxDB := db.GetDB()
	userRepo := repository.NewUserRepository(sqlxDB)
	credRepo := repository.NewCredentialsRepository(sqlxDB)
	hasher := crypto.NewArgonHasher()

	ctx := context.Background()

	if _, err := userRepo.GetByEmail(ctx, email); err == nil {
		return fmt.Errorf("user with email %s already exists", email)
	}

	if phone != "" {
		if _, err := userRepo.GetByPhone(ctx, phone); err == nil {
			return fmt.Errorf("user with phone %s already exists", phone)
		}
	}

	hash, salt, err := hasher.Hash(password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	userID := ulid.New()
	now := time.Now().UTC()

	newUser := &user.User{
		ID:         userID,
		FirstName:  firstName,
		LastName:   lastName,
		Email:      email,
		Phone:      phone,
		Role:       user.RoleAdmin,
		IsVerified: true,
		IsActive:   true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := userRepo.Create(ctx, newUser); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	cred := &user.Credentials{
		UserID:       userID,
		PasswordHash: hash,
		Salt:         salt,
		UpdatedAt:    now,
	}

	if err := credRepo.Create(ctx, cred); err != nil {
		return fmt.Errorf("create credentials: %w", err)
	}

	fmt.Printf("Admin user created: %s (%s %s)\n", email, firstName, lastName)
	return nil
}
