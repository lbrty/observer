package search

// PersonResult is a DTO for a matching person.
type PersonResult struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// PetResult is a DTO for a matching pet.
type PetResult struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ProjectResult is a DTO for a matching project.
type ProjectResult struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ProjectGroup groups search results under one project.
type ProjectGroup struct {
	ProjectID   string          `json:"project_id"`
	ProjectName string          `json:"project_name"`
	People      []PersonResult  `json:"people"`
	Pets        []PetResult     `json:"pets"`
	Projects    []ProjectResult `json:"projects"`
}

// SearchOutput is the response DTO for the search endpoint.
type SearchOutput struct {
	Query   string         `json:"query"`
	Results []ProjectGroup `json:"results"`
	Total   int            `json:"total"`
}
