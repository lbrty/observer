export interface BarLegendItem {
  short: string;
  full: string;
}

interface ChartLegendProps {
  legend: BarLegendItem[];
}

export function ChartLegend({ legend }: ChartLegendProps) {
  if (legend.length === 0) return null;
  return (
    <div className="mt-3 flex flex-wrap gap-x-4 gap-y-1 border-t border-border-secondary pt-3">
      {legend.map((item) => (
        <span key={item.short} className="text-[11px] text-fg-tertiary">
          <span className="font-medium text-fg-secondary">{item.short}</span> — {item.full}
        </span>
      ))}
    </div>
  );
}
