import { describe, expect, it, mock } from "bun:test";
import { renderHook, waitFor } from "@testing-library/react";

import { TestWrapper } from "@/test/wrapper";

mock.module("@/lib/api", () => ({
  api: {
    get: () => ({
      json: () => Promise.resolve({ countries: [{ id: "co1", name: "Ukraine", code: "UA" }] }),
    }),
    post: () => ({
      json: () => Promise.resolve({ id: "co-new", name: "Poland", code: "PL" }),
    }),
    patch: () => ({
      json: () => Promise.resolve({ id: "co1", name: "Updated", code: "UA" }),
    }),
    delete: () => ({ json: () => Promise.resolve() }),
  },
  HTTPError: class extends Error {},
}));

const { useCountries, useCreateCountry, useUpdateCountry, useDeleteCountry } = await import(
  "@/hooks/use-countries"
);

describe("useCountries", () => {
  it("fetches countries", async () => {
    const { result } = renderHook(() => useCountries(), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toHaveLength(1);
    expect(result.current.data![0].name).toBe("Ukraine");
  });
});

describe("useCreateCountry", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useCreateCountry(), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useUpdateCountry", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useUpdateCountry(), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useDeleteCountry", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useDeleteCountry(), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});
