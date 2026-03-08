import { describe, expect, it, mock } from "bun:test";
import { renderHook, waitFor } from "@testing-library/react";

import { TestWrapper } from "@/test/wrapper";

mock.module("@/lib/api", () => ({
  api: {
    get: () => ({
      json: () => Promise.resolve({ states: [{ id: "s1", name: "Kyiv Oblast" }] }),
    }),
    post: () => ({
      json: () => Promise.resolve({ id: "s-new", name: "New State" }),
    }),
    patch: () => ({
      json: () => Promise.resolve({ id: "s1", name: "Updated" }),
    }),
    delete: () => ({ json: () => Promise.resolve() }),
  },
  HTTPError: class extends Error {},
}));

const { useStates, useCreateState, useUpdateState, useDeleteState } = await import(
  "@/hooks/use-states"
);

describe("useStates", () => {
  it("fetches states", async () => {
    const { result } = renderHook(() => useStates(), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data?.states).toHaveLength(1);
    expect(result.current.data?.states[0].name).toBe("Kyiv Oblast");
  });

  it("fetches states with countryId filter", async () => {
    const { result } = renderHook(() => useStates("co1"), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });
});

describe("useCreateState", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useCreateState(), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useUpdateState", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useUpdateState(), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useDeleteState", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useDeleteState(), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});
