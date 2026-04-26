package search

import (
	"context"
	"fmt"
	"time"

	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/repository"
)

const searchTimeout = 30 * time.Second

// SearchUseCase executes global cross-project search.
type SearchUseCase struct {
	repo repository.SearchRepository
}

// NewSearchUseCase creates a SearchUseCase.
func NewSearchUseCase(repo repository.SearchRepository) *SearchUseCase {
	return &SearchUseCase{repo: repo}
}

// Execute runs a global search scoped to the projects the user is authorised to see.
func (uc *SearchUseCase) Execute(ctx context.Context, userID string, role user.Role, query string, limit int) (*SearchOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, searchTimeout)
	defer cancel()

	var (
		projectIDs []string
		err        error
	)

	if role == user.RoleAdmin || role == user.RoleStaff {
		projectIDs, err = uc.repo.ListAllProjectIDs(ctx)
	} else {
		projectIDs, err = uc.repo.ListProjectIDsByUser(ctx, userID)
	}

	if err != nil {
		return nil, fmt.Errorf("resolve project ids: %w", err)
	}

	if len(projectIDs) == 0 {
		return &SearchOutput{Query: query, Results: []ProjectGroup{}, Total: 0}, nil
	}

	hits, err := uc.repo.Search(ctx, projectIDs, query, limit)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	// Collect all unique project IDs referenced by the hits.
	seen := make(map[string]struct{})
	for _, h := range hits.Projects {
		seen[h.ID] = struct{}{}
	}

	for _, h := range hits.People {
		seen[h.ProjectID] = struct{}{}
	}

	for _, h := range hits.Pets {
		seen[h.ProjectID] = struct{}{}
	}

	if len(seen) == 0 {
		return &SearchOutput{Query: query, Results: []ProjectGroup{}, Total: 0}, nil
	}

	// Fetch project names for all referenced project IDs.
	allIDs := make([]string, 0, len(seen))
	for id := range seen {
		allIDs = append(allIDs, id)
	}

	nameMap, err := uc.repo.ListProjectsByIDs(ctx, allIDs)
	if err != nil {
		return nil, fmt.Errorf("fetch project names: %w", err)
	}

	// Build groups keyed by project ID.
	groups := make(map[string]*ProjectGroup, len(seen))
	getGroup := func(projectID string) *ProjectGroup {
		g, ok := groups[projectID]
		if !ok {
			g = &ProjectGroup{
				ProjectID:   projectID,
				ProjectName: nameMap[projectID],
				People:      []PersonResult{},
				Pets:        []PetResult{},
				Projects:    []ProjectResult{},
			}
			groups[projectID] = g
		}
		return g
	}

	for _, h := range hits.Projects {
		g := getGroup(h.ID)
		g.Projects = append(g.Projects, ProjectResult{ID: h.ID, Name: h.Name})
	}

	for _, h := range hits.People {
		g := getGroup(h.ProjectID)
		g.People = append(g.People, PersonResult{ID: h.ID, FirstName: h.FirstName, LastName: h.LastName})
	}

	for _, h := range hits.Pets {
		g := getGroup(h.ProjectID)
		g.Pets = append(g.Pets, PetResult{ID: h.ID, Name: h.Name})
	}

	results := make([]ProjectGroup, 0, len(groups))
	total := 0
	for _, g := range groups {
		total += len(g.People) + len(g.Pets) + len(g.Projects)
		results = append(results, *g)
	}

	return &SearchOutput{Query: query, Results: results, Total: total}, nil
}
