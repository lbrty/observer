import { UISwitch } from "@/components/ui/ui-switch";

interface PermissionToggleRowProps {
  label: string;
  description: string;
  checked: boolean;
  onCheckedChange: (v: boolean) => void;
}

export function PermissionToggleRow({
  label,
  description,
  checked,
  onCheckedChange,
}: PermissionToggleRowProps) {
  return (
    <div className="space-y-2">
      <UISwitch checked={checked} onCheckedChange={onCheckedChange} label={label} />
      <p className="ml-11.5 text-xs text-fg-tertiary">{description}</p>
    </div>
  );
}
