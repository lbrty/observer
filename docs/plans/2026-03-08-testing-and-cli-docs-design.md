# Testing & CLI Documentation Design

Date: 2026-03-08

## Backend Tests

### Handler Tests (new)
- httptest + Gin test mode, mock use cases via gomock
- All 17 handler files: auth, admin, permission, project, country, state, place, office, category, tag, person, support record, migration record, household, note, document, pet
- Error paths first (401, 403, 400, 404, 500), then happy paths
- One `_test.go` per handler, shared test helper for Gin context setup

### Use Case Tests (fill gaps)
- Missing: tag, person category, person tag, support record, household, pet
- gomock for repos, testify assert/require
- Error paths: repo failures, not-found, permission denied, validation

### Repository Integration Tests (new)
- testcontainers-go with Postgres, guarded by `testing.Short()`
- All repository implementations
- Constraint violations, unique conflicts, cascade deletes, empty results

## Frontend Tests

### Hook Tests (fill gaps)
- All missing hooks: use-persons, use-tags, use-pets, use-households, use-support-records, use-projects, use-auth, use-users
- mock.module() for api, renderHook with QueryClient wrapper
- Error responses (401, 404, 500), network failures, empty data

### Component Tests
- Auth: login/register forms — validation, submission failures
- Tables/lists: empty, loading, error states, pagination
- Forms: required fields, invalid input, server-side errors
- Charts: no data, edge-case data
- @testing-library/react + user-event + Happy DOM

### Route Smoke Tests
- Each route renders without crashing
- Auth routes: render, redirect when authenticated
- App routes: loading/error states with mocked API
- Route guards, 404 handling

## CLI Documentation

### --help Improvements
- Enhance cobra Long descriptions with usage examples
- Add Example field per command
- Consistent flag descriptions with defaults and env var refs

### docs/guides/cli.md
- Overview / quick start
- Per command: description, flags table, env vars, examples
- Commands: serve, migrate (up/create/version), keygen, create-admin, seed
- Common workflows: first-time setup, adding migrations, seeding dev data
