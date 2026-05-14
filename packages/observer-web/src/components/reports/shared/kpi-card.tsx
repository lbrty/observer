import { cardGradient } from "@/lib/card-gradient";

export function KpiCard({ label, value, index }: { label: string; value: number; index?: number }) {
  const gradClass = index !== undefined ? cardGradient(index) : "";
  return (
    <div className={`rounded-xl border border-border-secondary bg-bg-secondary p-4 ${gradClass}`}>
      <p className="relative text-2xl font-bold tabular-nums text-fg">{value.toLocaleString()}</p>
      <p className="relative mt-0.5 text-xs font-medium text-fg-tertiary">{label}</p>
    </div>
  );
}
