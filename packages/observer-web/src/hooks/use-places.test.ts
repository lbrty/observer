import { describe, expect, it, mock } from "bun:test";
import { renderHook, waitFor } from "@testing-library/react";

import { TestWrapper } from "@/test/wrapper";

mock.module("@/lib/api", () => ({
  api: {
    get: () => ({
      json: () => Promise.resolve({ places: [{ id: "pl1", name: "Kyiv" }] }),
    }),
    post: () => ({
      json: () => Promise.resolve({ id: "pl-new", name: "New Place" }),
    }),
    patch: () => ({
      json: () => Promise.resolve({ id: "pl1", name: "Updated" }),
    }),
    delete: () => ({ json: () => Promise.resolve() }),
  },
  HTTPError: class extends Error {},
}));

const { usePlaces, useCreatePlace, useUpdatePlace, useDeletePlace } = await import(
  "@/hooks/use-places"
);

describe("usePlaces", () => {
  it("fetches places", async () => {
    const { result } = renderHook(() => usePlaces(), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data?.places).toHaveLength(1);
    expect(result.current.data?.places[0].name).toBe("Kyiv");
  });

  it("fetches places with stateId filter", async () => {
    const { result } = renderHook(() => usePlaces("s1"), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });
});

describe("useCreatePlace", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useCreatePlace(), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useUpdatePlace", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useUpdatePlace(), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useDeletePlace", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useDeletePlace(), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});
