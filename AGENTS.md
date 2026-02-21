# Coding Style & Preferences Guide

## Language & Naming

- **Variable naming**: Prefer shorter, clear names
  - `idx` → `ix` for index variables
  - `uniq` → `uq` for unique identifiers
  - Keep names concise but meaningful
- **Default language**: Kyrgyz Latin (ky), not Russian

## Code Style


## Code Style

- **Comments**: Simple docstrings only.
    - No decorative separators
        - `//-----`
        - `//=====`
        - `/* ── Sidebar ────────── */`
    - no ASCII art
- **Complex logic**: Use mermaid diagrams over lengthy text explanations. Add a README in the module if needed.
- **Architecture**: Domain-Driven Design + Clean Architecture. Manual dependency injection, no frameworks. Business logic in application layer, not database.
- **Dependencies**: Prefer well-maintained, widely-known libraries. Leverage existing solutions over building custom.
- **Philosophy**: Pragmatic MVP — core functionality over complex features, essential security controls with practical implementation.

- **Documentation**:
  - Keep docstrings brief and functional
  - For complex logic: Create mermaid diagrams instead of lengthy explanations
  - Visual diagrams > walls of text
  - For complex module logic: Add README file documenting business logic with mermaid diagrams where needed

- **Structure**: Follow Go best practices
  - `cmd/` for entrypoints
  - `internal/` for app-specific packages

## Dependencies & Libraries

- **Prioritize**: Well-maintained, widely-known projects
- **Current stack**: Gin, testify, testcontainers-go, Redis client
- **Avoid**: Obscure or poorly-maintained libraries
- **Principle**: Leverage existing solutions over building custom

## Architecture Patterns

- **Design**: Domain-Driven Design + Clean Architecture
- **Dependency injection**: Manual, no frameworks
- **Business logic**: Application layer (Go), not database
- **Handlers**: Thin HTTP layers, logic in use cases

## MVP Philosophy

- **Approach**: Pragmatic simplicity
- **Priority**: Core functionality over complex features
- **Defer**: Advanced features (MEK/DEK encryption, detailed audit logs) to Phase 2
- **Security**: Balance essential controls with practical implementation

## Testing

- Integration tests with testcontainers-go
- Unit tests with testify
- Test real database interactions

## Build Tools

- Prefer Justfile over Makefile
- For frontend bun is the default package manager and bundler

## Project related variables

Use variables defined in: docs/variables.md whenever you implement new feature or bootstrap part of project.
Also actively suggest which variables can be added if there are variables which appear multiple times.

 skip verification. If tests fail, fix them before proceeding.

## Frontend Imports

- Use the `@/` alias for all imports (e.g. `@/stores/auth`, `@/components/auth`).
- Only exception: colocated files (siblings in the same directory) use relative `./` imports (e.g. `./constants`).
- Style module imports (`.module.css`) follow the same rules — alias for distant files, relative for colocated.
- Extract shared constants (config arrays, magic values, labels) `constants.ts` when they clutter component files.
- Extract shared constants which can be re-used by other parts of the application to the root level `constants.ts`.
- Import sorting order (each group separated by a blank line):
    1. **React** — `react`, `react-dom`
    2. **External libs** – `@tanstack/*`, `@zxcvbn-ts/*`, etc.
    3. **Workspace packages** — `@observer/*`
    4. **App aliases** — `@/components/*`, `@/stores/*`, `@/hooks/*`, etc.
    5. **Colocated** — `./constants`, `./types`, etc.
    6. **Styles** (always last, separated by blank line) — `.module.css` imports

## Component implementation

We have `base-ui` and `@phosphor-icons/react` first check if components and icons exist in both and try to take ready 
to use options.

## React compiler

Observer uses React compiler, consider omitting specifying dependencies to effects or skip effects altogether when possible.

## When writing tailwind styles

Use `@apply` helper from tailwind/css and always follow the following structure and the order and each group
must be on separate lines if the ruleset is >10 rules.

`@apply` should be separate per each group.

- positioninig `@apply absolute|relative` then, go top, left, bottom, right if needed.
- layout display options like `@apply flex|block`
- sizing/dimensions `@apply w-full` or `w-[X]` depending on situation.
- border configurations `@apply rounded-sm border-none`.
- background configuration `@apply bg-transparent` etc.
- paddings and margins `@apply p-1 m-1` etc.
- text styles `@apply text-lg;` or custom sizes literals or bound to variables.
- translation options `-translate-y-1/2` etc.
- here comes the rest but always lookup to the list and try to group by responsibilities.
