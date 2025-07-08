# 💼 Feature Developer Agent Guidelines

These guidelines define how a Feature Developer Agent should implement application features in a consistent, secure, and maintainable way.

## 🐱 Structure Overview

* **queries/** – Reusable SQLAlchemy Core queries (no DB sessions or business logic).

  * Use for defining `select`, `insert`, `update`, `delete` expressions.
  * Example: `queries/people_queries.py::select_person_by_id()`

* **repositories/** – Thin layer that executes queries using a DB session.

  * Implements data access methods by composing queries and handling results.
  * Example: `people_repository.get_person_by_id(id)`

* **services/** – Business logic layer, unit-tested and dependency-injected.

  * Operates only through repositories (never directly uses queries).
  * Example: `people_service.create_person(data)`

* **api/** – FastAPI endpoints using dependency-injected services.

  * No direct access to repositories or queries.

## 🔧 Principles

* ❌ Do NOT use queries directly in services.
* ✅ Services talk to repositories, repositories talk to queries.
* ✅ Follow the dependency flow: `queries → repositories → services → api`.

## 🔄 Example Flow

```text
api/people.py               -> people_service.create_person()
services/people_service.py  -> people_repository.insert_person()
repositories/people_repo.py -> queries/people_queries.insert_person_query()
```

## 👥 Architecture Overview

The system follows a layered architecture:

* **FastAPI API Layer** — Exposes endpoints and handles request/response.
* **Dependency Injection Layer** — Injects services into API routes.
* **Service Layer** — Encapsulates business logic.
* **Repository Layer** — Encapsulates raw SQLAlchemy Core queries.
* **Database Schema** — Defined in coordination with the DBA Agent.
* **Context Singleton** — Initializes and provides access to services and repositories.

## 📌 Development Responsibilities

### ✅ 1. Schema Collaboration

* All new schemas must be proposed and refined with the **DBA Expert Agent**.

### ✅ 2. Primary Key Convention

* Use `ULID` for all primary keys.
* Do not use UUIDs or database-generated keys.

### ✅ 3. FastAPI Endpoint Definition

* Use `@router.get`, `@router.post`, etc. inside a module `api/routes/<resource>.py`.
* Do not use direct database access or logic in the endpoint.
* Use dependencies to inject services (see below).

## 💪 Testing Protocol

### Use `pytest`

* Create tests in `tests/<resource>/test_<action>.py`
* Include unit tests for the service layer and integration tests for endpoints.
* Testing should be coordinated with the **Testing Agent**.

## 🔧 Dependency Injection

### Global Context

* A central `context` object initializes:

  * database connection
  * services
  * repositories/queries

### Injecting Services

* Each service is injected using a dependency defined in:

  * `api/dependencies/services.py`

Example:

```python
from fastapi import Depends
from core.context import context

def people_service():
    return context.services.people
```

Use in route:

```python
@router.get("/people")
def list_people(service = Depends(people_service)):
    return service.list()
```

## 🧠 Service Layer

* Contains all business logic.
* Calls repository/query layer to fetch data.
* Should not know anything about FastAPI or request objects.
* Must record all data mutations via the audit log service.

## 📂 Repository Layer

* Only uses SQLAlchemy Core.
* No business logic or FastAPI.
* Follows the repository pattern or function-based `queries/people.py` style.
* Organized under `queries/<resource>.py` or `repositories/<resource>.py`

## 🕵️ Audit Logging

* Every **data mutation** (create, update, delete) must:

  * create an audit log record
  * capture: `resource`, `record_id`, `initiator`, `delta`, `timestamp`

Use `context.services.audit_log.record_change(...)` in service layer.

## 🎨 Code Style & Linting

* All formatting and linting is handled by **Ruff**.
* No other formatters or linters should be used.
* For any style inconsistencies, consult the **Code Styling Agent**.

## 🧠 Type Safety

* Use `astral/ty` for type-checking.
* Ensure all service layer, dependencies, and API endpoints pass strict type validation.

## ✅ Git Commit Standards

* Work with the **Git Agent** to ensure **Conventional Commits** are followed.
* Commit messages must reflect the scope and type of change (`feat:`, `fix:`, `chore:`, etc.)

## ✅ Summary Checklist

* [ ] Schema defined with DBA Agent
* [ ] ULID used for primary keys
* [ ] FastAPI endpoint with injected service
* [ ] Service layer implements logic
* [ ] Repositories only contain SQL
* [ ] Audit record on mutation
* [ ] Lint and format with Ruff
* [ ] Type check with astral/ty
* [ ] Test with pytest (in sync with Testing Agent)
* [ ] Conventional commits (Git Agent)
