import * as d3 from "d3";
import { useEffect, useRef, useState } from "react";

import type { CountResult } from "@/types/report";

import { getColor } from "./colors";

interface Tooltip {
  visible: boolean;
  x: number;
  y: number;
  label: string;
  count: number;
}

export interface BarLegendItem {
  short: string;
  full: string;
}

interface BarChartProps {
  data: CountResult[];
  width?: number;
  height?: number;
  yAxisLabel?: string;
  legend?: BarLegendItem[];
  direction?: "vertical" | "horizontal" | "auto";
  colorMap?: Record<string, string>;
}

export function BarChart({
  data,
  width = 500,
  height = 300,
  yAxisLabel = "Count",
  legend,
  direction,
  colorMap,
}: BarChartProps) {
  const ref = useRef<SVGSVGElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const [tooltip, setTooltip] = useState<Tooltip>({
    visible: false,
    x: 0,
    y: 0,
    label: "",
    count: 0,
  });
  const [selectedLabel, setSelectedLabel] = useState<string | null>(null);

  const resolvedDirection =
    direction === "auto" ? (data.length > 6 ? "horizontal" : "vertical") : (direction ?? "vertical");

  const resolvedHeight =
    resolvedDirection === "horizontal" ? Math.max(300, data.length * 32) : height;

  useEffect(() => {
    if (!ref.current || data.length === 0) return;

    const svg = d3.select(ref.current);
    svg.selectAll("*").remove();

    const axisColor = "var(--fg-tertiary, #6b7280)";

    const handleMouseOver = function (event: MouseEvent, d: CountResult) {
      const bounds = containerRef.current?.getBoundingClientRect();
      if (!bounds) return;
      setTooltip({
        visible: true,
        x: event.clientX - bounds.left + 12,
        y: event.clientY - bounds.top - 10,
        label: d.label,
        count: d.count,
      });
    };

    const handleMouseOut = function () {
      setTooltip((prev) => ({ ...prev, visible: false }));
    };

    const handleClick = function (_event: MouseEvent, d: CountResult) {
      setSelectedLabel((prev) => (prev === d.label ? null : d.label));
    };

    const opacityFn = (d: CountResult) =>
      selectedLabel === null || selectedLabel === d.label ? 1 : 0.3;

    if (resolvedDirection === "horizontal") {
      const margin = { top: 20, right: 40, bottom: 20, left: 120 };
      const w = width - margin.left - margin.right;
      const h = resolvedHeight - margin.top - margin.bottom;

      const g = svg
        .attr("viewBox", `0 0 ${width} ${resolvedHeight}`)
        .append("g")
        .attr("transform", `translate(${margin.left},${margin.top})`);

      const yBand = d3
        .scaleBand<string>()
        .domain(data.map((d) => d.label))
        .range([0, h])
        .padding(0.3);

      const xLinear = d3
        .scaleLinear()
        .domain([0, d3.max(data, (d) => d.count) ?? 0])
        .nice()
        .range([0, w]);

      const xAxis = g.append("g").attr("transform", `translate(0,${h})`).call(d3.axisBottom(xLinear).ticks(5));
      xAxis.selectAll("text").style("font-size", "9px").style("fill", axisColor);
      xAxis.selectAll("line").style("stroke", axisColor);
      xAxis.select(".domain").style("stroke", axisColor);

      const yAxis = g.append("g").call(d3.axisLeft(yBand));
      yAxis.selectAll("text").style("font-size", "9px").style("fill", axisColor);
      yAxis.selectAll("line").style("stroke", axisColor);
      yAxis.select(".domain").style("stroke", axisColor);

      g.append("clipPath")
        .attr("id", "bar-clip")
        .append("rect")
        .attr("x", 0)
        .attr("y", 0)
        .attr("width", w)
        .attr("height", h);

      const barGroup = g.append("g").attr("clip-path", "url(#bar-clip)");

      barGroup
        .selectAll(".bar")
        .data(data)
        .join("rect")
        .attr("class", "bar")
        .attr("x", 0)
        .attr("y", (d) => yBand(d.label) ?? 0)
        .attr("width", (d) => xLinear(d.count) + 3)
        .attr("height", yBand.bandwidth())
        .attr("rx", 3)
        .attr("fill", (d, i) => getColor(d.label, colorMap, i))
        .attr("opacity", opacityFn)
        .style("cursor", "pointer")
        .on("mouseover", handleMouseOver)
        .on("mousemove", handleMouseOver)
        .on("mouseout", handleMouseOut)
        .on("click", handleClick);

      g.selectAll(".label")
        .data(data)
        .join("text")
        .attr("class", "label")
        .attr("x", (d) => xLinear(d.count) + 6)
        .attr("y", (d) => (yBand(d.label) ?? 0) + yBand.bandwidth() / 2)
        .attr("dominant-baseline", "central")
        .attr("text-anchor", "start")
        .style("font-size", "9px")
        .style("fill", "var(--fg-secondary, #6b7280)")
        .attr("opacity", opacityFn)
        .text((d) => d.count);
    } else {
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

      const xAxis = g.append("g").attr("transform", `translate(0,${h})`).call(d3.axisBottom(x));
      xAxis
        .selectAll("text")
        .attr("transform", "rotate(-35)")
        .style("text-anchor", "end")
        .style("font-size", "9px")
        .style("fill", axisColor);
      xAxis.selectAll("line").style("stroke", axisColor);
      xAxis.select(".domain").style("stroke", axisColor);

      const yAxis = g.append("g").call(d3.axisLeft(y).ticks(5));
      yAxis.selectAll("text").style("font-size", "9px").style("fill", axisColor);
      yAxis.selectAll("line").style("stroke", axisColor);
      yAxis.select(".domain").style("stroke", axisColor);

      g.append("text")
        .attr("transform", "rotate(-90)")
        .attr("x", -h / 2)
        .attr("y", -margin.left + 14)
        .attr("text-anchor", "middle")
        .style("font-size", "10px")
        .style("fill", "var(--fg-tertiary, #6b7280)")
        .text(yAxisLabel);

      g.append("clipPath")
        .attr("id", "bar-clip")
        .append("rect")
        .attr("x", 0)
        .attr("y", 0)
        .attr("width", w)
        .attr("height", h);

      const barGroup = g.append("g").attr("clip-path", "url(#bar-clip)");

      barGroup
        .selectAll(".bar")
        .data(data)
        .join("rect")
        .attr("class", "bar")
        .attr("x", (d) => x(d.label) ?? 0)
        .attr("y", (d) => y(d.count))
        .attr("width", x.bandwidth())
        .attr("height", (d) => h - y(d.count) + 3)
        .attr("rx", 3)
        .attr("fill", (d, i) => getColor(d.label, colorMap, i))
        .attr("opacity", opacityFn)
        .style("cursor", "pointer")
        .on("mouseover", handleMouseOver)
        .on("mousemove", handleMouseOver)
        .on("mouseout", handleMouseOut)
        .on("click", handleClick);

      g.selectAll(".label")
        .data(data)
        .join("text")
        .attr("class", "label")
        .attr("x", (d) => (x(d.label) ?? 0) + x.bandwidth() / 2)
        .attr("y", (d) => y(d.count) - 4)
        .attr("text-anchor", "middle")
        .style("font-size", "9px")
        .style("fill", "var(--fg-secondary, #6b7280)")
        .attr("opacity", opacityFn)
        .text((d) => d.count);
    }
  }, [data, width, height, resolvedHeight, yAxisLabel, selectedLabel, resolvedDirection, colorMap]);

  if (data.length === 0) return null;

  return (
    <div ref={containerRef} className="relative w-full">
      <svg ref={ref} className="w-full" />
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
      {legend && legend.length > 0 && (
        <div className="mt-3 flex flex-wrap gap-x-4 gap-y-1 border-t border-border-secondary pt-3">
          {legend.map((item) => (
            <span key={item.short} className="text-[11px] text-fg-tertiary">
              <span className="font-medium text-fg-secondary">{item.short}</span> — {item.full}
            </span>
          ))}
        </div>
      )}
    </div>
  );
}
