import type { Icon } from "@/components/icons";
import type { ReactNode } from "react";

interface EmptyStateProps {
  icon: Icon;
  title: string;
  description?: string;
  action?: ReactNode;
}

export function EmptyState({ icon: IconComp, title, description, action }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center py-16 text-center">
      <span className="mb-4 inline-flex size-12 items-center justify-center rounded-xl bg-bg-tertiary text-fg-tertiary">
        <IconComp size={24} />
      </span>
      <h3 className="text-sm font-medium text-fg">{title}</h3>
      {description && <p className="mt-1 max-w-xs text-sm text-fg-tertiary">{description}</p>}
      {action && <div className="mt-4">{action}</div>}
    </div>
  );
}
