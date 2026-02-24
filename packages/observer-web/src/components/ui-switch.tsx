import { Switch } from "@base-ui/react/switch";

interface UISwitchProps {
  checked: boolean;
  onCheckedChange: (checked: boolean) => void;
  label: string;
  disabled?: boolean;
}

export function UISwitch({
  checked,
  onCheckedChange,
  label,
  disabled,
}: UISwitchProps) {
  return (
    <label className="inline-flex cursor-pointer items-center gap-2.5 text-sm text-fg-secondary select-none">
      <Switch.Root
        checked={checked}
        onCheckedChange={onCheckedChange}
        disabled={disabled}
        className="relative inline-flex h-5 w-9 shrink-0 cursor-pointer items-center rounded-full border border-transparent bg-bg-tertiary transition-colors data-[checked]:bg-accent data-[disabled]:cursor-not-allowed data-[disabled]:opacity-50"
      >
        <Switch.Thumb className="pointer-events-none block size-3.5 translate-x-0.5 rounded-full bg-fg-tertiary shadow-sm transition-transform data-[checked]:translate-x-[18px] data-[checked]:bg-accent-fg" />
      </Switch.Root>
      {label}
    </label>
  );
}
