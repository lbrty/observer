# Seed Command Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add a `observer seed` CLI command that truncates the database and populates it with realistic mock data for development.

**Architecture:** Single `cmd/observer/cmd/seed.go` cobra command that connects to the DB, truncates all tables in FK-safe order via raw SQL, then uses repository interfaces to insert reference data, users, projects, people, and related records. A `gofakeit.Faker` instance with optional seed provides reproducible randomization.

**Tech Stack:** Go, cobra CLI, gofakeit v7, sqlx (raw TRUNCATE), repository layer for inserts, argon hasher for passwords.


### Task 2: Create seed command

**Files:**

- Create: `cmd/observer/cmd/seed.go`
- Modify: `cmd/observer/main.go` (add `rootCmd.AddCommand(cmd.SeedCmd)`)

**Step 1: Create `cmd/observer/cmd/seed.go`**

The file implements the entire seed command. Structure:

```go
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"github.com/lbrty/observer/internal/config"
	"github.com/lbrty/observer/internal/crypto"
	"github.com/lbrty/observer/internal/database"
	"github.com/lbrty/observer/internal/domain/household"
	"github.com/lbrty/observer/internal/domain/migration"
	"github.com/lbrty/observer/internal/domain/note"
	"github.com/lbrty/observer/internal/domain/person"
	"github.com/lbrty/observer/internal/domain/pet"
	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/domain/reference"
	"github.com/lbrty/observer/internal/domain/support"
	"github.com/lbrty/observer/internal/domain/tag"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
)
```

**SeedCmd declaration:**

```go
var SeedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed database with mock data (destructive — truncates all tables first)",
	RunE:  runSeed,
}

func init() {
	SeedCmd.Flags().Int("people", 50, "Number of people per project")
	SeedCmd.Flags().Int("projects", 2, "Number of projects")
	SeedCmd.Flags().Int64("seed", 0, "Random seed (0 = random)")
}
```

**runSeed function — high-level flow:**

```go
func runSeed(cmd *cobra.Command, _ []string) error {
	peopleCount, _ := cmd.Flags().GetInt("people")
	projectCount, _ := cmd.Flags().GetInt("projects")
	seed, _ := cmd.Flags().GetInt64("seed")

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

	faker := gofakeit.New(uint64(seed))

	ctx := context.Background()

	fmt.Println("Truncating all tables...")
	if err := truncateAll(ctx, sqlxDB); err != nil {
		return fmt.Errorf("truncate: %w", err)
	}

	fmt.Println("Seeding reference data...")
	countries, states, places, offices, categories := seedReferenceData(ctx, sqlxDB)

	fmt.Println("Seeding users...")
	users := seedUsers(ctx, sqlxDB)

	fmt.Printf("Seeding %d projects...\n", projectCount)
	for i := range projectCount {
		proj := seedProject(ctx, sqlxDB, faker, users, i)
		seedPermissions(ctx, sqlxDB, proj.ID, users)
		tags := seedTags(ctx, sqlxDB, faker, proj.ID)

		fmt.Printf("  Project %q: seeding %d people...\n", proj.Name, peopleCount)
		people := seedPeople(ctx, sqlxDB, faker, proj.ID, users, places, offices, peopleCount)
		seedPersonCategories(ctx, sqlxDB, faker, people, categories)
		seedPersonTags(ctx, sqlxDB, faker, people, tags)
		seedSupportRecords(ctx, sqlxDB, faker, proj.ID, people, users, offices)
		seedMigrationRecords(ctx, sqlxDB, faker, people, places)
		seedNotes(ctx, sqlxDB, faker, people, users)
		seedPets(ctx, sqlxDB, faker, proj.ID, people)
		seedHouseholds(ctx, sqlxDB, faker, proj.ID, people)
	}

	fmt.Println("Seed complete.")
	return nil
}
```

**truncateAll — raw SQL TRUNCATE CASCADE:**

```go
func truncateAll(ctx context.Context, db *sqlx.DB) error {
	tables := []string{
		"household_members", "households",
		"pets", "documents", "person_notes",
		"migration_records", "support_records",
		"person_categories", "person_tags", "people",
		"tags", "project_permissions", "projects",
		"offices", "places", "states", "countries", "categories",
		"credentials", "sessions", "mfa_configs", "verification_tokens",
		"users",
	}
	for _, t := range tables {
		if _, err := db.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", t)); err != nil {
			return fmt.Errorf("truncate %s: %w", t, err)
		}
	}
	return nil
}
```

**seedReferenceData — hardcoded countries/states/places/offices/categories:**

Uses repository.NewCountryRepository(sqlxDB), etc. Returns slices of created entities for later use.

Countries, states, places as described in design. Each entity gets `ID: ulid.NewString()`, `CreatedAt/UpdatedAt: time.Now().UTC()`.

Categories (12):

```go
categoryNames := []string{
	"Internally Displaced Person", "Refugee", "Asylum Seeker", "Returnee",
	"Stateless Person", "Conflict-Affected Civilian", "Unaccompanied Minor",
	"Single-Headed Household", "Person with Disability", "Elderly at Risk",
	"GBV Survivor", "Chronic Illness",
}
```

**seedUsers — 4 fixed users:**

```go
type seedUser struct {
	email string
	role  user.Role
	first string
	last  string
	phone string
}

seedUsers := []seedUser{
	{"admin@example.com", user.RoleAdmin, "Admin", "User", "+10000000001"},
	{"staff@example.com", user.RoleStaff, "Staff", "User", "+10000000002"},
	{"consultant@example.com", user.RoleConsultant, "Consultant", "User", "+10000000003"},
	{"guest@example.com", user.RoleGuest, "Guest", "User", "+10000000004"},
}
```

All get password `"password"`, hashed with `crypto.NewArgonHasher()`. Returns `[]*user.User`.

**seedProject — create project with faker name:**

```go
proj := &project.Project{
	ID:      ulid.NewString(),
	Name:    fmt.Sprintf("Project %s", faker.City()),
	OwnerID: adminUser.ID.String(),
	Status:  project.ProjectStatusActive,
	...
}
```

**seedPermissions — assign roles to each project:**

- admin → owner (all flags true)
- staff → manager (all flags true)
- consultant → consultant (can_view_contact=true, others false)
- guest → viewer (all flags false)

**seedTags — 5-8 random tags per project:**

```go
tagPool := []string{"urgent", "follow-up", "legal-aid", "housing", "medical", "employment", "education", "documentation", "family", "vulnerable"}
count := faker.IntRange(5, 8)
```

**seedPeople — N people per project:**
Each person gets:

- `FirstName`, `LastName` from faker
- Random `Sex` from enum values
- Random `AgeGroup` from enum values
- Random `CaseStatus` (weighted: 60% active, 20% new, 10% closed, 10% archived)
- Random `OriginPlaceID` from places slice
- Random `CurrentPlaceID` from places slice
- `PrimaryPhone` from faker
- `Email` from faker
- `ConsultantID` = consultant user's ID (50% chance)
- `OfficeID` = random office (70% chance)
- `ConsentGiven` = true (80% chance), with `ConsentDate` if true
- `RegisteredAt` = random date in past 2 years

**seedPersonCategories — 1-3 random categories per person**

**seedPersonTags — 0-2 random tags per person**

**seedSupportRecords — 0-3 per person:**
Each record:

- Random `Type` from support types
- Random `Sphere` (70% chance)
- `ProvidedAt` = random date after person's `RegisteredAt`
- `Notes` from faker.Sentence() (50% chance)
- `ConsultantID` = consultant user (60% chance)
- `OfficeID` = random office (50% chance)
- `ReferralStatus` (30% chance, random from enum)

**seedMigrationRecords — 0-2 per person:**
Each record:

- `FromPlaceID` = person's origin place
- `DestinationPlaceID` = random place
- `MigrationDate` = random date
- `MovementReason` = random from enum
- `HousingAtDestination` = random from enum
- `Notes` from faker (40% chance)

**seedNotes — 0-2 per person:**

- `Body` = faker.Sentence(faker.IntRange(5, 20))
- `AuthorID` = consultant or staff user

**seedPets — 0-1 per person (30% chance):**

- `Name` = faker.PetName()
- `Status` = random from pet statuses
- `OwnerID` = person ID

**seedHouseholds — ~20% of people grouped:**
Take every 5th person as household head, group with 1-3 following people:

- Create household with `HeadPersonID`
- Add head as `RelationshipHead`
- Add other members with random relationships (spouse, child, parent, sibling)

**Step 2: Register the command in main.go**

Add `rootCmd.AddCommand(cmd.SeedCmd)` in `init()`.

**Step 3: Build and verify**

Run: `go build ./...`

**Step 4: Test manually**

Run: `observer seed --people 10 --projects 1 --seed 42`
Expected output:

```
Truncating all tables...
Seeding reference data...
Seeding users...
Seeding 1 projects...
  Project "Project <city>": seeding 10 people...
Seed complete.
```

Verify with: `psql` → `SELECT count(*) FROM people;` → 10

**Step 5: Commit**

```bash
git add cmd/observer/cmd/seed.go cmd/observer/main.go
git commit -m "add seed command for dev database population"
```

