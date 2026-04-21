import type { ReactNode } from "react";

export function FilterField({ label, children }: { label: string; children: ReactNode }) {
  return (
    <div className="min-w-48 space-y-1.5">
      <span className="block text-xs font-medium text-fg-secondary">{label}</span>
      {children}
    </div>
  );
}
