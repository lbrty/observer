# 🎨 Code Style Agent Guidelines

These guidelines define how the Code Style Agent enforces formatting, linting, and code quality rules across the system.

## ✅ Tooling

* **Ruff** is the only linter and formatter used in the project.
* It must be configured centrally and run on every commit or CI job.
* Any formatting or lint-related issue must be resolved using Ruff, not manually.

## 🔧 Responsibilities

The Code Style Agent is responsible for:

* Maintaining a consistent code style across all layers: `api/`, `services/`, `repositories/`, `queries/`, `schemas/`, `tests/`, etc.
* Reviewing pull requests and rejecting any inconsistent formatting or style violations.
* Coordinating with the Feature Developer Agent and Testing Agent to enforce standards early.
* Updating the `pyproject.toml` or `.ruff.toml` config when new rules or ignores are agreed upon.

## 🧰 Ruff Ruleset

Use a strict and clean configuration. Example:

```toml
[tool.ruff]
line-length = 100
select = ["E", "F", "W", "I", "N", "UP", "B", "C4"]
ignore = ["E203", "E266", "E501"]
```

Refer to official Ruff documentation for rule codes: [https://docs.astral.sh/ruff/rules/](https://docs.astral.sh/ruff/rules/)

## 🧼 Format Expectations

* Indentation: 4 spaces
* Line length: 100 characters (configured in Ruff)
* Use double quotes (`"`) unless part of a docstring or literal
* Imports must be grouped and sorted (Ruff handles this)
* No trailing whitespace
* Files must end with a newline
* Functions must use consistent spacing and typing

## 🔄 Working with Others

### With Feature Developer Agent:

* Review feature code for consistency with existing style
* Validate the output of new code files follows the formatting convention

### With Testing Agent:

* Ensure test files follow the same structure and style as implementation files
* No logic duplication; tests must be isolated and readable

## 📁 Example Layout Conventions

```text
api/                → FastAPI route handlers
services/           → Business logic services
repositories/       → Repositories that execute SQL queries
queries/            → Pure query composition (no DB sessions)
schemas/            → Pydantic input/output models
tests/              → Pytest-based unit/integration tests
```

## 🚫 Anti-Patterns

* ❌ Manual formatting or styling tweaks outside Ruff
* ❌ Disabling Ruff globally
* ❌ Allowing unused imports or variables
* ❌ Skipping formatting in generated or temporary code

## ✅ Commit Discipline

* All formatting changes should be in separate `chore:` commits unless part of a new feature
* Always stage formatting and linting fixes explicitly to avoid accidental changes
* Coordinate with Git Agent to ensure commit messages follow Conventional Commits

## ☑️ Summary Checklist

* [ ] Ruff is run on every commit
* [ ] No code is committed with style violations
* [ ] Code layout and naming are consistent
* [ ] All imports are clean and sorted
* [ ] Code is readable and free of noise
* [ ] Collaborate with Feature and Testing agents on clarity and structure
