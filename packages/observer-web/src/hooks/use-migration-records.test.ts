import { describe, expect, it, mock, beforeEach } from "bun:test";
import { renderHook, waitFor } from "@testing-library/react";

import { TestWrapper } from "@/test/wrapper";

// Mock the api module
mock.module("@/lib/api", () => ({
  api: {
    get: () => ({
      json: () => Promise.resolve({ records: [{ id: "mr1", person_id: "p1" }] }),
    }),
    post: () => ({
      json: () => Promise.resolve({ id: "mr-new", person_id: "p1" }),
    }),
    patch: () => ({
      json: () => Promise.resolve({ id: "mr1", person_id: "p1", movement_reason: "economic" }),
    }),
  },
  HTTPError: class extends Error {},
}));

// Import after mocking
const { useMigrationRecords, useCreateMigrationRecord, useUpdateMigrationRecord } = await import(
  "@/hooks/use-migration-records"
);

describe("useMigrationRecords", () => {
  it("fetches migration records for a person", async () => {
    const { result } = renderHook(() => useMigrationRecords("proj1", "p1"), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data?.records).toHaveLength(1);
    expect(result.current.data?.records[0].id).toBe("mr1");
  });

  it("does not fetch when projectId is empty", () => {
    const { result } = renderHook(() => useMigrationRecords("", "p1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.fetchStatus).toBe("idle");
  });

  it("does not fetch when personId is empty", () => {
    const { result } = renderHook(() => useMigrationRecords("proj1", ""), {
      wrapper: TestWrapper,
    });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useCreateMigrationRecord", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useCreateMigrationRecord("proj1", "p1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
    expect(result.current.isPending).toBe(false);
  });
});

describe("useUpdateMigrationRecord", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useUpdateMigrationRecord("proj1", "p1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
    expect(result.current.isPending).toBe(false);
  });
});
