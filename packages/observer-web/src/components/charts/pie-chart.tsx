import * as d3 from "d3";
import { useEffect, useRef, useState } from "react";

import type { CountResult } from "@/types/report";

interface PieChartProps {
  data: CountResult[];
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

export function PieChart({ data }: PieChartProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const svgRef = useRef<SVGSVGElement>(null);
  const [size, setSize] = useState(220);

  useEffect(() => {
    if (!containerRef.current) return;
    const observer = new ResizeObserver((entries) => {
      const w = entries[0].contentRect.width;
      setSize(Math.min(w * 0.55, 340));
    });
    observer.observe(containerRef.current);
    return () => observer.disconnect();
  }, []);

  useEffect(() => {
    if (!svgRef.current || data.length === 0) return;

    const svg = d3.select(svgRef.current);
    svg.selectAll("*").remove();

    const dim = size;
    const radius = dim / 2 - 4;

    const g = svg
      .attr("viewBox", `0 0 ${dim} ${dim}`)
      .attr("width", dim)
      .attr("height", dim)
      .append("g")
      .attr("transform", `translate(${dim / 2},${dim / 2})`);

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
      .innerRadius(radius * 0.45)
      .outerRadius(radius);

    const labelArc = d3
      .arc<d3.PieArcDatum<CountResult>>()
      .innerRadius(radius * 0.72)
      .outerRadius(radius * 0.72);

    g.selectAll(".slice")
      .data(pie(data))
      .join("path")
      .attr("class", "slice")
      .attr("d", arc)
      .attr("fill", (d) => color(d.data.label))
      .attr("stroke", "var(--bg, #fff)")
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
  }, [data, size]);

  if (data.length === 0) return null;

  return (
    <div ref={containerRef} className="flex items-center gap-6">
      <svg ref={svgRef} className="shrink-0" />
      <ul className="space-y-1.5 text-sm">
        {data.map((d, i) => (
          <li key={d.label} className="flex items-center gap-2">
            <span
              className="inline-block size-2.5 rounded-full shrink-0"
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
