interface PersonDetailProps {
  label: string;
  value?: string | null;
}

export function PersonDetail({ label, value }: PersonDetailProps) {
  if (!value) return null;
  return (
    <div>
      <dt className="text-xs font-medium text-fg-tertiary">{label}</dt>
      <dd className="mt-0.5 text-sm text-fg">{value}</dd>
    </div>
  );
}
