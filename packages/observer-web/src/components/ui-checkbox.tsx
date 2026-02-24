import { Check } from "@phosphor-icons/react";
import { Checkbox } from "@base-ui/react/checkbox";

interface UICheckboxProps {
  checked: boolean;
  onCheckedChange: (checked: boolean) => void;
  label: string;
  disabled?: boolean;
}

export function UICheckbox({
  checked,
  onCheckedChange,
  label,
  disabled,
}: UICheckboxProps) {
  return (
    <label className="inline-flex cursor-pointer items-center gap-2 text-sm text-fg-secondary select-none">
      <Checkbox.Root
        checked={checked}
        onCheckedChange={onCheckedChange}
        disabled={disabled}
        className="flex size-4.5 shrink-0 items-center justify-center rounded border border-border-secondary bg-bg-secondary transition-colors data-[checked]:border-accent data-[checked]:bg-accent data-[disabled]:cursor-not-allowed data-[disabled]:opacity-50"
      >
        <Checkbox.Indicator className="flex items-center justify-center text-accent-fg">
          <Check size={12} weight="bold" />
        </Checkbox.Indicator>
      </Checkbox.Root>
      {label}
    </label>
  );
}
