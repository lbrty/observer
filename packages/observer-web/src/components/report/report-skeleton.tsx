export function ReportSkeleton({ kpiCount = 6, chartCount = 4 }: { kpiCount?: number; chartCount?: number }) {
  return (
    <div className="space-y-6">
      <div className={`grid gap-4 ${kpiCount > 4 ? "grid-cols-3 lg:grid-cols-6" : "grid-cols-2 sm:grid-cols-4"}`}>
        {Array.from({ length: kpiCount }).map((_, i) => (
          <div key={i} className="h-20 animate-pulse rounded-xl bg-bg-tertiary" />
        ))}
      </div>
      <div className="grid gap-6 lg:grid-cols-2">
        {Array.from({ length: chartCount }).map((_, i) => (
          <div key={i} className="h-72 animate-pulse rounded-xl bg-bg-tertiary" />
        ))}
      </div>
    </div>
  );
}
