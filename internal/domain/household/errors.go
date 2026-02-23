package household

import "errors"

var (
	ErrHouseholdNotFound = errors.New("household not found")
	ErrMemberNotFound    = errors.New("household member not found")
	ErrMemberExists      = errors.New("person is already a member of this household")
)
