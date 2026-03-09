# ADR-009: Global Search

| Field      | Value                                           |
| ---------- | ----------------------------------------------- |
| Status     | Accepted                                        |
| Date       | 2026-03-09                                      |
| Supersedes | —                                               |
| Components | observer (search endpoint), observer-web (cmdk) |

## Decision

Implement global search using a single `GET /search?q=<query>` endpoint backed by a
two-stage authorization filter, returning grouped results (people, pets) per project.
The frontend uses [cmdk](https://cmdk.paco.me/) triggered by `⌘K` / `Ctrl+K`.

## Scope

Search is **global** — it runs across all projects the requesting user has at least
read access to. Results are grouped by project, then by entity type within each project.
This lets users find any record from anywhere in the app without first navigating into
a project.

Searchable entities (Phase 1): **people**, **pets**, **projects**.

## Backend Design

### Two-stage authorization

**Stage 1 — resolve authorized project IDs:**

Query the `project_permissions` table for the authenticated user and collect all
`project_id` values where the user has any role (owner, manager, consultant, viewer).
Platform `admin` and `staff` roles implicitly have access to all projects and skip
this step.

This query is small, fast, and reusable — it can be extracted into a shared helper
used by other endpoints that need the same list.

**Stage 2 — search within authorized projects:**

Run parameterized `ILIKE` queries against the searchable entity tables, scoped to
the collected project IDs:

```sql
-- people: match on first_name, last_name, external_id
SELECT id, first_name, last_name, project_id, 'person' AS kind
FROM people
WHERE project_id = ANY($1)
  AND (
    first_name ILIKE $2 OR
    last_name  ILIKE $2 OR
    external_id ILIKE $2
  )
LIMIT 5;

-- pets: match on name, registration_id
SELECT id, name, project_id, 'pet' AS kind
FROM pets
WHERE project_id = ANY($1)
  AND (name ILIKE $2 OR registration_id ILIKE $2)
LIMIT 5;

-- projects: match on name (from authorized set only)
SELECT id, name, 'project' AS kind
FROM projects
WHERE id = ANY($1)
  AND name ILIKE $2
LIMIT 5;
```

`$2` is `'%' || query || '%'`. A minimum query length of **2 characters** is enforced
server-side to avoid full-table scans.

The three queries run concurrently (Go goroutines) and results are merged in the use
case layer before returning.

### Response shape

```json
{
  "query": "roxanne",
  "results": [
    {
      "project_id": "01J...",
      "project_name": "Lviv Program 2025",
      "people": [{ "id": "01J...", "first_name": "Roxanne", "last_name": "Doe" }],
      "pets": [],
      "projects": []
    }
  ],
  "total": 1
}
```

Results with no matching entities for a project are omitted from the array.

### Endpoint

```
GET /search?q=<query>
Authorization: Bearer <token>   (authenticated users only)
```

Minimum query length: 2 characters — returns `400` otherwise.
No project scoping parameter — authorization determines the scope.

### Architecture placement

- Repository method: `SearchRepository.Search(ctx, projectIDs []ulid.ULID, query string) (*SearchResults, error)`
- Use case: `internal/usecase/search/search_usecase.go` — resolves project IDs, fans out queries, merges results
- Handler: `internal/handler/search_handler.go` — thin adapter, validates `q` param

## Frontend Design

### Trigger

`⌘K` (macOS) / `Ctrl+K` (all platforms) opens the command palette from anywhere in
the app. A search icon button in the sidebar provides a mouse-accessible alternative.

### Component

[cmdk](https://cmdk.paco.me/) renders the palette. Results are grouped by project,
with entity-type sub-groups (People, Pets, Projects) inside each project group.

```
┌─────────────────────────────────────┐
│ 🔍 Search...                        │
├─────────────────────────────────────┤
│ Lviv Program 2025                   │
│   People                            │
│     Roxanne Doe                     │
│     Roxy Smith                      │
│   Pets                              │
│     Roxie (dog)                     │
├─────────────────────────────────────┤
│ → View all results                  │
└─────────────────────────────────────┘
```

Selecting a result navigates to the entity's detail page and closes the palette.
"View all results" navigates to `/search?q=<query>` — a dedicated results page.

### Result cap and "View all results"

The API returns at most **5 results per entity type per project**. When any entity
type returns 5 results the frontend shows a "View all results" item at the bottom of
the palette that links to the full results page. The results page calls the same
endpoint but with a higher server-side limit (configurable, default 50 per type).

### Debounce

Queries are debounced **300 ms** after the last keystroke. Requests with fewer than
2 characters are not sent.

### Hook

```ts
// hooks/use-search.ts
export function useSearch(query: string) {
  return useQuery({
    queryKey: ["search", query],
    queryFn: () => api.get(`/search?q=${encodeURIComponent(query)}`).json(),
    enabled: query.length >= 2,
    staleTime: 30_000,
  });
}
```

## Alternatives Considered

### A. Per-project search only

Scope search to the currently active project, no cross-project results.

**Rejected because**: users managing multiple projects (e.g. a coordinator across
three programs) would have to navigate into each project to search. Global scope is
more useful and the two-stage auth approach keeps it safe.

### B. Single JOIN query (permissions + entities)

One large SQL query joining `project_permissions`, `people`, `pets`, and `projects`.

**Rejected because**: the query becomes difficult to read, test, and extend. Pagination
across a multi-table union join is non-trivial. The two-stage approach separates
concerns cleanly — permission resolution is reusable and the search queries are
straightforward parameterized lookups.

### C. Full-text search (pg_trgm / tsvector)

Use PostgreSQL `pg_trgm` trigram indexes or `tsvector` full-text search for fuzzy
matching.

**Not rejected outright** — the `pg_trgm` extension is already enabled (migration
000001). For Phase 1, `ILIKE` is sufficient given expected data volumes (hundreds to
low thousands of records per project). Migrating the search queries to trigram index
scans (`%` operator) is a drop-in improvement if `ILIKE` becomes a bottleneck.
