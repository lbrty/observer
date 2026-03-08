import { test, type Page } from "@playwright/test";
import { copyFile, mkdir } from "node:fs/promises";
import { join } from "node:path";

const OUT = join(import.meta.dirname!, "out");

type PageEntry = [
  name: string,
  path: string,
  setup?: (page: Page) => Promise<void>,
];

const PROJECT = "01KJZDSG53QYEJ5JHD7YRHERV7";
const USER = "01KJZDSG33GM2VRF33K8Y5JZ0M";
const PERSON = "01KJZDSG5GHNTQQ1TA1GF5S6FP";
const COUNTRY = "01KJZDSG16FR09T487RBEYZ9V7";
const STATE = "01KJZDSG26VNJKETRQ3TGKR8A1";

const ACCOUNTS = [
  { role: "admin", email: "admin@example.com", password: "password" },
  { role: "staff", email: "staff@example.com", password: "password" },
  { role: "consultant", email: "consultant@example.com", password: "password" },
  { role: "guest", email: "guest@example.com", password: "password" },
] as const;

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
];

const APP_PAGES: PageEntry[] = [
  ["dashboard", "/"],
  ["profile", "/profile"],
  ["project-people", `/projects/${PROJECT}/people`],
  ["project-people-drawer", `/projects/${PROJECT}/people`, (p) => openDrawer(p, /register person/i)],
  ["project-person-detail", `/projects/${PROJECT}/people/${PERSON}`],
  ["project-person-documents", `/projects/${PROJECT}/people/${PERSON}/documents`],
  ["project-person-support-records", `/projects/${PROJECT}/people/${PERSON}/support-records`],
  ["project-person-support-records-drawer", `/projects/${PROJECT}/people/${PERSON}/support-records`, (p) => openDrawer(p, /new record/i)],
  ["project-person-migration-records", `/projects/${PROJECT}/people/${PERSON}/migration-records`],
  ["project-person-migration-records-drawer", `/projects/${PROJECT}/people/${PERSON}/migration-records`, (p) => openDrawer(p, /^add$/i)],
  ["project-person-notes", `/projects/${PROJECT}/people/${PERSON}/notes`],
  ["project-person-stats", `/projects/${PROJECT}/people/${PERSON}/stats`],
  ["project-documents", `/projects/${PROJECT}/documents`],
  ["project-support-records", `/projects/${PROJECT}/support-records`],
  ["project-support-records-drawer", `/projects/${PROJECT}/support-records`, (p) => openDrawer(p, /new record/i)],
  ["project-households", `/projects/${PROJECT}/households`],
  ["project-households-drawer", `/projects/${PROJECT}/households`, (p) => openDrawer(p, /new household/i)],
  ["project-tags", `/projects/${PROJECT}/tags`],
  ["project-tags-drawer", `/projects/${PROJECT}/tags`, (p) => openDrawer(p, /add tag/i)],
  ["project-pets", `/projects/${PROJECT}/pets`],
  ["project-pets-drawer", `/projects/${PROJECT}/pets`, (p) => openDrawer(p, /register pet/i)],
  ["project-reports", `/projects/${PROJECT}/reports`],
  ["project-my-stats", `/projects/${PROJECT}/my-stats`],
];

async function openDrawer(page: Page, buttonText: RegExp) {
  const btn = page.getByRole("button", { name: buttonText });
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

    const pages =
      account.role === "admin" ? [...APP_PAGES, ...ADMIN_PAGES] : APP_PAGES;

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
  "admin/project-reports.png": "reports.png",
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
