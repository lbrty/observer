package audit

type LogInput struct {
	ProjectID  *string
	UserID     *string
	Action     string
	EntityType string
	EntityID   *string
	Summary    string
	Details    map[string]any
	IP         string
	UserAgent  string
}

type ListInput struct {
	ProjectID  *string `form:"project_id"`
	UserID     *string `form:"user_id"`
	Action     *string `form:"action"`
	EntityType *string `form:"entity_type"`
	DateFrom   *string `form:"date_from"`
	DateTo     *string `form:"date_to"`
	Page       int     `form:"page"`
	PerPage    int     `form:"per_page"`
}

type EntryDTO struct {
	ID            string         `json:"id"`
	ProjectID     *string        `json:"project_id"`
	UserID        *string        `json:"user_id"`
	Action        string         `json:"action"`
	EntityType    string         `json:"entity_type"`
	EntityID      *string        `json:"entity_id"`
	Summary       string         `json:"summary"`
	Details       map[string]any `json:"details,omitempty"`
	IP            *string        `json:"ip,omitempty"`
	UserAgent     *string        `json:"user_agent,omitempty"`
	CreatedAt     string         `json:"created_at"`
	UserFirstName *string        `json:"user_first_name"`
	UserLastName  *string        `json:"user_last_name"`
	UserEmail     *string        `json:"user_email"`
}

type ListOutput struct {
	Entries []EntryDTO `json:"entries"`
	Total   int        `json:"total"`
	Page    int        `json:"page"`
	PerPage int        `json:"per_page"`
}
