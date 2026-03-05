package usecase

const (
	DefaultPage   = 1
	DefaultPerPage = 20
	MaxPerPage     = 100
)

// ClampPagination normalises page and perPage to safe defaults.
func ClampPagination(page, perPage int) (int, int) {
	if page < 1 {
		page = DefaultPage
	}
	if perPage < 1 {
		perPage = DefaultPerPage
	}
	if perPage > MaxPerPage {
		perPage = MaxPerPage
	}
	return page, perPage
}
