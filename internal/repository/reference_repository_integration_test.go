//go:build !short

package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lbrty/observer/internal/domain/reference"
	"github.com/lbrty/observer/internal/repository"
	iulid "github.com/lbrty/observer/internal/ulid"
)

func TestCountryRepo_CRUD(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewCountryRepository(db)
	ctx := context.Background()

	c := &reference.Country{
		ID:   iulid.NewString(),
		Name: "Kyrgyzstan",
		Code: "KG",
	}

	// Create
	require.NoError(t, repo.Create(ctx, c))

	// Get
	got, err := repo.GetByID(ctx, c.ID)
	require.NoError(t, err)
	assert.Equal(t, c.ID, got.ID)
	assert.Equal(t, "Kyrgyzstan", got.Name)
	assert.Equal(t, "KG", got.Code)

	// Update
	c.Name = "Kyrgyz Republic"
	require.NoError(t, repo.Update(ctx, c))
	got, err = repo.GetByID(ctx, c.ID)
	require.NoError(t, err)
	assert.Equal(t, "Kyrgyz Republic", got.Name)

	// List
	list, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, list, 1)

	// Delete
	require.NoError(t, repo.Delete(ctx, c.ID))

	// Get after delete
	_, err = repo.GetByID(ctx, c.ID)
	assert.ErrorIs(t, err, reference.ErrCountryNotFound)
}

func TestCountryRepo_DuplicateCode(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewCountryRepository(db)
	ctx := context.Background()

	c1 := &reference.Country{ID: iulid.NewString(), Name: "Country A", Code: "XX"}
	require.NoError(t, repo.Create(ctx, c1))

	c2 := &reference.Country{ID: iulid.NewString(), Name: "Country B", Code: "XX"}
	err := repo.Create(ctx, c2)
	assert.ErrorIs(t, err, reference.ErrCountryCodeExists)
}

func TestStateRepo_CRUD(t *testing.T) {
	db := setupTestDB(t)
	countryRepo := repository.NewCountryRepository(db)
	stateRepo := repository.NewStateRepository(db)
	ctx := context.Background()

	country := &reference.Country{ID: iulid.NewString(), Name: "Ukraine", Code: "UA"}
	require.NoError(t, countryRepo.Create(ctx, country))

	zone := "eastern conflict zone"
	s := &reference.State{
		ID:           iulid.NewString(),
		CountryID:    country.ID,
		Name:         "Donetsk Oblast",
		ConflictZone: &zone,
	}

	// Create
	require.NoError(t, stateRepo.Create(ctx, s))

	// Get
	got, err := stateRepo.GetByID(ctx, s.ID)
	require.NoError(t, err)
	assert.Equal(t, s.ID, got.ID)
	assert.Equal(t, "Donetsk Oblast", got.Name)
	assert.Equal(t, country.ID, got.CountryID)
	require.NotNil(t, got.ConflictZone)
	assert.Equal(t, zone, *got.ConflictZone)

	// Update
	s.Name = "Donetska Oblast"
	require.NoError(t, stateRepo.Update(ctx, s))
	got, err = stateRepo.GetByID(ctx, s.ID)
	require.NoError(t, err)
	assert.Equal(t, "Donetska Oblast", got.Name)

	// List by country
	list, err := stateRepo.List(ctx, country.ID)
	require.NoError(t, err)
	assert.Len(t, list, 1)

	// Delete
	require.NoError(t, stateRepo.Delete(ctx, s.ID))
	_, err = stateRepo.GetByID(ctx, s.ID)
	assert.ErrorIs(t, err, reference.ErrStateNotFound)
}

func TestPlaceRepo_CRUD(t *testing.T) {
	db := setupTestDB(t)
	countryRepo := repository.NewCountryRepository(db)
	stateRepo := repository.NewStateRepository(db)
	placeRepo := repository.NewPlaceRepository(db)
	ctx := context.Background()

	country := &reference.Country{ID: iulid.NewString(), Name: "Ukraine", Code: "UA"}
	require.NoError(t, countryRepo.Create(ctx, country))

	state := &reference.State{ID: iulid.NewString(), CountryID: country.ID, Name: "Kyiv Oblast"}
	require.NoError(t, stateRepo.Create(ctx, state))

	lat := 50.4501
	lon := 30.5234
	p := &reference.Place{
		ID:      iulid.NewString(),
		StateID: state.ID,
		Name:    "Kyiv",
		Lat:     &lat,
		Lon:     &lon,
	}

	// Create
	require.NoError(t, placeRepo.Create(ctx, p))

	// Get
	got, err := placeRepo.GetByID(ctx, p.ID)
	require.NoError(t, err)
	assert.Equal(t, "Kyiv", got.Name)
	assert.Equal(t, state.ID, got.StateID)
	require.NotNil(t, got.Lat)
	assert.InDelta(t, lat, *got.Lat, 0.0001)

	// Update
	p.Name = "Kyiv City"
	require.NoError(t, placeRepo.Update(ctx, p))
	got, err = placeRepo.GetByID(ctx, p.ID)
	require.NoError(t, err)
	assert.Equal(t, "Kyiv City", got.Name)

	// List by state
	list, err := placeRepo.List(ctx, state.ID)
	require.NoError(t, err)
	assert.Len(t, list, 1)

	// Delete
	require.NoError(t, placeRepo.Delete(ctx, p.ID))
	_, err = placeRepo.GetByID(ctx, p.ID)
	assert.ErrorIs(t, err, reference.ErrPlaceNotFound)
}

func TestOfficeRepo_CRUD(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewOfficeRepository(db)
	ctx := context.Background()

	o := &reference.Office{
		ID:   iulid.NewString(),
		Name: "Main Office",
	}

	// Create
	require.NoError(t, repo.Create(ctx, o))

	// Get
	got, err := repo.GetByID(ctx, o.ID)
	require.NoError(t, err)
	assert.Equal(t, "Main Office", got.Name)
	assert.Nil(t, got.PlaceID)

	// Update
	o.Name = "Central Office"
	require.NoError(t, repo.Update(ctx, o))
	got, err = repo.GetByID(ctx, o.ID)
	require.NoError(t, err)
	assert.Equal(t, "Central Office", got.Name)

	// List
	list, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, list, 1)

	// Delete
	require.NoError(t, repo.Delete(ctx, o.ID))
	_, err = repo.GetByID(ctx, o.ID)
	assert.ErrorIs(t, err, reference.ErrOfficeNotFound)
}

func TestCategoryRepo_CRUD(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewCategoryRepository(db)
	ctx := context.Background()

	desc := "Internally displaced persons"
	c := &reference.Category{
		ID:          iulid.NewString(),
		Name:        "IDP",
		Description: &desc,
	}

	// Create
	require.NoError(t, repo.Create(ctx, c))

	// Get
	got, err := repo.GetByID(ctx, c.ID)
	require.NoError(t, err)
	assert.Equal(t, "IDP", got.Name)
	require.NotNil(t, got.Description)
	assert.Equal(t, desc, *got.Description)

	// Update
	c.Name = "IDP Category"
	require.NoError(t, repo.Update(ctx, c))
	got, err = repo.GetByID(ctx, c.ID)
	require.NoError(t, err)
	assert.Equal(t, "IDP Category", got.Name)

	// List
	list, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, list, 1)

	// Delete
	require.NoError(t, repo.Delete(ctx, c.ID))
	_, err = repo.GetByID(ctx, c.ID)
	assert.ErrorIs(t, err, reference.ErrCategoryNotFound)
}

func TestCategoryRepo_DuplicateName(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewCategoryRepository(db)
	ctx := context.Background()

	c1 := &reference.Category{ID: iulid.NewString(), Name: "Unique Name"}
	require.NoError(t, repo.Create(ctx, c1))

	c2 := &reference.Category{ID: iulid.NewString(), Name: "Unique Name"}
	err := repo.Create(ctx, c2)
	assert.ErrorIs(t, err, reference.ErrCategoryNameExists)
}
