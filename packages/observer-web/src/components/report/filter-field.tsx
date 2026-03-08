import type { ReactNode } from "react";

export function FilterField({ label, children }: { label: string; children: ReactNode }) {
  return (
    <div className="space-y-1.5">
      <span className="block text-xs font-medium text-fg-secondary">{label}</span>
      {children}
    </div>
  );
}
