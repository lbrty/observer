import type { ReactNode } from "react";

import { SectionHeading } from "@/components/section-heading";

interface FormSectionProps {
  title: string;
  columns?: 1 | 2 | 3;
  children: ReactNode;
  className?: string;
}

const gridCols = {
  1: "grid-cols-1",
  2: "grid-cols-1 sm:grid-cols-2",
  3: "grid-cols-1 sm:grid-cols-3",
};

export function FormSection({ title, columns = 2, children, className }: FormSectionProps) {
  return (
    <>
      <SectionHeading>{title}</SectionHeading>
      <div className={`grid gap-4 ${gridCols[columns]} ${className ?? ""}`}>
        {children}
      </div>
    </>
  );
}
