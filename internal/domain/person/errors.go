package person

import "errors"

var (
	ErrPersonNotFound   = errors.New("person not found")
	ErrExternalIDExists = errors.New("external ID already exists in this project")
	// Both constraints are enforced by DB CHECK constraints (chk_people_consent, chk_people_age_xor).
	// These errors translate raw DB violations into typed domain errors.
	// TODO: add a Person.Validate() method to catch violations before the DB round-trip.
	ErrConsentConstraint = errors.New("consent_date requires consent_given to be true")
	ErrAgeConstraint     = errors.New("birth_date and age_group are mutually exclusive")
)
