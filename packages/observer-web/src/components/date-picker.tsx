import { Popover } from "@base-ui/react/popover";
import { useRef, useState } from "react";
import { DayPicker, type DateRange } from "react-day-picker";

import { CalendarBlankIcon, CaretLeftIcon, CaretRightIcon } from "@/components/icons";

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
  "flex h-9 w-full items-center gap-2 rounded-lg border border-border-secondary bg-bg-secondary px-3 text-sm text-fg outline-none focus:border-accent disabled:opacity-50";

interface DatePickerProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  disabled?: boolean;
  className?: string;
}

export function DatePicker({
  value,
  onChange,
  placeholder = "dd.mm.yyyy",
  disabled,
  className,
}: DatePickerProps) {
  const [open, setOpen] = useState(false);
  const selected = parseISO(value);

  return (
    <Popover.Root open={open} onOpenChange={setOpen}>
      <Popover.Trigger disabled={disabled} className={`${triggerClass} ${className ?? ""}`}>
        <CalendarBlankIcon className="size-4 shrink-0 text-fg-tertiary" />
        <span className={value ? "text-fg" : "text-fg-tertiary"}>
          {value ? formatDisplay(value) : placeholder}
        </span>
      </Popover.Trigger>
      <Popover.Portal>
        <Popover.Positioner sideOffset={4}>
          <Popover.Popup className="rdp-popup rounded-xl border border-border-secondary bg-bg-secondary p-3 shadow-elevated outline-none">
            <DayPicker
              mode="single"
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
        <Popover.Positioner sideOffset={4}>
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
