import { describe, expect, it, mock } from "bun:test";
import { renderHook, waitFor } from "@testing-library/react";

import { TestWrapper } from "@/test/wrapper";

mock.module("@/lib/api", () => ({
  api: {
    get: () => ({
      json: () => Promise.resolve({ projects: [{ id: "proj1", name: "Test" }], total: 1 }),
    }),
    post: () => ({
      json: () => Promise.resolve({ id: "proj-new", name: "New Project" }),
    }),
    patch: () => ({
      json: () => Promise.resolve({ id: "proj1", name: "Updated" }),
    }),
  },
  HTTPError: class extends Error {},
}));

const { useProjects, useProject, useCreateProject, useUpdateProject } = await import(
  "@/hooks/use-projects"
);

describe("useProjects", () => {
  it("fetches projects", async () => {
    const { result } = renderHook(() => useProjects(), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data?.projects).toHaveLength(1);
    expect(result.current.data?.projects[0].name).toBe("Test");
  });
});

describe("useProject", () => {
  it("fetches a single project", async () => {
    const { result } = renderHook(() => useProject("proj1"), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });

  it("does not fetch when id is empty", () => {
    const { result } = renderHook(() => useProject(""), {
      wrapper: TestWrapper,
    });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useCreateProject", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useCreateProject(), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useUpdateProject", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useUpdateProject(), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});
