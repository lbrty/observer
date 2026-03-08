export function KpiCard({ label, value }: { label: string; value: number }) {
  return (
    <div className="rounded-xl border border-border-secondary bg-bg-secondary p-4">
      <p className="text-2xl font-bold tabular-nums text-fg">{value.toLocaleString()}</p>
      <p className="mt-0.5 text-xs font-medium text-fg-tertiary">{label}</p>
    </div>
  );
}
