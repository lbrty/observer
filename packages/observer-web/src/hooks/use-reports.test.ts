import { describe, expect, it, mock } from "bun:test";
import { renderHook, waitFor } from "@testing-library/react";

import { TestWrapper } from "@/test/wrapper";

mock.module("@/lib/api", () => ({
  api: {
    get: () => ({
      json: () =>
        Promise.resolve({
          consultations: { group: "consultations", rows: [], total: 0 },
        }),
    }),
  },
  HTTPError: class extends Error {},
}));

const { useReport } = await import("@/hooks/use-reports");

describe("useReport", () => {
  it("fetches report for a project", async () => {
    const { result } = renderHook(() => useReport("proj1"), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toHaveProperty("consultations");
  });

  it("does not fetch when projectId is empty", () => {
    const { result } = renderHook(() => useReport(""), {
      wrapper: TestWrapper,
    });
    expect(result.current.fetchStatus).toBe("idle");
  });
});
