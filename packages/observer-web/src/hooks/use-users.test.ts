import { describe, expect, it, mock } from "bun:test";
import { renderHook, waitFor } from "@testing-library/react";

import { TestWrapper } from "@/test/wrapper";

mock.module("@/lib/api", () => ({
  api: {
    get: () => ({
      json: () =>
        Promise.resolve({ users: [{ id: "u1", email: "test@test.com" }], total: 1 }),
    }),
    post: () => ({
      json: () => Promise.resolve({ id: "u-new", email: "new@test.com" }),
    }),
    patch: () => ({
      json: () => Promise.resolve({ id: "u1", email: "updated@test.com" }),
    }),
  },
  HTTPError: class extends Error {},
}));

const { useUsers, useUser, useSearchUsers, useCreateUser, useUpdateUser } = await import(
  "@/hooks/use-users"
);

describe("useUsers", () => {
  it("fetches users", async () => {
    const { result } = renderHook(() => useUsers(), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data?.users).toHaveLength(1);
    expect(result.current.data?.users[0].email).toBe("test@test.com");
  });
});

describe("useUser", () => {
  it("fetches a single user", async () => {
    const { result } = renderHook(() => useUser("u1"), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });

  it("does not fetch when id is empty", () => {
    const { result } = renderHook(() => useUser(""), {
      wrapper: TestWrapper,
    });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useSearchUsers", () => {
  it("fetches when search is long enough", async () => {
    const { result } = renderHook(() => useSearchUsers("ab"), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data?.users).toHaveLength(1);
  });

  it("does not fetch when search is too short", () => {
    const { result } = renderHook(() => useSearchUsers("a"), {
      wrapper: TestWrapper,
    });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useCreateUser", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useCreateUser(), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useUpdateUser", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useUpdateUser(), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});
