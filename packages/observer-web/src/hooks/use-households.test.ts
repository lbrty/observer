import { describe, expect, it, mock } from "bun:test";
import { renderHook, waitFor } from "@testing-library/react";

import { TestWrapper } from "@/test/wrapper";

mock.module("@/lib/api", () => ({
  api: {
    get: () => ({
      json: () =>
        Promise.resolve({ households: [{ id: "h1", member_count: 3 }], total: 1 }),
    }),
    post: () => ({
      json: () => Promise.resolve({ id: "h-new", member_count: 0 }),
    }),
    patch: () => ({
      json: () => Promise.resolve({ id: "h1", member_count: 3 }),
    }),
    delete: () => Promise.resolve(),
  },
  HTTPError: class extends Error {},
}));

const {
  useHouseholds,
  useHousehold,
  useCreateHousehold,
  useUpdateHousehold,
  useAddHouseholdMember,
  useRemoveHouseholdMember,
} = await import("@/hooks/use-households");

describe("useHouseholds", () => {
  it("fetches households for a project", async () => {
    const { result } = renderHook(() => useHouseholds("proj1"), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data?.households).toHaveLength(1);
    expect(result.current.data?.households[0].member_count).toBe(3);
  });

  it("does not fetch when projectId is empty", () => {
    const { result } = renderHook(() => useHouseholds(""), {
      wrapper: TestWrapper,
    });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useHousehold", () => {
  it("fetches a single household", async () => {
    const { result } = renderHook(() => useHousehold("proj1", "h1"), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });

  it("does not fetch when id is empty", () => {
    const { result } = renderHook(() => useHousehold("proj1", ""), {
      wrapper: TestWrapper,
    });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useCreateHousehold", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useCreateHousehold("proj1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useUpdateHousehold", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useUpdateHousehold("proj1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useAddHouseholdMember", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useAddHouseholdMember("proj1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useRemoveHouseholdMember", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useRemoveHouseholdMember("proj1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});
