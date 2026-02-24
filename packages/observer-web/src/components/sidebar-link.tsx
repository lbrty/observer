import type { Icon } from "@/components/icons";
import { Link, useMatchRoute } from "@tanstack/react-router";

interface SidebarLinkProps {
  to: string;
  label: string;
  icon: Icon;
}

export function SidebarLink({ to, label, icon: Icon }: SidebarLinkProps) {
  const matchRoute = useMatchRoute();
  const isActive = !!matchRoute({ to, fuzzy: true });

  return (
    <Link
      to={to}
      className={`flex items-center gap-2.5 rounded-lg px-3 py-2 text-sm transition-colors ${
        isActive
          ? "bg-accent/10 font-medium text-accent"
          : "text-fg-secondary hover:bg-bg-tertiary hover:text-fg"
      }`}
    >
      <Icon size={18} weight={isActive ? "fill" : "regular"} />
      {label}
    </Link>
  );
}
