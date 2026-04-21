import { useEffect, useId, useRef, useState } from "react";

import * as d3 from "d3";

import type { CountResult } from "@/types/report";

import { ChartTooltip } from "../chart-tooltip";
import type { BarLegendItem } from "./chart-legend";
import { ChartLegend } from "./chart-legend";
import { renderHorizontalBars } from "./render-horizontal";
import { renderVerticalBars } from "./render-vertical";

export type { BarLegendItem };

interface Tooltip {
  visible: boolean;
  x: number;
  y: number;
  label: string;
  count: number;
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
  width = 700,
  height = 260,
  yAxisLabel = "Count",
  legend,
  direction,
  colorMap,
}: BarChartProps) {
  const ref = useRef<SVGSVGElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const clipId = useId().replace(/:/g, "_");
  const [tooltip, setTooltip] = useState<Tooltip>({
    visible: false,
    x: 0,
    y: 0,
    label: "",
    count: 0,
  });
  const [selectedLabel, setSelectedLabel] = useState<string | null>(null);

  const resolvedDirection =
    direction === "auto"
      ? data.length > 6
        ? "horizontal"
        : "vertical"
      : (direction ?? "vertical");

  const resolvedHeight =
    resolvedDirection === "horizontal" ? Math.max(200, data.length * 24) : height;

  useEffect(() => {
    if (!ref.current || data.length === 0) return;

    const svg = d3.select(ref.current);
    svg.selectAll("*").remove();

    const handleMouseOver = (event: MouseEvent, d: CountResult) => {
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

    const handleMouseOut = () => {
      setTooltip((prev) => ({ ...prev, visible: false }));
    };

    const handleClick = (_event: MouseEvent, d: CountResult) => {
      setSelectedLabel((prev) => (prev === d.label ? null : d.label));
    };

    const opts = {
      width,
      height,
      resolvedHeight,
      yAxisLabel,
      selectedLabel,
      colorMap,
      clipId,
      onMouseOver: handleMouseOver,
      onMouseOut: handleMouseOut,
      onClick: handleClick,
    };

    if (resolvedDirection === "horizontal") {
      renderHorizontalBars(svg, data, opts);
    } else {
      renderVerticalBars(svg, data, opts);
    }
  }, [
    data,
    width,
    height,
    resolvedHeight,
    yAxisLabel,
    selectedLabel,
    resolvedDirection,
    colorMap,
    clipId,
  ]);

  if (data.length === 0) return null;

  return (
    <div ref={containerRef} className="relative w-full">
      <svg ref={ref} className="w-full" style={{ maxHeight: 360 }} />
      <ChartTooltip visible={tooltip.visible} x={tooltip.x} y={tooltip.y}>
        <strong>{tooltip.label}</strong>: {tooltip.count}
      </ChartTooltip>
      {legend && <ChartLegend legend={legend} />}
    </div>
  );
}
