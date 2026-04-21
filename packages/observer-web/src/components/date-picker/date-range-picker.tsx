import { useRef, useState } from "react";

import { Popover } from "@base-ui/react/popover";
import { DayPicker, type DateRange } from "react-day-picker";

import { CalendarBlankIcon, CaretLeftIcon, CaretRightIcon } from "@/components/icons";

import { formatDisplay, toISO, parseISO, triggerClass } from "./utils";
import "./date-picker.css";

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
              captionLayout="dropdown"
              startMonth={new Date(2000, 0)}
              endMonth={new Date(new Date().getFullYear() + 1, 11)}
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
