import { XIcon } from "@/components/icons";

export function FilterChip({ label, value, onRemove }: { label: string; value: string; onRemove: () => void }) {
  return (
    <button
      type="button"
      onClick={onRemove}
      className="inline-flex items-center gap-1 rounded-md bg-bg-tertiary px-2 py-0.5 text-xs font-medium text-fg-secondary transition-colors hover:text-fg"
    >
      <span className="text-fg-tertiary">{label}:</span> {value}
      <XIcon size={10} />
    </button>
  );
}
