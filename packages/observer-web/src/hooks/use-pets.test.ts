import { describe, expect, it, mock } from "bun:test";
import { renderHook, waitFor } from "@testing-library/react";

import { TestWrapper } from "@/test/wrapper";

mock.module("@/lib/api", () => ({
  api: {
    get: () => ({
      json: () => Promise.resolve({ pets: [{ id: "pet1", name: "Rex" }], total: 1 }),
    }),
    post: () => ({
      json: () => Promise.resolve({ id: "pet-new", name: "Rex" }),
    }),
    patch: () => ({
      json: () => Promise.resolve({ id: "pet1", name: "Rex Updated" }),
    }),
    delete: () => Promise.resolve(),
  },
  HTTPError: class extends Error {},
}));

const { usePets, usePet, useCreatePet, useUpdatePet, useDeletePet } = await import(
  "@/hooks/use-pets"
);

describe("usePets", () => {
  it("fetches pets for a project", async () => {
    const { result } = renderHook(() => usePets("proj1"), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data?.pets).toHaveLength(1);
    expect(result.current.data?.pets[0].name).toBe("Rex");
  });

  it("does not fetch when projectId is empty", () => {
    const { result } = renderHook(() => usePets(""), {
      wrapper: TestWrapper,
    });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("usePet", () => {
  it("fetches a single pet", async () => {
    const { result } = renderHook(() => usePet("proj1", "pet1"), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });

  it("does not fetch when id is empty", () => {
    const { result } = renderHook(() => usePet("proj1", ""), {
      wrapper: TestWrapper,
    });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useCreatePet", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useCreatePet("proj1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useUpdatePet", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useUpdatePet("proj1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useDeletePet", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useDeletePet("proj1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});
