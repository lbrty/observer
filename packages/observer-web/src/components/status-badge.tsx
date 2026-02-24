const variants = {
  foam: "bg-foam/15 text-foam",
  gold: "bg-gold/15 text-gold",
  rose: "bg-rose/15 text-rose",
  neutral: "bg-bg-tertiary text-fg-secondary",
} as const;

type Variant = keyof typeof variants;

interface StatusBadgeProps {
  label: string;
  variant?: Variant;
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

export function StatusBadge({ label, variant }: StatusBadgeProps) {
  const resolved =
    variant ?? roleVariants[label] ?? statusVariants[label] ?? "neutral";

  return (
    <span
      className={`inline-block rounded-full px-2 py-0.5 text-xs font-medium ${variants[resolved]}`}
    >
      {label}
    </span>
  );
}
