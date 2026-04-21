import { describe, expect, it, mock } from "bun:test";
import { renderHook, waitFor } from "@testing-library/react";

import { TestWrapper } from "@/test/wrapper";

mock.module("@/lib/api", () => ({
  api: {
    get: () => ({
      json: () =>
        Promise.resolve({
          projects: [{ id: "proj1", name: "My Project", role: "manager" }],
        }),
    }),
  },
  HTTPError: class extends Error {},
}));

const { useMyProjects } = await import("@/hooks/use-my-projects");

describe("useMyProjects", () => {
  it("fetches my projects", async () => {
    const { result } = renderHook(() => useMyProjects(), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data?.projects).toHaveLength(1);
    expect(result.current.data?.projects[0].name).toBe("My Project");
  });
});
