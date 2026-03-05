import * as d3 from "d3";
import { useEffect, useRef, useState } from "react";

import type { CountResult } from "@/types/report";

interface Tooltip {
  visible: boolean;
  x: number;
  y: number;
  label: string;
  count: number;
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

export function PieChart({ data }: { data: CountResult[] }) {
  const containerRef = useRef<HTMLDivElement>(null);
  const svgRef = useRef<SVGSVGElement>(null);
  const [size, setSize] = useState(220);
  const [tooltip, setTooltip] = useState<Tooltip>({
    visible: false,
    x: 0,
    y: 0,
    label: "",
    count: 0,
  });
  const [selectedLabel, setSelectedLabel] = useState<string | null>(null);

  useEffect(() => {
    if (!containerRef.current) return;
    const observer = new ResizeObserver((entries) => {
      const w = entries[0].contentRect.width;
      setSize(Math.min(w * 0.7, 425));
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
      .sort(null)
      .padAngle(0.02);

    const arc = d3
      .arc<d3.PieArcDatum<CountResult>>()
      .innerRadius(radius * 0.45)
      .outerRadius(radius)
      .cornerRadius(2);

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
      .attr("opacity", (d) => (selectedLabel === null || selectedLabel === d.data.label ? 1 : 0.3))
      .style("cursor", "pointer")
      .on("mouseover", function (event: MouseEvent, d: d3.PieArcDatum<CountResult>) {
        const bounds = containerRef.current?.getBoundingClientRect();
        if (!bounds) return;
        setTooltip({
          visible: true,
          x: event.clientX - bounds.left + 12,
          y: event.clientY - bounds.top - 10,
          label: d.data.label,
          count: d.data.count,
        });
      })
      .on("mousemove", function (event: MouseEvent, d: d3.PieArcDatum<CountResult>) {
        const bounds = containerRef.current?.getBoundingClientRect();
        if (!bounds) return;
        setTooltip({
          visible: true,
          x: event.clientX - bounds.left + 12,
          y: event.clientY - bounds.top - 10,
          label: d.data.label,
          count: d.data.count,
        });
      })
      .on("mouseout", function () {
        setTooltip((prev) => ({ ...prev, visible: false }));
      })
      .on("click", function (_event: MouseEvent, d: d3.PieArcDatum<CountResult>) {
        setSelectedLabel((prev) => (prev === d.data.label ? null : d.data.label));
      });

    g.selectAll(".pie-label")
      .data(pie(data))
      .join("text")
      .attr("class", "pie-label")
      .attr("transform", (d) => `translate(${labelArc.centroid(d)})`)
      .attr("text-anchor", "middle")
      .style("font-size", "9px")
      .style("fill", "#fff")
      .style("font-weight", "600")
      .attr("opacity", (d) => (selectedLabel === null || selectedLabel === d.data.label ? 1 : 0.3))
      .text((d) => (d.data.count > 0 ? d.data.count : ""));
  }, [data, size, selectedLabel]);

  if (data.length === 0) return null;

  return (
    <div ref={containerRef} className="relative flex items-center gap-6">
      <svg ref={svgRef} className="shrink-0" />
      {tooltip.visible && (
        <div
          style={{
            position: "absolute",
            left: tooltip.x,
            top: tooltip.y,
            background: "var(--bg-secondary)",
            border: "1px solid var(--border-secondary)",
            borderRadius: 8,
            padding: "6px 10px",
            fontSize: 12,
            color: "var(--fg)",
            boxShadow: "var(--shadow-elevated)",
            pointerEvents: "none",
          }}
        >
          <strong>{tooltip.label}</strong>: {tooltip.count}
        </div>
      )}
      <ul className="space-y-1.5 text-sm">
        {data.map((d, i) => (
          <li
            key={d.label}
            className="flex items-center gap-2"
            style={{
              opacity: selectedLabel === null || selectedLabel === d.label ? 1 : 0.3,
              cursor: "pointer",
            }}
            onClick={() => setSelectedLabel((prev) => (prev === d.label ? null : d.label))}
          >
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
