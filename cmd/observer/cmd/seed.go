package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"github.com/lbrty/observer/internal/config"
	"github.com/lbrty/observer/internal/crypto"
	"github.com/lbrty/observer/internal/database"
	"github.com/lbrty/observer/internal/domain/household"
	"github.com/lbrty/observer/internal/domain/migration"
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

// SeedCmd populates the database with mock data for development.
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

const batchSize = 500

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
	places, offices, categories := seedReferenceData(ctx, sqlxDB)

	fmt.Println("Seeding users...")
	users := seedUsers(ctx, sqlxDB)

	fmt.Printf("Seeding %d project(s)...\n", projectCount)
	for i := range projectCount {
		proj := seedProject(ctx, sqlxDB, faker, users[0], i)
		seedPermissions(ctx, sqlxDB, proj.ID, users)
		tags := seedTags(ctx, sqlxDB, faker, proj.ID)

		fmt.Printf("  Project %q: seeding %d people...\n", proj.Name, peopleCount)
		people := genPeople(faker, proj.ID, users, places, offices, peopleCount)
		bulkInsertPeople(ctx, sqlxDB, people)
		genAndInsertStatusHistory(ctx, sqlxDB, faker, people)

		genAndInsertPersonCategories(ctx, sqlxDB, faker, people, categories)
		genAndInsertPersonTags(ctx, sqlxDB, faker, people, tags)
		genAndInsertSupportRecords(ctx, sqlxDB, faker, proj.ID, people, users, offices)
		genAndInsertMigrationRecords(ctx, sqlxDB, faker, people, places)
		genAndInsertNotes(ctx, sqlxDB, faker, people, users)
		genAndInsertPets(ctx, sqlxDB, faker, proj.ID, people)
		seedHouseholds(ctx, sqlxDB, faker, proj.ID, people)
	}

	fmt.Println("Seed complete.")
	return nil
}

func truncateAll(ctx context.Context, db *sqlx.DB) error {
	tables := []string{
		"household_members", "households",
		"pets", "documents", "person_notes",
		"migration_records", "support_records",
		"person_status_history",
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

// bulkInsert executes batched multi-row INSERTs.
func bulkInsert(ctx context.Context, db *sqlx.DB, table string, cols []string, rows [][]any) {
	colList := strings.Join(cols, ", ")
	colCount := len(cols)

	for start := 0; start < len(rows); start += batchSize {
		end := start + batchSize
		if end > len(rows) {
			end = len(rows)
		}
		batch := rows[start:end]

		var placeholders []string
		var args []any
		for i, row := range batch {
			var ph []string
			for j := range colCount {
				ph = append(ph, fmt.Sprintf("$%d", i*colCount+j+1))
			}
			placeholders = append(placeholders, "("+strings.Join(ph, ",")+")")
			args = append(args, row...)
		}

		q := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", table, colList, strings.Join(placeholders, ","))
		if _, err := db.ExecContext(ctx, q, args...); err != nil {
			panic(fmt.Sprintf("seed bulk insert %s: %v", table, err))
		}
	}
}

func ptr[T any](v T) *T { return &v }

func now() time.Time { return time.Now().UTC() }

// seedReferenceData creates countries, states, places, offices, and categories.
func seedReferenceData(ctx context.Context, db *sqlx.DB) ([]*reference.Place, []*reference.Office, []*reference.Category) {
	countryRepo := repository.NewCountryRepository(db)
	stateRepo := repository.NewStateRepository(db)
	placeRepo := repository.NewPlaceRepository(db)
	officeRepo := repository.NewOfficeRepository(db)
	catRepo := repository.NewCategoryRepository(db)

	type stateSpec struct {
		name         string
		conflictZone *string
		places       []string
	}
	type countrySpec struct {
		name   string
		code   string
		states []stateSpec
	}

	specs := []countrySpec{
		{
			name: "Ukraine", code: "UA",
			states: []stateSpec{
				{name: "Kyiv Oblast", places: []string{"Kyiv", "Brovary", "Irpin"}},
				{name: "Kharkiv Oblast", places: []string{"Kharkiv", "Izium", "Chuhuiv"}},
				{name: "Donetsk Oblast", conflictZone: ptr("active conflict zone"), places: []string{"Kramatorsk", "Sloviansk", "Bakhmut"}},
				{name: "Zaporizhzhia Oblast", conflictZone: ptr("frontline zone"), places: []string{"Zaporizhzhia", "Melitopol", "Enerhodar"}},
			},
		},
		{
			name: "Kyrgyzstan", code: "KG",
			states: []stateSpec{
				{name: "Bishkek City", places: []string{"Bishkek"}},
				{name: "Osh Oblast", places: []string{"Osh", "Jalal-Abad", "Uzgen"}},
				{name: "Chuy Oblast", places: []string{"Tokmok", "Kara-Balta", "Kemin"}},
			},
		},
		{
			name: "Germany", code: "DE",
			states: []stateSpec{
				{name: "Berlin", places: []string{"Berlin", "Spandau"}},
				{name: "Nordrhein-Westfalen", places: []string{"Cologne", "Dusseldorf", "Bonn"}},
			},
		},
	}

	var allPlaces []*reference.Place
	n := now()

	for _, cs := range specs {
		c := &reference.Country{ID: ulid.NewString(), Name: cs.name, Code: cs.code, CreatedAt: n, UpdatedAt: n}
		must(countryRepo.Create(ctx, c))

		for _, ss := range cs.states {
			s := &reference.State{ID: ulid.NewString(), CountryID: c.ID, Name: ss.name, ConflictZone: ss.conflictZone, CreatedAt: n, UpdatedAt: n}
			must(stateRepo.Create(ctx, s))

			for _, pn := range ss.places {
				p := &reference.Place{ID: ulid.NewString(), StateID: s.ID, Name: pn, CreatedAt: n, UpdatedAt: n}
				must(placeRepo.Create(ctx, p))
				allPlaces = append(allPlaces, p)
			}
		}
	}

	type officeSpec struct {
		name    string
		placeIx int
	}
	officeSpecs := []officeSpec{
		{"Bishkek Main Office", 12},
		{"Osh Field Office", 13},
		{"Kyiv Hub", 0},
		{"Berlin Liaison", 19},
	}

	var allOffices []*reference.Office
	for _, os := range officeSpecs {
		var placeID *string
		if os.placeIx < len(allPlaces) {
			placeID = &allPlaces[os.placeIx].ID
		}
		o := &reference.Office{ID: ulid.NewString(), Name: os.name, PlaceID: placeID, CreatedAt: n, UpdatedAt: n}
		must(officeRepo.Create(ctx, o))
		allOffices = append(allOffices, o)
	}

	categoryNames := []string{
		"Internally Displaced Person", "Refugee", "Asylum Seeker", "Returnee",
		"Stateless Person", "Conflict-Affected Civilian", "Unaccompanied Minor",
		"Single-Headed Household", "Person with Disability", "Elderly at Risk",
		"GBV Survivor", "Chronic Illness",
	}
	var allCategories []*reference.Category
	for _, cn := range categoryNames {
		c := &reference.Category{ID: ulid.NewString(), Name: cn, CreatedAt: n, UpdatedAt: n}
		must(catRepo.Create(ctx, c))
		allCategories = append(allCategories, c)
	}

	return allPlaces, allOffices, allCategories
}

func seedUsers(ctx context.Context, db *sqlx.DB) []*user.User {
	userRepo := repository.NewUserRepository(db)
	credRepo := repository.NewCredentialsRepository(db)
	hasher := crypto.NewArgonHasher()

	type spec struct {
		email string
		role  user.Role
		first string
		last  string
		phone string
	}

	specs := []spec{
		{"admin@example.com", user.RoleAdmin, "Admin", "User", "+10000000001"},
		{"staff@example.com", user.RoleStaff, "Staff", "User", "+10000000002"},
		{"consultant@example.com", user.RoleConsultant, "Consultant", "User", "+10000000003"},
		{"guest@example.com", user.RoleGuest, "Guest", "User", "+10000000004"},
	}

	var users []*user.User
	n := now()

	for _, s := range specs {
		uid := ulid.New()
		u := &user.User{
			ID:         uid,
			FirstName:  s.first,
			LastName:   s.last,
			Email:      s.email,
			Phone:      s.phone,
			Role:       s.role,
			IsVerified: true,
			IsActive:   true,
			CreatedAt:  n,
			UpdatedAt:  n,
		}
		must(userRepo.Create(ctx, u))

		hash, salt, err := hasher.Hash("password")
		must(err)
		must(credRepo.Create(ctx, &user.Credentials{
			UserID:       uid,
			PasswordHash: hash,
			Salt:         salt,
			UpdatedAt:    n,
		}))

		users = append(users, u)
	}

	return users
}

func seedProject(ctx context.Context, db *sqlx.DB, faker *gofakeit.Faker, admin *user.User, ix int) *project.Project {
	projRepo := repository.NewProjectRepository(db)
	n := now()
	desc := fmt.Sprintf("Development project #%d for testing", ix+1)
	p := &project.Project{
		ID:          ulid.NewString(),
		Name:        fmt.Sprintf("Project %s", faker.City()),
		Description: &desc,
		OwnerID:     admin.ID.String(),
		Status:      project.ProjectStatusActive,
		CreatedAt:   n,
		UpdatedAt:   n,
	}
	must(projRepo.Create(ctx, p))
	return p
}

func seedPermissions(ctx context.Context, db *sqlx.DB, projectID string, users []*user.User) {
	permRepo := repository.NewProjectPermissionRepository(db)
	n := now()

	type permSpec struct {
		userIx       int
		role         project.ProjectRole
		viewContact  bool
		viewPersonal bool
		viewDocs     bool
	}

	specs := []permSpec{
		{0, project.ProjectRoleOwner, true, true, true},
		{1, project.ProjectRoleManager, true, true, true},
		{2, project.ProjectRoleConsultant, true, false, false},
		{3, project.ProjectRoleViewer, false, false, false},
	}

	for _, s := range specs {
		must(permRepo.Create(ctx, &project.ProjectPermission{
			ID:               ulid.NewString(),
			ProjectID:        projectID,
			UserID:           users[s.userIx].ID.String(),
			Role:             s.role,
			CanViewContact:   s.viewContact,
			CanViewPersonal:  s.viewPersonal,
			CanViewDocuments: s.viewDocs,
			CreatedAt:        n,
			UpdatedAt:        n,
		}))
	}
}

func seedTags(ctx context.Context, db *sqlx.DB, faker *gofakeit.Faker, projectID string) []*tag.Tag {
	tagRepo := repository.NewTagRepository(db)
	pool := []string{"urgent", "follow-up", "legal-aid", "housing", "medical", "employment", "education", "documentation", "family", "vulnerable"}
	faker.ShuffleStrings(pool)

	count := faker.IntRange(5, 8)
	var tags []*tag.Tag
	for i := range count {
		t := &tag.Tag{ID: ulid.NewString(), ProjectID: projectID, Name: pool[i], CreatedAt: now()}
		must(tagRepo.Create(ctx, t))
		tags = append(tags, t)
	}
	return tags
}

// genPeople generates person structs in memory without inserting.
func genPeople(
	faker *gofakeit.Faker,
	projectID string, users []*user.User,
	places []*reference.Place, offices []*reference.Office,
	count int,
) []*person.Person {
	consultant := users[2]

	sexes := []person.Sex{person.SexMale, person.SexFemale, person.SexOther, person.SexUnknown}
	ageGroups := []person.AgeGroup{
		person.AgeGroupInfant, person.AgeGroupToddler, person.AgeGroupPreSchool,
		person.AgeGroupMiddleChildhood, person.AgeGroupYoungTeen, person.AgeGroupTeenager,
		person.AgeGroupYoungAdult, person.AgeGroupEarlyAdult, person.AgeGroupMiddleAged,
		person.AgeGroupOldAdult,
	}
	statuses := []person.CaseStatus{
		person.CaseStatusActive, person.CaseStatusActive, person.CaseStatusActive,
		person.CaseStatusNew, person.CaseStatusNew,
		person.CaseStatusClosed,
		person.CaseStatusArchived,
	}

	n := now()
	people := make([]*person.Person, 0, count)
	for range count {
		regDate := faker.DateRange(n.AddDate(-2, 0, 0), n)

		var birthDateVal *time.Time
		var ageGroupVal *person.AgeGroup
		if faker.IntRange(1, 10) <= 7 {
			bd := faker.DateRange(time.Date(1940, 1, 1, 0, 0, 0, 0, time.UTC), n.AddDate(-1, 0, 0))
			birthDateVal = &bd
		} else {
			ag := ageGroups[faker.IntRange(0, len(ageGroups)-1)]
			ageGroupVal = &ag
		}

		p := &person.Person{
			ID:             ulid.NewString(),
			ProjectID:      projectID,
			FirstName:      faker.FirstName(),
			LastName:       ptr(faker.LastName()),
			Sex:            sexes[faker.IntRange(0, len(sexes)-1)],
			AgeGroup:       ageGroupVal,
			BirthDate:      birthDateVal,
			PrimaryPhone:   ptr(faker.Phone()),
			PhoneNumbers:   json.RawMessage(`[]`),
			Email:          ptr(faker.Email()),
			CaseStatus:     statuses[faker.IntRange(0, len(statuses)-1)],
			OriginPlaceID:  &places[faker.IntRange(0, len(places)-1)].ID,
			CurrentPlaceID: &places[faker.IntRange(0, len(places)-1)].ID,
			ConsentGiven:   faker.IntRange(1, 10) <= 8,
			RegisteredAt:   &regDate,
			CreatedAt:      n,
			UpdatedAt:      n,
		}

		if p.ConsentGiven {
			cd := faker.DateRange(regDate, n)
			p.ConsentDate = &cd
		}

		if faker.IntRange(1, 2) == 1 {
			p.ConsultantID = ptr(consultant.ID.String())
		}

		if faker.IntRange(1, 10) <= 7 {
			p.OfficeID = &offices[faker.IntRange(0, len(offices)-1)].ID
		}

		people = append(people, p)
	}

	return people
}

func bulkInsertPeople(ctx context.Context, db *sqlx.DB, people []*person.Person) {
	cols := []string{
		"id", "project_id", "consultant_id", "office_id", "current_place_id", "origin_place_id",
		"external_id", "first_name", "last_name", "patronymic", "email", "birth_date", "sex", "age_group",
		"primary_phone", "phone_numbers", "case_status", "consent_given", "consent_date", "registered_at",
		"created_at", "updated_at",
	}
	rows := make([][]any, 0, len(people))
	for _, p := range people {
		rows = append(rows, []any{
			p.ID, p.ProjectID, p.ConsultantID, p.OfficeID, p.CurrentPlaceID, p.OriginPlaceID,
			p.ExternalID, p.FirstName, p.LastName, p.Patronymic, p.Email, p.BirthDate, p.Sex, p.AgeGroup,
			p.PrimaryPhone, p.PhoneNumbers, p.CaseStatus, p.ConsentGiven, p.ConsentDate, p.RegisteredAt,
			p.CreatedAt, p.UpdatedAt,
		})
	}
	bulkInsert(ctx, db, "people", cols, rows)
}

func genAndInsertPersonCategories(ctx context.Context, db *sqlx.DB, faker *gofakeit.Faker, people []*person.Person, categories []*reference.Category) {
	cols := []string{"person_id", "category_id"}
	var rows [][]any
	for _, p := range people {
		count := faker.IntRange(1, 3)
		used := map[int]bool{}
		for range count {
			ix := faker.IntRange(0, len(categories)-1)
			if used[ix] {
				continue
			}
			used[ix] = true
			rows = append(rows, []any{p.ID, categories[ix].ID})
		}
	}
	bulkInsert(ctx, db, "person_categories", cols, rows)
}

func genAndInsertPersonTags(ctx context.Context, db *sqlx.DB, faker *gofakeit.Faker, people []*person.Person, tags []*tag.Tag) {
	if len(tags) == 0 {
		return
	}
	cols := []string{"person_id", "tag_id"}
	var rows [][]any
	for _, p := range people {
		count := faker.IntRange(0, 2)
		used := map[int]bool{}
		for range count {
			ix := faker.IntRange(0, len(tags)-1)
			if used[ix] {
				continue
			}
			used[ix] = true
			rows = append(rows, []any{p.ID, tags[ix].ID})
		}
	}
	bulkInsert(ctx, db, "person_tags", cols, rows)
}

func genAndInsertSupportRecords(
	ctx context.Context, db *sqlx.DB, faker *gofakeit.Faker,
	projectID string, people []*person.Person,
	users []*user.User, offices []*reference.Office,
) {
	consultant := users[2]

	types := []support.SupportType{
		support.SupportTypeHumanitarian, support.SupportTypeLegal,
		support.SupportTypeSocial, support.SupportTypePsychological,
		support.SupportTypeMedical, support.SupportTypeGeneral,
	}
	spheres := []support.SupportSphere{
		support.SphereHousingAssistance, support.SphereDocumentRecovery,
		support.SphereSocialBenefits, support.SpherePropertyRights,
		support.SphereEmploymentRights, support.SphereFamilyLaw,
		support.SphereHealthcareAccess, support.SphereEducationAccess,
		support.SphereFinancialAid, support.SpherePsychSupport, support.SphereOther,
	}
	referralStatuses := []support.ReferralStatus{
		support.ReferralPending, support.ReferralAccepted,
		support.ReferralCompleted, support.ReferralDeclined,
		support.ReferralNoResponse,
	}

	cols := []string{
		"id", "person_id", "project_id", "consultant_id", "recorded_by", "office_id",
		"referred_to_office", "type", "sphere", "referral_status", "provided_at", "notes",
		"created_at", "updated_at",
	}
	n := now()
	var rows [][]any

	for _, p := range people {
		count := faker.IntRange(5, 20)
		for range count {
			var sphere *support.SupportSphere
			if faker.IntRange(1, 10) <= 7 {
				sphere = &spheres[faker.IntRange(0, len(spheres)-1)]
			}

			var providedAt *time.Time
			if p.RegisteredAt != nil {
				pd := faker.DateRange(*p.RegisteredAt, n)
				providedAt = &pd
			}

			var notes *string
			if faker.IntRange(1, 2) == 1 {
				notes = ptr(faker.Sentence(faker.IntRange(5, 15)))
			}

			var consultantID *string
			if faker.IntRange(1, 10) <= 6 {
				consultantID = ptr(consultant.ID.String())
			}

			var officeID *string
			if faker.IntRange(1, 2) == 1 {
				officeID = &offices[faker.IntRange(0, len(offices)-1)].ID
			}

			var referralStatus *support.ReferralStatus
			var referredToOffice *string
			if faker.IntRange(1, 10) <= 3 {
				referralStatus = &referralStatuses[faker.IntRange(0, len(referralStatuses)-1)]
				referredToOffice = &offices[faker.IntRange(0, len(offices)-1)].ID
			}

			rows = append(rows, []any{
				ulid.NewString(), p.ID, projectID, consultantID, nil, officeID,
				referredToOffice, types[faker.IntRange(0, len(types)-1)], sphere, referralStatus, providedAt, notes,
				n, n,
			})
		}
	}

	bulkInsert(ctx, db, "support_records", cols, rows)
}

func genAndInsertMigrationRecords(ctx context.Context, db *sqlx.DB, faker *gofakeit.Faker, people []*person.Person, places []*reference.Place) {
	reasons := []migration.MovementReason{
		migration.ReasonConflict, migration.ReasonSecurity,
		migration.ReasonServiceAccess, migration.ReasonReturn,
		migration.ReasonRelocationProgram, migration.ReasonEconomic,
		migration.ReasonOther,
	}
	housing := []migration.HousingAtDestination{
		migration.HousingOwnProperty, migration.HousingRenting,
		migration.HousingWithRelatives, migration.HousingCollectiveSite,
		migration.HousingHotel, migration.HousingOther, migration.HousingUnknown,
	}

	cols := []string{
		"id", "person_id", "from_place_id", "destination_place_id",
		"migration_date", "movement_reason", "housing_at_destination", "notes", "created_at",
	}
	n := now()
	var rows [][]any

	for _, p := range people {
		count := faker.IntRange(0, 2)
		for range count {
			var notes *string
			if faker.IntRange(1, 10) <= 4 {
				notes = ptr(faker.Sentence(faker.IntRange(5, 12)))
			}
			md := faker.DateRange(n.AddDate(-2, 0, 0), n)

			rows = append(rows, []any{
				ulid.NewString(), p.ID, p.OriginPlaceID, &places[faker.IntRange(0, len(places)-1)].ID,
				&md, &reasons[faker.IntRange(0, len(reasons)-1)], &housing[faker.IntRange(0, len(housing)-1)], notes, n,
			})
		}
	}

	bulkInsert(ctx, db, "migration_records", cols, rows)
}

func genAndInsertNotes(ctx context.Context, db *sqlx.DB, faker *gofakeit.Faker, people []*person.Person, users []*user.User) {
	authors := []string{users[1].ID.String(), users[2].ID.String()}
	cols := []string{"id", "person_id", "author_id", "body", "created_at"}
	n := now()
	var rows [][]any

	for _, p := range people {
		count := faker.IntRange(0, 2)
		for range count {
			rows = append(rows, []any{
				ulid.NewString(), p.ID, &authors[faker.IntRange(0, len(authors)-1)],
				faker.Sentence(faker.IntRange(5, 20)), n,
			})
		}
	}

	bulkInsert(ctx, db, "person_notes", cols, rows)
}

func genAndInsertPets(ctx context.Context, db *sqlx.DB, faker *gofakeit.Faker, projectID string, people []*person.Person) {
	statuses := []pet.PetStatus{
		pet.PetStatusRegistered, pet.PetStatusAdopted,
		pet.PetStatusOwnerFound, pet.PetStatusNeedsShelter,
	}
	cols := []string{"id", "project_id", "owner_id", "name", "status", "created_at", "updated_at"}
	n := now()
	var rows [][]any

	for _, p := range people {
		if faker.IntRange(1, 10) > 3 {
			continue
		}
		rows = append(rows, []any{
			ulid.NewString(), projectID, &p.ID, faker.PetName(),
			statuses[faker.IntRange(0, len(statuses)-1)], n, n,
		})
	}

	bulkInsert(ctx, db, "pets", cols, rows)
}

func seedHouseholds(ctx context.Context, db *sqlx.DB, faker *gofakeit.Faker, projectID string, people []*person.Person) {
	hhRepo := repository.NewHouseholdRepository(db)
	relationships := []household.Relationship{
		household.RelationshipSpouse, household.RelationshipChild,
		household.RelationshipParent, household.RelationshipSibling,
	}

	cols := []string{"household_id", "person_id", "relationship"}
	var memberRows [][]any

	for i := 0; i < len(people)-3; i += 5 {
		head := people[i]
		hhID := ulid.NewString()

		must(hhRepo.Create(ctx, &household.Household{
			ID:           hhID,
			ProjectID:    projectID,
			HeadPersonID: &head.ID,
			CreatedAt:    now(),
			UpdatedAt:    now(),
		}))

		memberRows = append(memberRows, []any{hhID, head.ID, household.RelationshipHead})

		memberCount := faker.IntRange(1, 3)
		for j := 1; j <= memberCount && i+j < len(people); j++ {
			memberRows = append(memberRows, []any{
				hhID, people[i+j].ID,
				relationships[faker.IntRange(0, len(relationships)-1)],
			})
		}
	}

	bulkInsert(ctx, db, "household_members", cols, memberRows)
}

func genAndInsertStatusHistory(ctx context.Context, db *sqlx.DB, faker *gofakeit.Faker, people []*person.Person) {
	cols := []string{"person_id", "from_status", "to_status", "changed_at"}
	var rows [][]any

	for _, p := range people {
		regAt := p.CreatedAt
		if p.RegisteredAt != nil {
			regAt = *p.RegisteredAt
		}

		status := p.CaseStatus
		if status == person.CaseStatusNew {
			continue
		}

		roll := faker.IntRange(1, 100)
		cursor := regAt

		switch {
		case roll <= 8:
			// 8% — new → closed directly (quick resolution)
			offset := time.Duration(faker.IntRange(1, 7)) * 24 * time.Hour
			cursor = cursor.Add(offset)
			rows = append(rows, []any{p.ID, string(person.CaseStatusNew), string(person.CaseStatusClosed), cursor})
			if status == person.CaseStatusArchived {
				offset = time.Duration(faker.IntRange(14, 90)) * 24 * time.Hour
				cursor = cursor.Add(offset)
				rows = append(rows, []any{p.ID, string(person.CaseStatusClosed), string(person.CaseStatusArchived), cursor})
			}

		case roll <= 13:
			// 5% — new → archived directly (abandoned)
			offset := time.Duration(faker.IntRange(30, 120)) * 24 * time.Hour
			cursor = cursor.Add(offset)
			rows = append(rows, []any{p.ID, string(person.CaseStatusNew), string(person.CaseStatusArchived), cursor})

		case roll <= 20:
			// 7% — new → active → archived (skip closed)
			offset := time.Duration(faker.IntRange(1, 14)) * 24 * time.Hour
			cursor = cursor.Add(offset)
			rows = append(rows, []any{p.ID, string(person.CaseStatusNew), string(person.CaseStatusActive), cursor})
			offset = time.Duration(faker.IntRange(14, 60)) * 24 * time.Hour
			cursor = cursor.Add(offset)
			rows = append(rows, []any{p.ID, string(person.CaseStatusActive), string(person.CaseStatusArchived), cursor})

		case roll <= 28:
			// 8% — new → active → closed → active (reopened) → closed → archived
			offset := time.Duration(faker.IntRange(1, 14)) * 24 * time.Hour
			cursor = cursor.Add(offset)
			rows = append(rows, []any{p.ID, string(person.CaseStatusNew), string(person.CaseStatusActive), cursor})
			offset = time.Duration(faker.IntRange(7, 45)) * 24 * time.Hour
			cursor = cursor.Add(offset)
			rows = append(rows, []any{p.ID, string(person.CaseStatusActive), string(person.CaseStatusClosed), cursor})
			offset = time.Duration(faker.IntRange(7, 30)) * 24 * time.Hour
			cursor = cursor.Add(offset)
			rows = append(rows, []any{p.ID, string(person.CaseStatusClosed), string(person.CaseStatusActive), cursor})
			if status == person.CaseStatusClosed || status == person.CaseStatusArchived {
				offset = time.Duration(faker.IntRange(7, 60)) * 24 * time.Hour
				cursor = cursor.Add(offset)
				rows = append(rows, []any{p.ID, string(person.CaseStatusActive), string(person.CaseStatusClosed), cursor})
			}
			if status == person.CaseStatusArchived {
				offset = time.Duration(faker.IntRange(30, 180)) * 24 * time.Hour
				cursor = cursor.Add(offset)
				rows = append(rows, []any{p.ID, string(person.CaseStatusClosed), string(person.CaseStatusArchived), cursor})
			}

		default:
			// 72% — standard linear path: new → active → closed → archived
			offset := time.Duration(faker.IntRange(1, 14)) * 24 * time.Hour
			cursor = cursor.Add(offset)
			rows = append(rows, []any{p.ID, string(person.CaseStatusNew), string(person.CaseStatusActive), cursor})

			if status == person.CaseStatusActive {
				continue
			}

			offset = time.Duration(faker.IntRange(7, 90)) * 24 * time.Hour
			cursor = cursor.Add(offset)
			rows = append(rows, []any{p.ID, string(person.CaseStatusActive), string(person.CaseStatusClosed), cursor})

			if status == person.CaseStatusClosed {
				continue
			}

			offset = time.Duration(faker.IntRange(30, 180)) * 24 * time.Hour
			cursor = cursor.Add(offset)
			rows = append(rows, []any{p.ID, string(person.CaseStatusClosed), string(person.CaseStatusArchived), cursor})
		}
	}

	if len(rows) > 0 {
		bulkInsert(ctx, db, "person_status_history", cols, rows)
	}
}

func must(err error) {
	if err != nil {
		panic(fmt.Sprintf("seed: %v", err))
	}
}
