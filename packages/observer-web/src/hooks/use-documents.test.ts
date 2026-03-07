import { describe, expect, it, mock } from "bun:test";
import { renderHook, waitFor } from "@testing-library/react";

import { TestWrapper } from "@/test/wrapper";

mock.module("@/lib/api", () => ({
  api: {
    get: () => ({
      json: () =>
        Promise.resolve({
          documents: [{ id: "d1", name: "file.pdf", person_id: "p1", project_id: "proj1" }],
        }),
    }),
    post: () => ({
      json: () =>
        Promise.resolve({ id: "d-new", name: "new.pdf", person_id: "p1", project_id: "proj1" }),
    }),
    patch: () => ({
      json: () =>
        Promise.resolve({
          id: "d1",
          name: "updated.pdf",
          person_id: "p1",
          project_id: "proj1",
        }),
    }),
    delete: () => Promise.resolve(),
  },
}));

const { useDocuments, useUploadDocument, useUpdateDocument, useDeleteDocument } = await import(
  "@/hooks/use-documents"
);

describe("useDocuments", () => {
  it("fetches documents for a person", async () => {
    const { result } = renderHook(() => useDocuments("proj1", "p1"), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data?.documents).toHaveLength(1);
    expect(result.current.data?.documents[0].name).toBe("file.pdf");
  });

  it("does not fetch when projectId is empty", () => {
    const { result } = renderHook(() => useDocuments("", "p1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.fetchStatus).toBe("idle");
  });

  it("does not fetch when personId is empty", () => {
    const { result } = renderHook(() => useDocuments("proj1", ""), {
      wrapper: TestWrapper,
    });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useUploadDocument", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useUploadDocument("proj1", "p1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useUpdateDocument", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useUpdateDocument("proj1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useDeleteDocument", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useDeleteDocument("proj1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});
