import { describe, expect, it, mock } from "bun:test";
import { renderHook, waitFor } from "@testing-library/react";

import { TestWrapper } from "@/test/wrapper";

mock.module("@/lib/api", () => ({
  api: {
    get: () => ({
      json: () => Promise.resolve({ categories: [{ id: "c1", name: "Social" }] }),
    }),
    post: () => ({
      json: () => Promise.resolve({ id: "c-new", name: "New Category" }),
    }),
    patch: () => ({
      json: () => Promise.resolve({ id: "c1", name: "Updated" }),
    }),
    delete: () => ({ json: () => Promise.resolve() }),
  },
  HTTPError: class extends Error {},
}));

const { useCategories, useCreateCategory, useUpdateCategory, useDeleteCategory } = await import(
  "@/hooks/use-categories"
);

describe("useCategories", () => {
  it("fetches categories", async () => {
    const { result } = renderHook(() => useCategories(), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toHaveLength(1);
    expect(result.current.data![0].name).toBe("Social");
  });
});

describe("useCreateCategory", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useCreateCategory(), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useUpdateCategory", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useUpdateCategory(), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useDeleteCategory", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useDeleteCategory(), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});
