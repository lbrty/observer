import { describe, expect, it, mock } from "bun:test";
import { renderHook, waitFor } from "@testing-library/react";

import { TestWrapper } from "@/test/wrapper";

mock.module("@/lib/api", () => ({
  api: {
    get: () => ({
      json: () =>
        Promise.resolve({
          by_status: { group: "by_status", rows: [], total: 0 },
        }),
    }),
  },
  HTTPError: class extends Error {},
}));

const { usePetReport } = await import("@/hooks/use-pet-reports");

describe("usePetReport", () => {
  it("fetches pet report for a project", async () => {
    const { result } = renderHook(() => usePetReport("proj1"), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toHaveProperty("by_status");
  });

  it("does not fetch when projectId is empty", () => {
    const { result } = renderHook(() => usePetReport(""), {
      wrapper: TestWrapper,
    });
    expect(result.current.fetchStatus).toBe("idle");
  });
});
