package search

// PersonHit is a search result for a person entity.
type PersonHit struct {
	ID        string
	FirstName string
	LastName  string
	ProjectID string
}

// PetHit is a search result for a pet entity.
type PetHit struct {
	ID        string
	Name      string
	ProjectID string
}

// ProjectHit is a search result for a project entity.
type ProjectHit struct {
	ID   string
	Name string
}

// SearchHits holds search results across all entity types.
type SearchHits struct {
	People   []PersonHit
	Pets     []PetHit
	Projects []ProjectHit
}
