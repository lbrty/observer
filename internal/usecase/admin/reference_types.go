package admin

import "time"

// CountryDTO is the admin-facing country representation.
type CountryDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateCountryInput holds data for creating a country.
type CreateCountryInput struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
}

// UpdateCountryInput holds data for updating a country.
type UpdateCountryInput struct {
	Name *string `json:"name"`
	Code *string `json:"code"`
}

// StateDTO is the admin-facing state representation.
type StateDTO struct {
	ID           string    `json:"id"`
	CountryID    string    `json:"country_id"`
	Name         string    `json:"name"`
	Code         *string   `json:"code,omitempty"`
	ConflictZone *string   `json:"conflict_zone,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateStateInput holds data for creating a state.
type CreateStateInput struct {
	Name         string  `json:"name" binding:"required"`
	Code         *string `json:"code"`
	ConflictZone *string `json:"conflict_zone"`
}

// UpdateStateInput holds data for updating a state.
type UpdateStateInput struct {
	Name         *string `json:"name"`
	Code         *string `json:"code"`
	ConflictZone *string `json:"conflict_zone"`
}

// PlaceDTO is the admin-facing place representation.
type PlaceDTO struct {
	ID        string    `json:"id"`
	StateID   string    `json:"state_id"`
	Name      string    `json:"name"`
	Lat       *float64  `json:"lat,omitempty"`
	Lon       *float64  `json:"lon,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreatePlaceInput holds data for creating a place.
type CreatePlaceInput struct {
	Name string   `json:"name" binding:"required"`
	Lat  *float64 `json:"lat"`
	Lon  *float64 `json:"lon"`
}

// UpdatePlaceInput holds data for updating a place.
type UpdatePlaceInput struct {
	Name *string  `json:"name"`
	Lat  *float64 `json:"lat"`
	Lon  *float64 `json:"lon"`
}

// OfficeDTO is the admin-facing office representation.
type OfficeDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	PlaceID   *string   `json:"place_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateOfficeInput holds data for creating an office.
type CreateOfficeInput struct {
	Name    string  `json:"name" binding:"required"`
	PlaceID *string `json:"place_id"`
}

// UpdateOfficeInput holds data for updating an office.
type UpdateOfficeInput struct {
	Name    *string `json:"name"`
	PlaceID *string `json:"place_id"`
}

// CategoryDTO is the admin-facing category representation.
type CategoryDTO struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateCategoryInput holds data for creating a category.
type CreateCategoryInput struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

// UpdateCategoryInput holds data for updating a category.
type UpdateCategoryInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
