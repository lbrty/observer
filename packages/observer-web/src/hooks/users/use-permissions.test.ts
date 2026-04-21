import { describe, expect, it, mock } from "bun:test";
import { renderHook, waitFor } from "@testing-library/react";

import { TestWrapper } from "@/test/wrapper";

mock.module("@/lib/api", () => ({
  api: {
    get: () => ({
      json: () =>
        Promise.resolve({
          permissions: [{ id: "perm1", user_id: "u1", role: "manager" }],
        }),
    }),
    post: () => ({
      json: () => Promise.resolve({ id: "perm-new", user_id: "u2", role: "viewer" }),
    }),
    patch: () => ({
      json: () => Promise.resolve({ id: "perm1", user_id: "u1", role: "editor" }),
    }),
    delete: () => ({ json: () => Promise.resolve() }),
  },
  HTTPError: class extends Error {},
}));

const { usePermissions, useAssignPermission, useUpdatePermission, useRevokePermission } =
  await import("@/hooks/use-permissions");

describe("usePermissions", () => {
  it("fetches permissions for a project", async () => {
    const { result } = renderHook(() => usePermissions("proj1"), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data?.permissions).toHaveLength(1);
    expect(result.current.data?.permissions[0].role).toBe("manager");
  });

  it("does not fetch when projectId is empty", () => {
    const { result } = renderHook(() => usePermissions(""), {
      wrapper: TestWrapper,
    });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useAssignPermission", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useAssignPermission(), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useUpdatePermission", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useUpdatePermission(), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useRevokePermission", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useRevokePermission(), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});
