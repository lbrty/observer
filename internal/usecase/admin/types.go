package admin

import "time"

// ListUsersInput holds query parameters for listing users.
type ListUsersInput struct {
	Page     int    `form:"page"`
	PerPage  int    `form:"per_page"`
	Search   string `form:"search"`
	Role     string `form:"role"`
	IsActive *bool  `form:"is_active"`
}

// ListUsersOutput is the paginated user list response.
type ListUsersOutput struct {
	Users   []UserDTO `json:"users"`
	Total   int       `json:"total"`
	Page    int       `json:"page"`
	PerPage int       `json:"per_page"`
}

// UpdateUserInput holds fields for a partial user update.
type UpdateUserInput struct {
	FirstName  *string `json:"first_name"`
	LastName   *string `json:"last_name"`
	Email      *string `json:"email"`
	Phone      *string `json:"phone"`
	OfficeID   *string `json:"office_id"`
	Role       *string `json:"role"`
	IsActive   *bool   `json:"is_active"`
	IsVerified *bool   `json:"is_verified"`
}

// UserDTO is the admin-facing user representation.
type UserDTO struct {
	ID         string    `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	OfficeID   *string   `json:"office_id,omitempty"`
	Role       string    `json:"role"`
	IsVerified bool      `json:"is_verified"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
