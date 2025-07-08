# AGENTS.md

## Roles

### Feature Implementor

- Implements features based on the specifications found in `agents/feature-developer.md`.
- Works in coordination with other agents to ensure complete and stable integration.
- Follows conventional commits for all git messages.

### Style Enforcer

- Ensures all code adheres to the standards defined in `agents/code-styler.md`.
- Guidelines include:
  - Use `snake_case` for Python code.
  - Prefer named imports over `*` imports.
  - Ensure imports are sorted consistently.
  - Use consistent whitespace, indentation, and formatting.
- Runs formatting tools as defined in the `Justfile` (e.g., `just format`, `just lint`).

### Test Writer

- Follows guidelines in `agents/tester.md`.
- Expands and maintains tests under the `tests/` directory.
- Ensures coverage for edge cases, critical paths, and integrations.
- Uses `pytest` for Python testing and follows test structure conventions.
- Coordinates with the Feature Implementor to cover all newly introduced logic.

### Thinker

- Follows guidelines in `agents/architect.md`.
- Contributes to early-stage design, system architecture, and specifications.
- Maintains clarity and consistency across `agents/` documents.
- Reviews proposed features and offers insights before implementation begins.

### Code Executor

- Uses `Justfile` tasks to execute checks and validations:
  - `just format` — applies code formatters.
  - `just lint` — runs static analysis tools.
  - `just typecheck` — validates type hints and annotations.
  - `just test` — executes the full test suite.

## Infrastructure & Database Agents

### DBA Expert

- Follows guidelines in `agents/dba.md`.
- Designs and evolves the relational database schema using SQLAlchemy Core and Alembic.
- Maintains referential integrity and enforces domain constraints.
- Reviews and optimizes query performance and indexing strategy.
- Owns `observer/db`, `observer/entities` `migrations/` and ensures compatibility with production data.
- Coordinates with the Feature Implementor to ensure database support for new features.
- Avoids implicit database-generated defaults unless explicitly necessary.

## Notes

- All agents must follow the repository’s structure and conventions.
- Communication between agents should happen via GitHub Issues and Pull Requests.
