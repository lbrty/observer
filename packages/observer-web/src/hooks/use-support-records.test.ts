import { describe, expect, it, mock } from "bun:test";
import { renderHook, waitFor } from "@testing-library/react";

import { TestWrapper } from "@/test/wrapper";

mock.module("@/lib/api", () => ({
  api: {
    get: () => ({
      json: () =>
        Promise.resolve({ records: [{ id: "sr1", type: "consultation" }], total: 1 }),
    }),
    post: () => ({
      json: () => Promise.resolve({ id: "sr-new", type: "consultation" }),
    }),
    patch: () => ({
      json: () => Promise.resolve({ id: "sr1", type: "consultation" }),
    }),
  },
  HTTPError: class extends Error {},
}));

const { useSupportRecords, useSupportRecord, useCreateSupportRecord, useUpdateSupportRecord } =
  await import("@/hooks/use-support-records");

describe("useSupportRecords", () => {
  it("fetches support records for a project", async () => {
    const { result } = renderHook(() => useSupportRecords("proj1"), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data?.records).toHaveLength(1);
    expect(result.current.data?.records[0].type).toBe("consultation");
  });

  it("does not fetch when projectId is empty", () => {
    const { result } = renderHook(() => useSupportRecords(""), {
      wrapper: TestWrapper,
    });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useSupportRecord", () => {
  it("fetches a single support record", async () => {
    const { result } = renderHook(() => useSupportRecord("proj1", "sr1"), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });

  it("does not fetch when id is empty", () => {
    const { result } = renderHook(() => useSupportRecord("proj1", ""), {
      wrapper: TestWrapper,
    });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useCreateSupportRecord", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useCreateSupportRecord("proj1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useUpdateSupportRecord", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useUpdateSupportRecord("proj1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});
