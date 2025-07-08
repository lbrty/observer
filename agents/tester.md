# 🧪 Testing Agent Guidelines

These guidelines define how the Testing Agent ensures reliable, maintainable, and thorough test coverage across the system.

## ✅ Tooling

* Use **pytest** for all testing.
* No usage of `unittest`, `nose`, or other test runners.
* Type checks must be handled separately by **astral/ty**.
* Linting and formatting of tests is managed via **Ruff**, enforced by the Code Style Agent.

## 🧰 Test Types

### Unit Tests

* Focus on individual functions, especially in the **service** and **repository** layers.
* Must mock external dependencies (DB, APIs, etc.).
* Fast and deterministic.

### Integration Tests

* Target the **API layer** using `TestClient` from FastAPI.
* Use isolated test databases.
* Setup and teardown must be consistent across test runs.

### Regression Tests

* Add regression tests for every fixed bug.
* Always reproduce the bug first, then assert correct behavior after the fix.

## 📁 Project Layout

```text
tests/
├── conftest.py         # Shared fixtures and setup
├── api/                # Integration tests for FastAPI endpoints
├── services/           # Unit tests for business logic
├── repositories/       # Tests for data access logic
├── schemas/            # Tests for validation logic
```

## 🔁 Test Naming & Structure

* File: `test_<functionality>.py`
* Function: `test_<behavior>_<condition>()`
* One assertion per test is ideal
* Group tests with related fixtures into classes when needed

Example:

```python
def test_create_person_valid_payload():
    ...

def test_create_person_missing_name():
    ...
```

## ⚙️ Fixtures

* Shared fixtures go in `conftest.py`
* Scope should be as narrow as possible (`function` > `module` > `session`)
* Prefer factory functions over hardcoded values

## 🛑 Anti-Patterns

* ❌ Asserting internal state of mocks
* ❌ Tests depending on shared global state
* ❌ Testing third-party library behavior
* ❌ Skipping test cleanup or DB rollback

## 🤝 Collaboration with Other Agents

### Feature Developer Agent

* Coordinate to define coverage expectations for every feature
* Review acceptance criteria and translate to test scenarios

### Code Style Agent

* Ensure test files conform to project style and organization
* Clean imports, naming, and structure

## ✅ Commit Standards

* Tests should be committed in the same PR as the feature unless agreed otherwise
* Fixes should include regression test cases
* Use `test:` or `fix:` prefixes for test-related commits, coordinated with Git Agent

## ☑️ Summary Checklist

* [ ] Every service and route has test coverage
* [ ] Tests are structured into appropriate layers
* [ ] Fixtures are scoped and reusable
* [ ] Integration tests use isolated DB
* [ ] All tests are fast, reliable, and deterministic
* [ ] Follow naming and structure conventions
* [ ] Tests are typed and formatted
