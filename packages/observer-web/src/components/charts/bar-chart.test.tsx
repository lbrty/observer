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

import { render, screen } from "@testing-library/react";

import { BarChart } from "./bar-chart";

describe("BarChart", () => {
  it("returns null when data is empty", () => {
    const { container } = render(<BarChart data={[]} />);
    expect(container.innerHTML).toBe("");
  });

  it("renders container div when data provided", () => {
    const data = [
      { label: "A", count: 5 },
      { label: "B", count: 3 },
    ];
    const { container } = render(<BarChart data={data} />);
    expect(container.querySelector("div")).toBeDefined();
    expect(container.querySelector("svg")).toBeDefined();
  });

  it("renders legend items when legend prop provided", () => {
    const data = [{ label: "A", count: 5 }];
    const legend = [
      { short: "A", full: "Alpha" },
      { short: "B", full: "Bravo" },
    ];
    render(<BarChart data={data} legend={legend} />);

    expect(screen.getByText("A")).toBeDefined();
    expect(screen.getByText(/Alpha/)).toBeDefined();
    expect(screen.getByText("B")).toBeDefined();
    expect(screen.getByText(/Bravo/)).toBeDefined();
  });
});
