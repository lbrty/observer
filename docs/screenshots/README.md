# Screenshots

Playwright-based screenshot capture for documentation. Captures every page across all roles (admin, staff, consultant, guest) and copies selected screenshots to `docs/assets/images/screenshots/` for use in Hugo docs.

## Prerequisites

- Observer backend running with a seeded database (`just seed`)
- Frontend dev server running (`cd packages/observer-web && bun dev`)
- Chromium installed for Playwright

## Setup

```bash
cd docs/screenshots
bun install
bun run install-browsers
```

## Capture

Make sure the backend and frontend are running:

```bash
# Terminal 1: backend with seeded data
just serve

# Terminal 2: frontend dev server
cd packages/observer-web && bun dev

# Terminal 3: capture screenshots
cd docs/screenshots
bun run capture
```

This will:

1. Screenshot all public pages (login, register)
2. Log in as each role (admin, staff, consultant, guest) and screenshot every app page
3. Copy selected screenshots to `docs/assets/images/screenshots/` for Hugo

## Using an existing deployment

If running against a deployed instance instead of localhost, set the `baseURL` in `playwright.config.ts`:

```ts
use: {
  baseURL: "https://your-instance.example.com",
},
```

Make sure the test accounts exist in the database (see `ACCOUNTS` in `capture.ts`).

## Output

- `out/` — raw screenshots per role (gitignored)
- `test-results/` — Playwright artifacts (gitignored)
- `docs/assets/images/screenshots/` — curated copies used by Hugo docs
