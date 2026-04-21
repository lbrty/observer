import * as d3 from "d3";

import type { CountResult } from "@/types/report";

import { getColor } from "../colors";
import type { RenderOptions } from "./render-horizontal";

export function renderVerticalBars(
  svg: d3.Selection<SVGSVGElement, unknown, null, undefined>,
  data: CountResult[],
  opts: RenderOptions,
) {
  const { width, height, yAxisLabel, selectedLabel, colorMap, clipId, onMouseOver, onMouseOut, onClick } = opts;

  const axisColor = "var(--fg-tertiary, #6b7280)";
  const margin = { top: 20, right: 20, bottom: 60, left: 50 };
  const w = width - margin.left - margin.right;
  const h = height - margin.top - margin.bottom;

  const opacityFn = (d: CountResult) =>
    selectedLabel === null || selectedLabel === d.label ? 1 : 0.3;

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
    .attr("id", clipId)
    .append("rect")
    .attr("x", 0)
    .attr("y", 0)
    .attr("width", w)
    .attr("height", h);

  const barGroup = g.append("g").attr("clip-path", `url(#${clipId})`);

  const rV = 3;
  barGroup
    .selectAll(".bar")
    .data(data)
    .join("path")
    .attr("class", "bar")
    .attr("d", (d) => {
      const bx = x(d.label) ?? 0;
      const by = y(d.count);
      const bw = x.bandwidth();
      const bh = h - y(d.count);
      return `M${bx},${by + rV} a${rV},${rV} 0 0 1 ${rV},${-rV} H${bx + bw - rV} a${rV},${rV} 0 0 1 ${rV},${rV} V${by + bh} H${bx} Z`;
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
    .attr("x", (d) => (x(d.label) ?? 0) + x.bandwidth() / 2)
    .attr("y", (d) => y(d.count) - 4)
    .attr("text-anchor", "middle")
    .style("font-size", "9px")
    .style("fill", "var(--fg-secondary, #6b7280)")
    .attr("opacity", opacityFn)
    .text((d) => d.count);
}
