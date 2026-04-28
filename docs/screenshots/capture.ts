import { test, type Page, type APIRequestContext } from "@playwright/test";
import { copyFile, mkdir } from "node:fs/promises";
import { join } from "node:path";

const OUT = join(import.meta.dirname!, "out");
const API = process.env.VITE_API_URL ?? "http://localhost:9000/api";

type PageEntry = [name: string, path: string, setup?: (page: Page) => Promise<void>];

const PROJECT = "01KK9YYEFQ3WXR4E6H1D1YH9MX";
const USER = "01KK9YYD3ZD1XJ70WV8SXHERZX";
const COUNTRY = "01KK9YYD0XY4XBH890XF5ZQE9Y";
const STATE = "01KK9YYD0ZFV5QGE4SG4FE0NZV";

// Resolved in beforeAll by fetching the first real person in the project.
let PERSON = "";

const ACCOUNTS = [
  { role: "admin", email: "admin@example.com", password: "password" },
  { role: "staff", email: "staff@example.com", password: "password" },
  { role: "consultant", email: "consultant@example.com", password: "password" },
  { role: "guest", email: "guest@example.com", password: "password" },
] as const;

// Project roles to assign per platform role so screenshots show full UI.
// guest is always read-only at the platform level regardless of project role.
const PROJECT_ROLE: Partial<Record<(typeof ACCOUNTS)[number]["role"], string>> = {
  admin: "owner",
  staff: "manager",
  consultant: "consultant",
};

async function ensurePermissions(request: APIRequestContext) {
  const loginRes = await request.post(`${API}/auth/login`, {
    data: { email: "admin@example.com", password: "password" },
  });
  const { tokens } = await loginRes.json();
  const headers = { Authorization: `Bearer ${tokens.access_token}` };

  const [usersRes, permsRes] = await Promise.all([
    request.get(`${API}/admin/users`, { headers }),
    request.get(`${API}/admin/projects/${PROJECT}/permissions`, { headers }),
  ]);

  const { users } = await usersRes.json();
  const { permissions } = await permsRes.json();

  const existingByUserId = new Map<string, { id: string }>(
    permissions.map((p: { user_id: string; id: string }) => [p.user_id, p]),
  );

  for (const account of ACCOUNTS) {
    const projectRole = PROJECT_ROLE[account.role];
    if (!projectRole) continue;

    const user = users.find((u: { email: string }) => u.email === account.email);
    if (!user) continue;

    const body = {
      role: projectRole,
      can_view_contact: true,
      can_view_personal: true,
      can_view_documents: true,
      can_export: true,
    };

    const existing = existingByUserId.get(user.id);
    if (existing) {
      await request.patch(`${API}/admin/projects/${PROJECT}/permissions/${existing.id}`, {
        data: body,
        headers,
      });
    } else {
      await request.post(`${API}/admin/projects/${PROJECT}/permissions`, {
        data: { user_id: user.id, ...body },
        headers,
      });
    }
  }

  // Resolve a real person ID so person sub-page screenshots aren't blank.
  const peopleRes = await request.get(
    `${API}/projects/${PROJECT}/people?page=1&per_page=1`,
    { headers },
  );
  const { people } = await peopleRes.json();
  if (people?.length) {
    PERSON = people[0].id;
  }
}

test.beforeAll(async ({ request }) => {
  await ensurePermissions(request);
});

const PUBLIC_PAGES: PageEntry[] = [
  ["login", "/login"],
  ["register", "/register"],
];

const ADMIN_PAGES: PageEntry[] = [
  ["admin-users", "/admin/users"],
  ["admin-user-detail", `/admin/users/${USER}`],
  ["admin-projects", "/admin/projects"],
  ["admin-project-detail", `/admin/projects/${PROJECT}`],
  ["admin-project-permissions", `/admin/projects/${PROJECT}/permissions`],
  ["admin-reference", "/admin/reference"],
  ["admin-reference-countries", "/admin/reference/countries"],
  ["admin-reference-country-detail", `/admin/reference/countries/${COUNTRY}`],
  ["admin-reference-country-state", `/admin/reference/countries/${COUNTRY}/states/${STATE}`],
  ["admin-reference-offices", "/admin/reference/offices"],
  ["admin-reference-categories", "/admin/reference/categories"],
  ["admin-audit-logs", "/admin/audit-logs"],
];

// Wait for all loading skeletons to clear, catching slow API responses.
const waitForData = async (page: Page) => {
  await page.waitForLoadState("networkidle");
  await page
    .waitForFunction(() => document.querySelectorAll(".animate-pulse").length === 0, undefined, {
      timeout: 10_000,
    })
    .catch(() => {});
};

function appPages(): PageEntry[] {
  const p = PERSON;
  return [
    ["dashboard", "/"],
    ["profile", "/profile"],
    ["project-people", `/projects/${PROJECT}/people`],
    ["project-people-drawer", `/projects/${PROJECT}/people`, (pg) => openDrawer(pg, /register person/i)],
    ["project-person-detail", `/projects/${PROJECT}/people/${p}`, waitForData],
    ["project-person-documents", `/projects/${PROJECT}/people/${p}/documents`, waitForData],
    ["project-person-support-records", `/projects/${PROJECT}/people/${p}/support-records`, waitForData],
    [
      "project-person-support-records-drawer",
      `/projects/${PROJECT}/people/${p}/support-records`,
      async (pg) => { await waitForData(pg); await openDrawer(pg, /new record/i); },
    ],
    ["project-person-migration-records", `/projects/${PROJECT}/people/${p}/migration-records`, waitForData],
    [
      "project-person-migration-records-drawer",
      `/projects/${PROJECT}/people/${p}/migration-records`,
      async (pg) => { await waitForData(pg); await openDrawer(pg, /^add$/i); },
    ],
    ["project-person-notes", `/projects/${PROJECT}/people/${p}/notes`, waitForData],
    ["project-person-stats", `/projects/${PROJECT}/people/${p}/stats`, waitForData],
    ["project-support-records", `/projects/${PROJECT}/support-records`],
    ["project-support-records-drawer", `/projects/${PROJECT}/support-records`, (pg) => openDrawer(pg, /new record/i)],
    ["project-households", `/projects/${PROJECT}/households`],
    ["project-households-drawer", `/projects/${PROJECT}/households`, (pg) => openDrawer(pg, /new household/i)],
    ["project-tags", `/projects/${PROJECT}/tags`],
    ["project-tags-drawer", `/projects/${PROJECT}/tags`, (pg) => openDrawer(pg, /add tag/i)],
    ["project-pets", `/projects/${PROJECT}/pets`],
    ["project-pets-drawer", `/projects/${PROJECT}/pets`, (pg) => openDrawer(pg, /register pet/i)],
    ["project-reports-people", `/projects/${PROJECT}/reports/people`, waitForData],
    ["project-reports-pets", `/projects/${PROJECT}/reports/pets`, waitForData],
    ["project-reports-custom", `/projects/${PROJECT}/reports/custom`],
    ["project-audit-logs", `/projects/${PROJECT}/audit-logs`],
    ["project-my-stats", `/projects/${PROJECT}/my-stats`],
  ];
}

async function openDrawer(page: Page, buttonText: RegExp) {
  const btn = page.getByRole("button", { name: buttonText });
  const count = await btn.count();
  if (count === 0) return;
  await btn.click();
  await page.waitForTimeout(500);
}

async function setupPage(page: Page) {
  await page.addInitScript(() => {
    localStorage.setItem("observer-lang", "en");
  });
}

async function login(page: Page, email: string, password: string) {
  await page.goto("/login");
  await page.waitForLoadState("networkidle");
  await page.locator('input[name="email"]').waitFor({ state: "visible", timeout: 15_000 });
  await page.locator('input[name="email"]').fill(email);
  await page.locator('input[name="password"]').fill(password);
  await page.locator('button[type="submit"]').click();
  await page.waitForURL((url) => !url.pathname.includes("/login"), {
    timeout: 15_000,
  });
}

async function capture(page: Page, dir: string, pages: PageEntry[]) {
  await mkdir(dir, { recursive: true });

  for (const [name, path, setup] of pages) {
    await page.goto(path);
    await page.waitForLoadState("networkidle");
    await page.waitForTimeout(2000);

    if (setup) {
      await setup(page);
    }

    await page.screenshot({
      path: join(dir, `${name}.png`),
      fullPage: false,
    });

    console.log(`captured: ${dir}/${name}.png`);
  }
}

test("capture public pages", async ({ page }) => {
  await setupPage(page);
  await capture(page, join(OUT, "public"), PUBLIC_PAGES);
});

for (const account of ACCOUNTS) {
  test(`capture ${account.role} pages`, async ({ page }) => {
    await setupPage(page);
    await login(page, account.email, account.password);

    const pages = account.role === "admin" ? [...appPages(), ...ADMIN_PAGES] : appPages();

    await capture(page, join(OUT, account.role), pages);
  });
}

const DOCS_DIR = join(import.meta.dirname!, "..", "assets", "images", "screenshots");

const COPY_MAP: Record<string, string> = {
  "public/login.png": "login.png",
  "admin/dashboard.png": "dashboard.png",
  "admin/project-people.png": "people-list.png",
  "admin/project-people-drawer.png": "people-register.png",
  "admin/project-person-detail.png": "person-detail.png",
  "admin/project-person-support-records.png": "support-records.png",
  "admin/project-person-support-records-drawer.png": "support-record-form.png",
  "admin/project-person-migration-records.png": "migration-records.png",
  "admin/project-person-migration-records-drawer.png": "migration-record-form.png",
  "admin/project-person-notes.png": "notes.png",
  "admin/project-person-documents.png": "documents.png",
  "admin/project-person-stats.png": "person-stats.png",
  "admin/project-households.png": "households.png",
  "admin/project-households-drawer.png": "household-form.png",
  "admin/project-tags.png": "tags.png",
  "admin/project-pets.png": "pets.png",
  "admin/project-reports-people.png": "reports.png",
  "admin/project-reports-pets.png": "reports-pets.png",
  "admin/project-reports-custom.png": "reports-custom.png",
  "admin/project-audit-logs.png": "project-audit-logs.png",
  "admin/admin-audit-logs.png": "admin-audit-logs.png",
  "admin/admin-users.png": "admin-users.png",
  "admin/admin-projects.png": "admin-projects.png",
  "admin/admin-project-permissions.png": "admin-permissions.png",
  "admin/admin-reference.png": "admin-reference.png",
  "consultant/project-my-stats.png": "my-stats.png",
  "staff/project-people.png": "people-list-staff.png",
  "consultant/project-people.png": "people-list-consultant.png",
  "guest/project-people.png": "people-list-guest.png",
};

test("copy to docs", async () => {
  await mkdir(DOCS_DIR, { recursive: true });

  for (const [src, dest] of Object.entries(COPY_MAP)) {
    await copyFile(join(OUT, src), join(DOCS_DIR, dest));
    console.log(`copied: ${src} -> ${dest}`);
  }
});
