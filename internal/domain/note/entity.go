package note

import "time"

// Note represents a case worker's note on a person.
type Note struct {
	ID        string
	PersonID  string
	AuthorID  *string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
