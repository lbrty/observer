package project

// setPtr sets *dst to src when src is non-nil. Use for pointer-typed entity fields.
func setPtr[T any](dst **T, src *T) {
	if src != nil {
		*dst = src
	}
}

// applyOpt dereferences src into dst when src is non-nil. Use for value-typed entity fields.
func applyOpt[T any](dst *T, src *T) {
	if src != nil {
		*dst = *src
	}
}
