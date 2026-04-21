import { afterEach, describe, expect, it } from "bun:test";
import { cleanup, render, screen } from "@testing-library/react";

import { StatusBadge, StatusDot } from "@/components/status-badge";

afterEach(cleanup);

describe("StatusBadge", () => {
  it("renders label text", () => {
    render(<StatusBadge label="active" />);

    expect(screen.getByText("active")).toBeDefined();
  });

  it("applies rose variant for admin role", () => {
    const { container } = render(<StatusBadge label="admin" />);
    const badge = container.firstElementChild as HTMLElement;

    expect(badge.className).toContain("bg-rose/15");
    expect(badge.className).toContain("text-rose");
  });

  it("applies foam variant for active status", () => {
    const { container } = render(<StatusBadge label="active" />);
    const badge = container.firstElementChild as HTMLElement;

    expect(badge.className).toContain("bg-foam/15");
    expect(badge.className).toContain("text-foam");
  });

  it("defaults to neutral for unknown label", () => {
    const { container } = render(<StatusBadge label="something-unknown" />);
    const badge = container.firstElementChild as HTMLElement;

    expect(badge.className).toContain("bg-bg-tertiary");
    expect(badge.className).toContain("text-fg-secondary");
  });

  it("hides dot when dot={false}", () => {
    const { container } = render(<StatusBadge label="active" dot={false} />);
    const dots = container.querySelectorAll(".rounded-full.size-1\\.5");

    expect(dots.length).toBe(0);
  });

  it("shows dot by default", () => {
    const { container } = render(<StatusBadge label="active" />);
    const badge = container.firstElementChild as HTMLElement;
    const dot = badge.querySelector("span span");

    expect(dot).toBeDefined();
  });
});

describe("StatusDot", () => {
  it("renders with foam color when active", () => {
    const { container } = render(<StatusDot active />);
    const dot = container.firstElementChild as HTMLElement;

    expect(dot.className).toContain("bg-foam");
  });

  it("renders with tertiary color when inactive", () => {
    const { container } = render(<StatusDot active={false} />);
    const dot = container.firstElementChild as HTMLElement;

    expect(dot.className).toContain("bg-fg-tertiary");
  });
});
