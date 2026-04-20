// packages/observer-web/src/components/charts/chart-tooltip.tsx
import type { ReactNode } from "react";

interface ChartTooltipProps {
  visible: boolean;
  x: number;
  y: number;
  children: ReactNode;
}

export function ChartTooltip({ visible, x, y, children }: ChartTooltipProps) {
  if (!visible) return null;
  return (
    <div
      style={{
        position: "absolute",
        left: x,
        top: y,
        background: "var(--bg-secondary)",
        border: "1px solid var(--border-secondary)",
        borderRadius: 8,
        padding: "6px 10px",
        fontSize: 12,
        color: "var(--fg)",
        boxShadow: "var(--shadow-elevated)",
        pointerEvents: "none",
        whiteSpace: "nowrap",
      }}
    >
      {children}
    </div>
  );
}
