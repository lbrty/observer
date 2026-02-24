const variants = {
  foam: "bg-foam/15 text-foam",
  gold: "bg-gold/15 text-gold",
  rose: "bg-rose/15 text-rose",
  neutral: "bg-bg-tertiary text-fg-secondary",
} as const;

const dotColors = {
  foam: "bg-foam",
  gold: "bg-gold",
  rose: "bg-rose",
  neutral: "bg-fg-tertiary",
} as const;

type Variant = keyof typeof variants;

interface StatusBadgeProps {
  label: string;
  variant?: Variant;
  dot?: boolean;
}

const roleVariants: Record<string, Variant> = {
  admin: "rose",
  staff: "gold",
  consultant: "neutral",
  guest: "neutral",
  owner: "rose",
  manager: "gold",
  viewer: "neutral",
};

const statusVariants: Record<string, Variant> = {
  active: "foam",
  archived: "neutral",
  closed: "rose",
  true: "foam",
  false: "neutral",
};

export function StatusBadge({ label, variant, dot }: StatusBadgeProps) {
  const resolved =
    variant ?? roleVariants[label] ?? statusVariants[label] ?? "neutral";

  return (
    <span
      className={`inline-flex items-center gap-1.5 rounded-full px-2 py-0.5 text-xs font-medium ${variants[resolved]}`}
    >
      {dot !== false && (
        <span className={`size-1.5 rounded-full ${dotColors[resolved]}`} />
      )}
      {label}
    </span>
  );
}

export function StatusDot({ active }: { active: boolean }) {
  return (
    <span
      className={`inline-block size-2 rounded-full ${active ? "bg-foam" : "bg-fg-tertiary/40"}`}
    />
  );
}
