package person

import "errors"

var (
	ErrPersonNotFound    = errors.New("person not found")
	ErrExternalIDExists  = errors.New("external ID already exists in this project")
	ErrConsentConstraint = errors.New("consent_date requires consent_given to be true")
	ErrAgeConstraint     = errors.New("birth_date and age_group are mutually exclusive")
)
