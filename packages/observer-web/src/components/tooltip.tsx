import { Tooltip as BaseTooltip } from "@base-ui/react/tooltip";
import type { ReactElement } from "react";

interface TooltipProps {
  label: string;
  children: ReactElement;
  side?: "top" | "bottom" | "left" | "right";
}

export function Tooltip({ label, children, side = "top" }: TooltipProps) {
  return (
    <BaseTooltip.Provider delay={200}>
      <BaseTooltip.Root>
        <BaseTooltip.Trigger render={children} />
        <BaseTooltip.Portal>
          <BaseTooltip.Positioner side={side} sideOffset={6}>
            <BaseTooltip.Popup className="rounded-md bg-fg px-2.5 py-1 text-xs font-medium text-bg shadow-elevated">
              {label}
            </BaseTooltip.Popup>
          </BaseTooltip.Positioner>
        </BaseTooltip.Portal>
      </BaseTooltip.Root>
    </BaseTooltip.Provider>
  );
}
