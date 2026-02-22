package reference

import "errors"

var (
	ErrCountryNotFound   = errors.New("country not found")
	ErrCountryCodeExists = errors.New("country code already exists")

	ErrStateNotFound = errors.New("state not found")

	ErrPlaceNotFound = errors.New("place not found")

	ErrOfficeNotFound = errors.New("office not found")

	ErrCategoryNotFound   = errors.New("category not found")
	ErrCategoryNameExists = errors.New("category name already exists")
)
