import { describe, expect, it, mock } from "bun:test";

mock.module("d3", () => {
  const chainable: any = new Proxy(() => chainable, {
    get: () => chainable,
    apply: () => chainable,
  });
  return {
    select: () => chainable,
    scaleBand: () => chainable,
    scaleLinear: () => chainable,
    axisBottom: () => chainable,
    axisLeft: () => chainable,
    max: () => 10,
    pie: () => chainable,
    arc: () => chainable,
  };
});

mock.module("./colors", () => ({
  getColor: (_label: string, _map: unknown, i: number) => `#color${i}`,
}));

globalThis.ResizeObserver = class {
  observe() {}
  unobserve() {}
  disconnect() {}
} as any;

import { render, screen } from "@testing-library/react";

import { PieChart } from "./pie-chart";

describe("PieChart", () => {
  it("returns null when data is empty", () => {
    const { container } = render(<PieChart data={[]} />);
    expect(container.innerHTML).toBe("");
  });

  it("renders legend items when data provided", () => {
    const data = [
      { label: "Female", count: 12 },
      { label: "Male", count: 8 },
    ];
    render(<PieChart data={data} />);

    expect(screen.getByText(/Female/)).toBeDefined();
    expect(screen.getByText(/12/)).toBeDefined();
    expect(screen.getByText(/Male/)).toBeDefined();
    expect(screen.getByText(/8/)).toBeDefined();
  });
});
