import { describe, expect, it } from "bun:test";

describe("test infrastructure", () => {
  it("runs bun tests", () => {
    expect(1 + 1).toBe(2);
  });

  it("has DOM available", () => {
    const el = document.createElement("div");
    el.textContent = "hello";
    expect(el.textContent).toBe("hello");
  });
});
