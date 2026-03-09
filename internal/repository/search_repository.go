package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type searchRepo struct {
	db *sqlx.DB
}

// NewSearchRepository creates a SearchRepository backed by Postgres.
func NewSearchRepository(db *sqlx.DB) SearchRepository {
	return &searchRepo{db: db}
}

func (r *searchRepo) ListAllProjectIDs(ctx context.Context) ([]string, error) {
	const q = `SELECT id FROM projects ORDER BY created_at`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list all project ids: %w", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *searchRepo) ListProjectIDsByUser(ctx context.Context, userID string) ([]string, error) {
	const q = `SELECT project_id FROM project_permissions WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list project ids by user: %w", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *searchRepo) ListProjectsByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	if len(ids) == 0 {
		return map[string]string{}, nil
	}

	const q = `SELECT id, name FROM projects WHERE id = ANY($1)`
	rows, err := r.db.QueryContext(ctx, q, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("list projects by ids: %w", err)
	}
	defer rows.Close()

	out := make(map[string]string, len(ids))
	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		out[id] = name
	}
	return out, rows.Err()
}

func (r *searchRepo) Search(ctx context.Context, projectIDs []string, query string, limit int) (*SearchHits, error) {
	if len(projectIDs) == 0 {
		return &SearchHits{}, nil
	}

	like := "%" + query + "%"
	var (
		mu   sync.Mutex
		hits SearchHits
		errs []error
		wg   sync.WaitGroup
	)

	wg.Add(3)

	go func() {
		defer wg.Done()
		results, err := r.searchPeople(ctx, projectIDs, like, limit)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			errs = append(errs, err)
			return
		}
		hits.People = results
	}()

	go func() {
		defer wg.Done()
		results, err := r.searchPets(ctx, projectIDs, like, limit)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			errs = append(errs, err)
			return
		}
		hits.Pets = results
	}()

	go func() {
		defer wg.Done()
		results, err := r.searchProjects(ctx, projectIDs, like, limit)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			errs = append(errs, err)
			return
		}
		hits.Projects = results
	}()

	wg.Wait()

	if len(errs) > 0 {
		return nil, fmt.Errorf("search queries: %w", errs[0])
	}
	return &hits, nil
}

func (r *searchRepo) searchPeople(ctx context.Context, projectIDs []string, like string, limit int) ([]PersonHit, error) {
	const q = `
		SELECT id, first_name, last_name, project_id
		FROM people
		WHERE project_id = ANY($1)
		  AND (first_name ILIKE $2 OR last_name ILIKE $2 OR external_id ILIKE $2)
		LIMIT $3`
	rows, err := r.db.QueryContext(ctx, q, pq.Array(projectIDs), like, limit)
	if err != nil {
		return nil, fmt.Errorf("search people: %w", err)
	}
	defer rows.Close()

	var results []PersonHit
	for rows.Next() {
		var h PersonHit
		if err := rows.Scan(&h.ID, &h.FirstName, &h.LastName, &h.ProjectID); err != nil {
			return nil, err
		}
		results = append(results, h)
	}
	return results, rows.Err()
}

func (r *searchRepo) searchPets(ctx context.Context, projectIDs []string, like string, limit int) ([]PetHit, error) {
	const q = `
		SELECT id, name, project_id
		FROM pets
		WHERE project_id = ANY($1)
		  AND (name ILIKE $2 OR registration_id ILIKE $2)
		LIMIT $3`
	rows, err := r.db.QueryContext(ctx, q, pq.Array(projectIDs), like, limit)
	if err != nil {
		return nil, fmt.Errorf("search pets: %w", err)
	}
	defer rows.Close()

	var results []PetHit
	for rows.Next() {
		var h PetHit
		if err := rows.Scan(&h.ID, &h.Name, &h.ProjectID); err != nil {
			return nil, err
		}
		results = append(results, h)
	}
	return results, rows.Err()
}

func (r *searchRepo) searchProjects(ctx context.Context, projectIDs []string, like string, limit int) ([]ProjectHit, error) {
	const q = `
		SELECT id, name
		FROM projects
		WHERE id = ANY($1)
		  AND name ILIKE $2
		LIMIT $3`
	rows, err := r.db.QueryContext(ctx, q, pq.Array(projectIDs), like, limit)
	if err != nil {
		return nil, fmt.Errorf("search projects: %w", err)
	}
	defer rows.Close()

	var results []ProjectHit
	for rows.Next() {
		var h ProjectHit
		if err := rows.Scan(&h.ID, &h.Name); err != nil {
			return nil, err
		}
		results = append(results, h)
	}
	return results, rows.Err()
}
