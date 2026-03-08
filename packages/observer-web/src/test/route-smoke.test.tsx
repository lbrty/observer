import { describe, expect, it } from "bun:test";

// Route files call createFileRoute(path)({component}) at the top level.
// The component function itself is never invoked at import time, so hooks
// and stores do NOT need mocking — only the router's createFileRoute must
// not throw. Since @tanstack/react-router's real createFileRoute works
// without a running router (it returns a builder), we can import route
// files without any mocks as long as all transitive modules resolve.
//
// This test catches: missing modules, syntax errors, and top-level
// execution errors across all route files.

describe("route imports", () => {
  it("imports __root", async () => {
    const mod = await import("@/routes/__root");
    expect(mod).toBeDefined();
  });

  it("imports _auth layout", async () => {
    const mod = await import("@/routes/_auth");
    expect(mod).toBeDefined();
  });

  it("imports _app layout", async () => {
    const mod = await import("@/routes/_app");
    expect(mod).toBeDefined();
  });

  it("imports _app/index (dashboard)", async () => {
    const mod = await import("@/routes/_app/index");
    expect(mod).toBeDefined();
  });

  it("imports _app/profile", async () => {
    const mod = await import("@/routes/_app/profile");
    expect(mod).toBeDefined();
  });

  it("imports _app/admin layout", async () => {
    const mod = await import("@/routes/_app/admin");
    expect(mod).toBeDefined();
  });

  it("imports _app/admin/index", async () => {
    const mod = await import("@/routes/_app/admin/index");
    expect(mod).toBeDefined();
  });

  it("imports admin/users/index", async () => {
    const mod = await import("@/routes/_app/admin/users/index");
    expect(mod).toBeDefined();
  });

  it("imports admin/users/$userId", async () => {
    const mod = await import("@/routes/_app/admin/users/$userId");
    expect(mod).toBeDefined();
  });

  it("imports admin/projects/index", async () => {
    const mod = await import("@/routes/_app/admin/projects/index");
    expect(mod).toBeDefined();
  });

  it("imports admin/projects/$projectId layout", async () => {
    const mod = await import("@/routes/_app/admin/projects/$projectId");
    expect(mod).toBeDefined();
  });

  it("imports admin/projects/$projectId/index", async () => {
    const mod = await import("@/routes/_app/admin/projects/$projectId/index");
    expect(mod).toBeDefined();
  });

  it("imports admin/projects/$projectId/permissions", async () => {
    const mod = await import("@/routes/_app/admin/projects/$projectId/permissions");
    expect(mod).toBeDefined();
  });

  it("imports admin/reference/index", async () => {
    const mod = await import("@/routes/_app/admin/reference/index");
    expect(mod).toBeDefined();
  });

  it("imports admin/reference/categories", async () => {
    const mod = await import("@/routes/_app/admin/reference/categories");
    expect(mod).toBeDefined();
  });

  it("imports admin/reference/countries/index", async () => {
    const mod = await import("@/routes/_app/admin/reference/countries/index");
    expect(mod).toBeDefined();
  });

  it("imports admin/reference/countries/$countryId layout", async () => {
    const mod = await import("@/routes/_app/admin/reference/countries/$countryId");
    expect(mod).toBeDefined();
  });

  it("imports admin/reference/countries/$countryId/index", async () => {
    const mod = await import("@/routes/_app/admin/reference/countries/$countryId/index");
    expect(mod).toBeDefined();
  });

  it("imports admin/reference/countries/$countryId/states/$stateId", async () => {
    const mod = await import(
      "@/routes/_app/admin/reference/countries/$countryId/states/$stateId"
    );
    expect(mod).toBeDefined();
  });

  it("imports admin/reference/offices", async () => {
    const mod = await import("@/routes/_app/admin/reference/offices");
    expect(mod).toBeDefined();
  });

  it("imports projects/$projectId layout", async () => {
    const mod = await import("@/routes/_app/projects/$projectId");
    expect(mod).toBeDefined();
  });

  it("imports projects/$projectId/people/index", async () => {
    const mod = await import("@/routes/_app/projects/$projectId/people/index");
    expect(mod).toBeDefined();
  });

  it("imports projects/$projectId/people/$personId layout", async () => {
    const mod = await import("@/routes/_app/projects/$projectId/people/$personId");
    expect(mod).toBeDefined();
  });

  it("imports projects/$projectId/people/$personId/index", async () => {
    const mod = await import("@/routes/_app/projects/$projectId/people/$personId/index");
    expect(mod).toBeDefined();
  });

  it("imports projects/$projectId/people/$personId/notes", async () => {
    const mod = await import("@/routes/_app/projects/$projectId/people/$personId/notes");
    expect(mod).toBeDefined();
  });

  it("imports projects/$projectId/people/$personId/documents", async () => {
    const mod = await import(
      "@/routes/_app/projects/$projectId/people/$personId/documents"
    );
    expect(mod).toBeDefined();
  });

  it("imports projects/$projectId/people/$personId/migration-records", async () => {
    const mod = await import(
      "@/routes/_app/projects/$projectId/people/$personId/migration-records"
    );
    expect(mod).toBeDefined();
  });

  it("imports projects/$projectId/people/$personId/support-records", async () => {
    const mod = await import(
      "@/routes/_app/projects/$projectId/people/$personId/support-records"
    );
    expect(mod).toBeDefined();
  });

  it("imports projects/$projectId/people/$personId/stats", async () => {
    const mod = await import(
      "@/routes/_app/projects/$projectId/people/$personId/stats"
    );
    expect(mod).toBeDefined();
  });

  it("imports projects/$projectId/tags/index", async () => {
    const mod = await import("@/routes/_app/projects/$projectId/tags/index");
    expect(mod).toBeDefined();
  });

  it("imports projects/$projectId/pets/index", async () => {
    const mod = await import("@/routes/_app/projects/$projectId/pets/index");
    expect(mod).toBeDefined();
  });

  it("imports projects/$projectId/pets/$status", async () => {
    const mod = await import("@/routes/_app/projects/$projectId/pets/$status");
    expect(mod).toBeDefined();
  });

  it("imports projects/$projectId/pets/-pets-page", async () => {
    const mod = await import("@/routes/_app/projects/$projectId/pets/-pets-page");
    expect(mod).toBeDefined();
  });

  it("imports projects/$projectId/households/index", async () => {
    const mod = await import("@/routes/_app/projects/$projectId/households/index");
    expect(mod).toBeDefined();
  });

  it("imports projects/$projectId/support-records/index", async () => {
    const mod = await import(
      "@/routes/_app/projects/$projectId/support-records/index"
    );
    expect(mod).toBeDefined();
  });

  it("imports projects/$projectId/support-records/$type", async () => {
    const mod = await import(
      "@/routes/_app/projects/$projectId/support-records/$type"
    );
    expect(mod).toBeDefined();
  });

  it("imports projects/$projectId/support-records/-support-records-page", async () => {
    const mod = await import(
      "@/routes/_app/projects/$projectId/support-records/-support-records-page"
    );
    expect(mod).toBeDefined();
  });

  it("imports projects/$projectId/documents/index", async () => {
    const mod = await import("@/routes/_app/projects/$projectId/documents/index");
    expect(mod).toBeDefined();
  });

  it("imports projects/$projectId/reports/index", async () => {
    const mod = await import("@/routes/_app/projects/$projectId/reports/index");
    expect(mod).toBeDefined();
  });

  it("imports projects/$projectId/reports/people", async () => {
    const mod = await import("@/routes/_app/projects/$projectId/reports/people");
    expect(mod).toBeDefined();
  });

  it("imports projects/$projectId/reports/pets", async () => {
    const mod = await import("@/routes/_app/projects/$projectId/reports/pets");
    expect(mod).toBeDefined();
  });

  it("imports projects/$projectId/my-stats/index", async () => {
    const mod = await import("@/routes/_app/projects/$projectId/my-stats/index");
    expect(mod).toBeDefined();
  });
});
