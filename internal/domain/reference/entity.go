package reference

import "time"

// Country represents a country record.
type Country struct {
	ID        string
	Name      string
	Code      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// State represents a state/oblast within a country.
type State struct {
	ID           string
	CountryID    string
	Name         string
	Code         *string
	ConflictZone *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Place represents a locality within a state.
type Place struct {
	ID        string
	StateID   string
	Name      string
	Lat       *float64
	Lon       *float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Office represents a physical office location.
type Office struct {
	ID        string
	Name      string
	PlaceID   *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Category represents a person classification category.
type Category struct {
	ID          string
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
