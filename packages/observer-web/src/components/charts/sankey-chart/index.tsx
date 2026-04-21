import { useEffect, useRef } from "react";

import * as d3 from "d3";
import { sankey, sankeyLinkHorizontal } from "d3-sankey";

import type { StatusFlow } from "@/types/report";

import { ChartTooltip } from "../chart-tooltip";
import { useChartTooltip } from "../use-chart-tooltip";

const NODE_COLORS: Record<string, string> = {
  new: "#5b8af8",
  active: "#30a46c",
  closed: "#f59e0b",
  archived: "#8b909e",
};

function nodeColor(name: string): string {
  return NODE_COLORS[name] ?? "#8b5cf6";
}

interface SankeyChartProps {
  data: StatusFlow[];
  width?: number;
  height?: number;
  translateLabel?: (label: string) => string;
}

export function SankeyChart({ data, width = 500, height = 260, translateLabel }: SankeyChartProps) {
  const svgRef = useRef<SVGSVGElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const { tooltip, show, hide } = useChartTooltip(containerRef);

  useEffect(() => {
    if (!svgRef.current || data.length === 0) return;

    const svg = d3.select(svgRef.current);
    svg.selectAll("*").remove();

    const margin = { top: 8, right: 8, bottom: 8, left: 8 };
    const w = width - margin.left - margin.right;
    const h = height - margin.top - margin.bottom;

    // Define forward order to prevent circular links
    const STATUS_ORDER: Record<string, number> = {
      new: 0,
      active: 1,
      closed: 2,
      archived: 3,
    };

    // Only keep forward transitions (from lower to higher order)
    const forwardData = data.filter((d) => {
      const fromIdx = STATUS_ORDER[d.from_status] ?? -1;
      const toIdx = STATUS_ORDER[d.to_status] ?? -1;
      return fromIdx < toIdx;
    });

    const nodeNames = Array.from(new Set(forwardData.flatMap((d) => [d.from_status, d.to_status])));

    if (nodeNames.length === 0) return;

    const nodes = nodeNames.map((name) => ({ name }));
    const links = forwardData.map((d) => ({
      source: d.from_status,
      target: d.to_status,
      value: d.count,
      avgDays: d.avg_days,
      fromStatus: d.from_status,
      toStatus: d.to_status,
    }));

    type N = { name: string };
    type L = (typeof links)[0];

    const layout = sankey<N, L>()
      .nodeId((d) => d.name)
      .nodeWidth(14)
      .nodePadding(16)
      .extent([
        [0, 0],
        [w, h],
      ]);

    const { nodes: sNodes, links: sLinks } = layout({
      nodes: nodes.map((d) => ({ ...d })),
      links: links.map((d) => ({ ...d })),
    });

    const g = svg
      .attr("viewBox", `0 0 ${width} ${height}`)
      .append("g")
      .attr("transform", `translate(${margin.left},${margin.top})`);

    g.selectAll(".link")
      .data(sLinks)
      .join("path")
      .attr("class", "link")
      .attr("d", sankeyLinkHorizontal())
      .attr("fill", "none")
      .attr("stroke", (d) => {
        const src = d.source as unknown as N;
        return nodeColor(src.name);
      })
      .attr("stroke-opacity", 0.3)
      .attr("stroke-width", (d) => Math.max(1, d.width ?? 0))
      .style("cursor", "pointer")
      .on("mouseover", function (event: MouseEvent, d) {
        d3.select(this).attr("stroke-opacity", 0.6);
        const tl = translateLabel ?? ((l: string) => l);
        const fromStatus = (d as unknown as { fromStatus: string }).fromStatus;
        const toStatus = (d as unknown as { toStatus: string }).toStatus;
        const value = (d as unknown as { value: number }).value;
        const avgDays = (d as unknown as { avgDays: number }).avgDays;
        show(event, `${tl(fromStatus)} → ${tl(toStatus)}: ${value} (~${avgDays}d)`);
      })
      .on("mousemove", function (event: MouseEvent, d) {
        const tl = translateLabel ?? ((l: string) => l);
        const fromStatus = (d as unknown as { fromStatus: string }).fromStatus;
        const toStatus = (d as unknown as { toStatus: string }).toStatus;
        const value = (d as unknown as { value: number }).value;
        const avgDays = (d as unknown as { avgDays: number }).avgDays;
        show(event, `${tl(fromStatus)} → ${tl(toStatus)}: ${value} (~${avgDays}d)`);
      })
      .on("mouseout", function () {
        d3.select(this).attr("stroke-opacity", 0.3);
        hide();
      });

    g.selectAll(".node")
      .data(sNodes)
      .join("rect")
      .attr("class", "node")
      .attr("x", (d) => d.x0 ?? 0)
      .attr("y", (d) => d.y0 ?? 0)
      .attr("width", (d) => (d.x1 ?? 0) - (d.x0 ?? 0))
      .attr("height", (d) => Math.max(1, (d.y1 ?? 0) - (d.y0 ?? 0)))
      .attr("rx", 2)
      .attr("fill", (d) => nodeColor(d.name));

    const tl = translateLabel ?? ((l: string) => l);
    g.selectAll(".node-label")
      .data(sNodes)
      .join("text")
      .attr("class", "node-label")
      .attr("x", (d) => {
        const x0 = d.x0 ?? 0;
        const x1 = d.x1 ?? 0;
        return x0 < w / 2 ? x1 + 6 : x0 - 6;
      })
      .attr("y", (d) => ((d.y0 ?? 0) + (d.y1 ?? 0)) / 2)
      .attr("dy", "0.35em")
      .attr("text-anchor", (d) => ((d.x0 ?? 0) < w / 2 ? "start" : "end"))
      .style("font-size", "9px")
      .style("font-weight", "600")
      .style("fill", "var(--fg-secondary, #6b7280)")
      .text((d) => tl(d.name));
  }, [data, width, height, translateLabel]);

  if (data.length === 0) return null;

  return (
    <div ref={containerRef} className="relative w-full">
      <svg ref={svgRef} className="w-full" />
      <ChartTooltip visible={tooltip.visible} x={tooltip.x} y={tooltip.y}>
        {tooltip.text}
      </ChartTooltip>
    </div>
  );
}
