import type { ReactNode } from "react";

interface PeopleChartSectionProps {
  title: string;
  filter?: ReactNode;
  children: ReactNode;
}

export function PeopleChartSection({ title, filter, children }: PeopleChartSectionProps) {
  return (
    <div className="rounded-xl border border-border-secondary bg-bg-secondary p-5">
      <div className="mb-4 flex items-center justify-between">
        <h3 className="text-sm font-semibold text-fg">{title}</h3>
        {filter}
      </div>
      {children}
    </div>
  );
}
