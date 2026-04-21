import { useState } from "react";

import { Popover } from "@base-ui/react/popover";
import { DayPicker } from "react-day-picker";

import { CalendarBlankIcon, CaretLeftIcon, CaretRightIcon, XIcon } from "@/components/ui/icons";

import { formatDisplay, toISO, parseISO, triggerClass } from "./utils";
import "./date-picker.css";

interface DatePickerProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  disabled?: boolean;
  clearable?: boolean;
  captionLayout?: "label" | "dropdown" | "dropdown-months" | "dropdown-years";
  startYear?: number;
  endYear?: number;
  className?: string;
}

export function DatePicker({
  value,
  onChange,
  placeholder = "dd.mm.yyyy",
  disabled,
  clearable,
  captionLayout = "dropdown",
  startYear,
  endYear,
  className,
}: DatePickerProps) {
  const [open, setOpen] = useState(false);
  const selected = parseISO(value);

  const now = new Date();
  const useDropdown =
    captionLayout === "dropdown" ||
    captionLayout === "dropdown-months" ||
    captionLayout === "dropdown-years";
  const startMonth = useDropdown ? new Date(startYear ?? 1920, 0) : undefined;
  const endMonth = useDropdown ? new Date(endYear ?? now.getFullYear() + 1, 11) : undefined;

  return (
    <div className="flex items-center gap-1">
      <Popover.Root open={open} onOpenChange={setOpen}>
        <Popover.Trigger
          disabled={disabled}
          className={`flex-1 ${triggerClass} ${className ?? ""}`}
        >
          <CalendarBlankIcon className="size-4 shrink-0 text-fg-tertiary" />
          <span className={value ? "text-fg" : "text-fg-tertiary"}>
            {value ? formatDisplay(value) : placeholder}
          </span>
        </Popover.Trigger>
        <Popover.Portal>
          <Popover.Positioner sideOffset={4} className="z-60">
            <Popover.Popup className="rdp-popup rounded-xl border border-border-secondary bg-bg-secondary p-3 shadow-elevated outline-none">
              <DayPicker
                mode="single"
                captionLayout={captionLayout}
                selected={selected}
                onSelect={(day) => {
                  if (day) {
                    onChange(toISO(day));
                  } else {
                    onChange("");
                  }
                  setOpen(false);
                }}
                defaultMonth={selected}
                startMonth={startMonth}
                endMonth={endMonth}
                components={{
                  Chevron: ({ orientation }) =>
                    orientation === "left" ? (
                      <CaretLeftIcon className="size-4" />
                    ) : (
                      <CaretRightIcon className="size-4" />
                    ),
                }}
              />
            </Popover.Popup>
          </Popover.Positioner>
        </Popover.Portal>
      </Popover.Root>
      {clearable && value && (
        <button
          type="button"
          onClick={() => onChange("")}
          className="inline-flex size-9 shrink-0 cursor-pointer items-center justify-center rounded-lg text-fg-tertiary hover:bg-bg-tertiary hover:text-fg"
        >
          <XIcon className="size-4" />
        </button>
      )}
    </div>
  );
}
