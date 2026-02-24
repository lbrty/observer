import { CaretUpDownIcon, CheckIcon } from "@/components/icons";
import { Select } from "@base-ui/react/select";

interface SelectOption {
  label: string;
  value: string;
}

interface UISelectProps {
  value: string;
  onValueChange: (value: string) => void;
  options: SelectOption[];
  placeholder?: string;
  disabled?: boolean;
  fullWidth?: boolean;
}

export function UISelect({
  value,
  onValueChange,
  options,
  placeholder,
  disabled,
  fullWidth,
}: UISelectProps) {
  return (
    <Select.Root
      value={value}
      onValueChange={(v) => {
        if (v !== null) onValueChange(v);
      }}
      disabled={disabled}
    >
      <Select.Trigger
        className={`inline-flex h-9 items-center justify-between gap-2 rounded-lg border border-border-secondary bg-bg-secondary pr-2 pl-3 text-sm text-fg outline-none hover:border-border-primary focus:border-accent disabled:opacity-50 data-popup-open:border-accent ${fullWidth ? "w-full" : ""}`}
      >
        <Select.Value placeholder={placeholder}>
          {value
            ? (options.find((o) => o.value === value)?.label ?? value)
            : null}
        </Select.Value>
        <Select.Icon className="text-fg-tertiary">
          <CaretUpDownIcon size={14} weight="bold" />
        </Select.Icon>
      </Select.Trigger>

      <Select.Portal>
        <Select.Positioner sideOffset={4} align="start" className="z-[60]">
          <Select.Popup className="origin-(--transform-origin) rounded-lg border border-border-secondary bg-bg-secondary py-1 shadow-elevated transition-[transform,scale,opacity] data-ending-style:scale-95 data-ending-style:opacity-0 data-starting-style:scale-95 data-starting-style:opacity-0">
            <Select.List>
              {options.map((opt) => (
                <Select.Item
                  key={opt.value}
                  value={opt.value}
                  className="flex cursor-pointer items-center gap-2 px-3 py-1.5 text-sm text-fg outline-none select-none data-highlighted:bg-bg-tertiary"
                >
                  <Select.ItemIndicator className="inline-flex w-4 items-center justify-center text-accent">
                    <CheckIcon size={14} weight="bold" />
                  </Select.ItemIndicator>
                  <Select.ItemText>{opt.label}</Select.ItemText>
                </Select.Item>
              ))}
            </Select.List>
          </Select.Popup>
        </Select.Positioner>
      </Select.Portal>
    </Select.Root>
  );
}
