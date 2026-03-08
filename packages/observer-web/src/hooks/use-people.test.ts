import { describe, expect, it, mock } from "bun:test";
import { renderHook, waitFor } from "@testing-library/react";

import { TestWrapper } from "@/test/wrapper";

mock.module("@/lib/api", () => ({
  api: {
    get: () => ({
      json: () => Promise.resolve({ people: [{ id: "p1", first_name: "Aida" }], total: 1 }),
    }),
    post: () => ({
      json: () => Promise.resolve({ id: "p-new", first_name: "Aida" }),
    }),
    patch: () => ({
      json: () => Promise.resolve({ id: "p1", first_name: "Aida Updated" }),
    }),
  },
  HTTPError: class extends Error {},
}));

const { usePeople, usePerson, useCreatePerson, useUpdatePerson, useSearchPeople } = await import(
  "@/hooks/use-people"
);

describe("usePeople", () => {
  it("fetches people for a project", async () => {
    const { result } = renderHook(() => usePeople("proj1"), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data?.people).toHaveLength(1);
    expect(result.current.data?.people[0].first_name).toBe("Aida");
  });

  it("does not fetch when projectId is empty", () => {
    const { result } = renderHook(() => usePeople(""), {
      wrapper: TestWrapper,
    });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("usePerson", () => {
  it("fetches a single person", async () => {
    const { result } = renderHook(() => usePerson("proj1", "p1"), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });

  it("does not fetch when personId is empty", () => {
    const { result } = renderHook(() => usePerson("proj1", ""), {
      wrapper: TestWrapper,
    });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useCreatePerson", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useCreatePerson("proj1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useUpdatePerson", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useUpdatePerson("proj1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useSearchPeople", () => {
  it("fetches when search is long enough", async () => {
    const { result } = renderHook(() => useSearchPeople("proj1", "ai"), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data?.people).toHaveLength(1);
  });

  it("does not fetch when search is too short", () => {
    const { result } = renderHook(() => useSearchPeople("proj1", "a"), {
      wrapper: TestWrapper,
    });
    expect(result.current.fetchStatus).toBe("idle");
  });
});
