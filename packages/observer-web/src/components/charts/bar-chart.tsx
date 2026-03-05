import * as d3 from "d3";
import { useEffect, useRef } from "react";

import type { CountResult } from "@/types/report";

interface BarChartProps {
  data: CountResult[];
  width?: number;
  height?: number;
}

export function BarChart({ data, width = 500, height = 300 }: BarChartProps) {
  const ref = useRef<SVGSVGElement>(null);

  useEffect(() => {
    if (!ref.current || data.length === 0) return;

    const svg = d3.select(ref.current);
    svg.selectAll("*").remove();

    const margin = { top: 20, right: 20, bottom: 60, left: 50 };
    const w = width - margin.left - margin.right;
    const h = height - margin.top - margin.bottom;

    const g = svg
      .attr("viewBox", `0 0 ${width} ${height}`)
      .append("g")
      .attr("transform", `translate(${margin.left},${margin.top})`);

    const x = d3
      .scaleBand<string>()
      .domain(data.map((d) => d.label))
      .range([0, w])
      .padding(0.3);

    const y = d3
      .scaleLinear()
      .domain([0, d3.max(data, (d) => d.count) ?? 0])
      .nice()
      .range([h, 0]);

    g.append("g")
      .attr("transform", `translate(0,${h})`)
      .call(d3.axisBottom(x))
      .selectAll("text")
      .attr("transform", "rotate(-35)")
      .style("text-anchor", "end")
      .style("font-size", "11px");

    g.append("g").call(d3.axisLeft(y).ticks(5)).style("font-size", "11px");

    g.selectAll(".bar")
      .data(data)
      .join("rect")
      .attr("class", "bar")
      .attr("x", (d) => x(d.label) ?? 0)
      .attr("y", (d) => y(d.count))
      .attr("width", x.bandwidth())
      .attr("height", (d) => h - y(d.count))
      .attr("rx", 3)
      .attr("fill", "var(--color-accent, #6366f1)");

    g.selectAll(".label")
      .data(data)
      .join("text")
      .attr("class", "label")
      .attr("x", (d) => (x(d.label) ?? 0) + x.bandwidth() / 2)
      .attr("y", (d) => y(d.count) - 4)
      .attr("text-anchor", "middle")
      .style("font-size", "11px")
      .style("fill", "var(--color-fg-secondary, #6b7280)")
      .text((d) => d.count);
  }, [data, width, height]);

  if (data.length === 0) return null;

  return <svg ref={ref} className="w-full" />;
}
