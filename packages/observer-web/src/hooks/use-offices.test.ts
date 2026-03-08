import { describe, expect, it, mock } from "bun:test";
import { renderHook, waitFor } from "@testing-library/react";

import { TestWrapper } from "@/test/wrapper";

mock.module("@/lib/api", () => ({
  api: {
    get: () => ({
      json: () => Promise.resolve({ offices: [{ id: "o1", name: "Main Office" }] }),
    }),
    post: () => ({
      json: () => Promise.resolve({ id: "o-new", name: "New Office" }),
    }),
    patch: () => ({
      json: () => Promise.resolve({ id: "o1", name: "Updated" }),
    }),
    delete: () => ({ json: () => Promise.resolve() }),
  },
  HTTPError: class extends Error {},
}));

const { useOffices, useCreateOffice, useUpdateOffice, useDeleteOffice } = await import(
  "@/hooks/use-offices"
);

describe("useOffices", () => {
  it("fetches offices", async () => {
    const { result } = renderHook(() => useOffices(), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toHaveLength(1);
    expect(result.current.data![0].name).toBe("Main Office");
  });
});

describe("useCreateOffice", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useCreateOffice(), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useUpdateOffice", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useUpdateOffice(), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useDeleteOffice", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useDeleteOffice(), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});
