import * as d3 from "d3";

import type { CountResult } from "@/types/report";

import { getColor } from "../colors";

export interface RenderOptions {
  width: number;
  height: number;
  resolvedHeight: number;
  yAxisLabel: string;
  selectedLabel: string | null;
  colorMap?: Record<string, string>;
  clipId: string;
  onMouseOver: (event: MouseEvent, d: CountResult) => void;
  onMouseOut: () => void;
  onClick: (event: MouseEvent, d: CountResult) => void;
}

export function renderHorizontalBars(
  svg: d3.Selection<SVGSVGElement, unknown, null, undefined>,
  data: CountResult[],
  opts: RenderOptions,
) {
  const { width, resolvedHeight, selectedLabel, colorMap, clipId, onMouseOver, onMouseOut, onClick } = opts;

  const axisColor = "var(--fg-tertiary, #6b7280)";
  const margin = { top: 10, right: 40, bottom: 20, left: 120 };
  const w = width - margin.left - margin.right;
  const h = resolvedHeight - margin.top - margin.bottom;

  const opacityFn = (d: CountResult) =>
    selectedLabel === null || selectedLabel === d.label ? 1 : 0.3;

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

  const xAxis = g
    .append("g")
    .attr("transform", `translate(0,${h})`)
    .call(d3.axisBottom(xLinear).ticks(5));
  xAxis.selectAll("text").style("font-size", "9px").style("fill", axisColor);
  xAxis.selectAll("line").style("stroke", axisColor);
  xAxis.select(".domain").style("stroke", axisColor);

  const yAxis = g.append("g").call(d3.axisLeft(yBand));
  yAxis.selectAll("text").style("font-size", "9px").style("fill", axisColor);
  yAxis.selectAll("line").style("stroke", axisColor);
  yAxis.select(".domain").style("stroke", axisColor);

  g.append("clipPath")
    .attr("id", clipId)
    .append("rect")
    .attr("x", 0)
    .attr("y", 0)
    .attr("width", w)
    .attr("height", h);

  const barGroup = g.append("g").attr("clip-path", `url(#${clipId})`);

  const rH = 3;
  barGroup
    .selectAll(".bar")
    .data(data)
    .join("path")
    .attr("class", "bar")
    .attr("d", (d) => {
      const bx = 0;
      const by = yBand(d.label) ?? 0;
      const bw = xLinear(d.count);
      const bh = yBand.bandwidth();
      return `M${bx},${by} H${bx + bw - rH} a${rH},${rH} 0 0 1 ${rH},${rH} V${by + bh - rH} a${rH},${rH} 0 0 1 ${-rH},${rH} H${bx} Z`;
    })
    .attr("fill", (d, i) => getColor(d.label, colorMap, i))
    .attr("opacity", opacityFn)
    .style("cursor", "pointer")
    .on("mouseover", onMouseOver)
    .on("mousemove", onMouseOver)
    .on("mouseout", onMouseOut)
    .on("click", onClick);

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
}
