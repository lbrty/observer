import { Popover } from "@base-ui/react/popover";
import { useRef, useState } from "react";
import { DayPicker, type DateRange } from "react-day-picker";

import { CalendarBlankIcon, CaretLeftIcon, CaretRightIcon, XIcon } from "@/components/icons";

import "./date-picker.css";

function formatDisplay(iso: string): string {
  if (!iso) return "";
  const [y, m, d] = iso.split("-");
  return `${d}.${m}.${y}`;
}

function toISO(date: Date): string {
  const y = date.getFullYear();
  const m = String(date.getMonth() + 1).padStart(2, "0");
  const d = String(date.getDate()).padStart(2, "0");
  return `${y}-${m}-${d}`;
}

function parseISO(iso: string): Date | undefined {
  if (!iso) return undefined;
  const [y, m, d] = iso.split("-").map(Number);
  return new Date(y, m - 1, d);
}

const triggerClass =
  "flex h-9 w-full items-center gap-2 rounded-lg border border-border-secondary bg-bg-secondary px-3 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg disabled:opacity-50";

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
  captionLayout,
  startYear,
  endYear,
  className,
}: DatePickerProps) {
  const [open, setOpen] = useState(false);
  const selected = parseISO(value);

  const now = new Date();
  const useDropdown = captionLayout === "dropdown" || captionLayout === "dropdown-months" || captionLayout === "dropdown-years";
  const startMonth = useDropdown ? new Date(startYear ?? 1920, 0) : undefined;
  const endMonth = useDropdown ? new Date(endYear ?? now.getFullYear(), 11) : undefined;

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

interface DateRangePickerProps {
  from: string;
  to: string;
  onChange: (range: { from?: string; to?: string }) => void;
  placeholderFrom?: string;
  placeholderTo?: string;
  disabled?: boolean;
  className?: string;
}

export function DateRangePicker({
  from,
  to,
  onChange,
  placeholderFrom = "dd.mm.yyyy",
  placeholderTo = "dd.mm.yyyy",
  disabled,
  className,
}: DateRangePickerProps) {
  const [open, setOpen] = useState(false);
  const selected: DateRange = {
    from: parseISO(from),
    to: parseISO(to),
  };
  const clickCount = useRef(0);

  return (
    <Popover.Root open={open} onOpenChange={setOpen}>
      <Popover.Trigger
        disabled={disabled}
        className={`inline-flex items-center gap-2 ${className ?? ""}`}
      >
        <span className={triggerClass}>
          <CalendarBlankIcon className="size-4 shrink-0 text-fg-tertiary" />
          <span className={from ? "text-fg" : "text-fg-tertiary"}>
            {from ? formatDisplay(from) : placeholderFrom}
          </span>
        </span>
        <span className="text-fg-tertiary">&ndash;</span>
        <span className={triggerClass}>
          <CalendarBlankIcon className="size-4 shrink-0 text-fg-tertiary" />
          <span className={to ? "text-fg" : "text-fg-tertiary"}>
            {to ? formatDisplay(to) : placeholderTo}
          </span>
        </span>
      </Popover.Trigger>
      <Popover.Portal>
        <Popover.Positioner sideOffset={4} className="z-60">
          <Popover.Popup className="rdp-popup rounded-xl border border-border-secondary bg-bg-secondary p-3 shadow-elevated outline-none">
            <DayPicker
              mode="range"
              selected={selected}
              onSelect={(range) => {
                clickCount.current += 1;
                onChange({
                  from: range?.from ? toISO(range.from) : undefined,
                  to: range?.to ? toISO(range.to) : undefined,
                });
                if (clickCount.current >= 2) {
                  clickCount.current = 0;
                  setOpen(false);
                }
              }}
              defaultMonth={parseISO(from)}
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
  );
}
