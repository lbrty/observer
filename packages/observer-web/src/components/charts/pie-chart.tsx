import * as d3 from "d3";
import { useEffect, useRef } from "react";

import type { CountResult } from "@/types/report";

interface PieChartProps {
  data: CountResult[];
  width?: number;
  height?: number;
}

const COLORS = [
  "#6366f1",
  "#f59e0b",
  "#10b981",
  "#ef4444",
  "#8b5cf6",
  "#ec4899",
  "#14b8a6",
  "#f97316",
];

export function PieChart({ data, width = 300, height = 300 }: PieChartProps) {
  const ref = useRef<SVGSVGElement>(null);

  useEffect(() => {
    if (!ref.current || data.length === 0) return;

    const svg = d3.select(ref.current);
    svg.selectAll("*").remove();

    const radius = Math.min(width, height) / 2 - 10;

    const g = svg
      .attr("viewBox", `0 0 ${width} ${height}`)
      .append("g")
      .attr("transform", `translate(${width / 2},${height / 2})`);

    const color = d3
      .scaleOrdinal<string>()
      .domain(data.map((d) => d.label))
      .range(COLORS);

    const pie = d3
      .pie<CountResult>()
      .value((d) => d.count)
      .sort(null);

    const arc = d3
      .arc<d3.PieArcDatum<CountResult>>()
      .innerRadius(radius * 0.4)
      .outerRadius(radius);

    const labelArc = d3
      .arc<d3.PieArcDatum<CountResult>>()
      .innerRadius(radius * 0.7)
      .outerRadius(radius * 0.7);

    g.selectAll(".slice")
      .data(pie(data))
      .join("path")
      .attr("class", "slice")
      .attr("d", arc)
      .attr("fill", (d) => color(d.data.label))
      .attr("stroke", "var(--color-bg, #fff)")
      .attr("stroke-width", 2);

    g.selectAll(".pie-label")
      .data(pie(data))
      .join("text")
      .attr("class", "pie-label")
      .attr("transform", (d) => `translate(${labelArc.centroid(d)})`)
      .attr("text-anchor", "middle")
      .style("font-size", "11px")
      .style("fill", "#fff")
      .style("font-weight", "600")
      .text((d) => (d.data.count > 0 ? d.data.count : ""));
  }, [data, width, height]);

  if (data.length === 0) return null;

  return (
    <div className="flex items-start gap-4">
      <svg ref={ref} className="w-48 shrink-0" />
      <ul className="space-y-1 pt-4 text-sm">
        {data.map((d, i) => (
          <li key={d.label} className="flex items-center gap-2">
            <span
              className="inline-block size-2.5 rounded-full"
              style={{ backgroundColor: COLORS[i % COLORS.length] }}
            />
            <span className="text-fg-secondary">
              {d.label} ({d.count})
            </span>
          </li>
        ))}
      </ul>
    </div>
  );
}
