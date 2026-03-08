import { describe, expect, it, mock } from "bun:test";
import { renderHook, waitFor } from "@testing-library/react";

import { TestWrapper } from "@/test/wrapper";

mock.module("@/lib/api", () => ({
  api: {
    get: () => ({
      json: () => Promise.resolve({ tags: [{ id: "t1", name: "urgent" }], tag_ids: ["t1"] }),
    }),
    post: () => ({
      json: () => Promise.resolve({ id: "t-new", name: "new tag" }),
    }),
    patch: () => ({
      json: () => Promise.resolve({ id: "t1", name: "updated" }),
    }),
    put: () => ({
      json: () => Promise.resolve({ tag_ids: ["t1", "t2"] }),
    }),
    delete: () => Promise.resolve(),
  },
  HTTPError: class extends Error {},
}));

const {
  useTags,
  useCreateTag,
  useUpdateTag,
  useDeleteTag,
  usePersonTags,
  useReplacePersonTags,
  usePetTags,
  useReplacePetTags,
} = await import("@/hooks/use-tags");

describe("useTags", () => {
  it("fetches tags for a project", async () => {
    const { result } = renderHook(() => useTags("proj1"), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data?.tags).toHaveLength(1);
    expect(result.current.data?.tags[0].name).toBe("urgent");
  });

  it("does not fetch when projectId is empty", () => {
    const { result } = renderHook(() => useTags(""), {
      wrapper: TestWrapper,
    });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("usePersonTags", () => {
  it("fetches tags for a person", async () => {
    const { result } = renderHook(() => usePersonTags("proj1", "p1"), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data?.tag_ids).toContain("t1");
  });

  it("does not fetch when personId is empty", () => {
    const { result } = renderHook(() => usePersonTags("proj1", ""), {
      wrapper: TestWrapper,
    });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("usePetTags", () => {
  it("fetches tags for a pet", async () => {
    const { result } = renderHook(() => usePetTags("proj1", "pet1"), {
      wrapper: TestWrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data?.tag_ids).toContain("t1");
  });

  it("does not fetch when petId is empty", () => {
    const { result } = renderHook(() => usePetTags("proj1", ""), {
      wrapper: TestWrapper,
    });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useCreateTag", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useCreateTag("proj1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useUpdateTag", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useUpdateTag("proj1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useDeleteTag", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useDeleteTag("proj1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useReplacePersonTags", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useReplacePersonTags("proj1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});

describe("useReplacePetTags", () => {
  it("returns a mutation hook", () => {
    const { result } = renderHook(() => useReplacePetTags("proj1"), {
      wrapper: TestWrapper,
    });
    expect(result.current.mutateAsync).toBeFunction();
  });
});
