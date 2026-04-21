// packages/observer-web/src/components/charts/use-chart-tooltip.ts
import { type RefObject, useState } from "react";

interface TooltipState {
  visible: boolean;
  x: number;
  y: number;
  text: string;
}

export function useChartTooltip(containerRef: RefObject<HTMLDivElement | null>) {
  const [tooltip, setTooltip] = useState<TooltipState>({
    visible: false,
    x: 0,
    y: 0,
    text: "",
  });

  function show(event: MouseEvent, text: string) {
    const bounds = containerRef.current?.getBoundingClientRect();
    if (!bounds) return;
    setTooltip({
      visible: true,
      x: event.clientX - bounds.left + 12,
      y: event.clientY - bounds.top - 10,
      text,
    });
  }

  function hide() {
    setTooltip((prev) => ({ ...prev, visible: false }));
  }

  return { tooltip, show, hide };
}
