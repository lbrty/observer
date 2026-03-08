import { describe, expect, it, mock } from "bun:test";
import { renderHook, waitFor } from "@testing-library/react";

import { TestWrapper } from "@/test/wrapper";

mock.module("@/lib/api", () => ({
  api: {
    get: () => ({
      json: () => Promise.resolve({ notes: [{ id: "n1", body: "test note", person_id: "p1" }] }),
    }),
    post: () => ({
      json: () => Promise.resolve({ id: "n-new", body: "new note", person_id: "p1" }),
    }),
    patch: () => ({
      json: () => Promise.resolve({ id: "n1", body: "updated note", person_id: "p1" }),
    }),
    delete: () => Promise.resolve(),
  },
  HTTPError: class extends Error {},
}));

const { useNotes, useCreateNote, useUpdateNote, useDeleteNote } = await import("@/hooks/use-notes");

describe("useNotes", () => {
  it("fetches notes for a person", async () => {
    const { result } = renderHook(() => useNotes("proj1", "p1"), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data?.notes).toHaveLength(1);
    expect(result.current.data?.notes[0].body).toBe("test note");
  });

  it("does not fetch when personId is empty", () => {
    const { result } = renderHook(() => useNotes("proj1", ""), {
      wrapper: TestWrapper,
    });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useCreateNote", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useCreateNote("proj1", "p1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useUpdateNote", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useUpdateNote("proj1", "p1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useDeleteNote", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useDeleteNote("proj1", "p1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});
