import { DrawerPreview as Drawer } from "@base-ui/react/drawer";
import type { FormEvent, ReactNode } from "react";
import { useTranslation } from "react-i18next";

import { XIcon } from "@/components/icons";

interface DrawerShellProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  onSubmit: (e: FormEvent) => void;
  isPending?: boolean;
  submitLabel?: string;
  savingLabel?: string;
  children: ReactNode;
  footer?: ReactNode;
}

export function DrawerShell({
  open,
  onOpenChange,
  title,
  onSubmit,
  isPending,
  submitLabel,
  savingLabel,
  children,
  footer,
}: DrawerShellProps) {
  const { t } = useTranslation();
  const saveText = submitLabel ?? t("admin.common.save");
  const savingText = savingLabel ?? t("admin.common.saving");

  return (
    <Drawer.Root open={open} onOpenChange={onOpenChange} swipeDirection="right">
      <Drawer.Portal>
        <Drawer.Backdrop className="fixed inset-0 z-50 bg-black/25 backdrop-blur-xs transition-opacity duration-200 data-ending-style:opacity-0 data-starting-style:opacity-0" />
        <Drawer.Viewport className="fixed inset-0 z-50">
          <Drawer.Popup className="fixed top-0 right-0 flex h-dvh w-full max-w-[840px] flex-col border-l border-border-secondary bg-bg-secondary shadow-elevated transition-transform duration-200 ease-out data-ending-style:translate-x-full data-starting-style:translate-x-full">
            <div className="flex shrink-0 items-center justify-between border-b border-border-secondary px-6 py-4">
              <Drawer.Title className="font-serif text-lg font-semibold text-fg">
                {title}
              </Drawer.Title>
              <Drawer.Close className="cursor-pointer rounded-lg p-1.5 text-fg-tertiary hover:bg-bg-tertiary hover:text-fg">
                <XIcon size={18} />
              </Drawer.Close>
            </div>

            <form onSubmit={onSubmit} className="flex min-h-0 flex-1 flex-col">
              <div className="flex-1 space-y-5 overflow-y-auto px-6 py-5">{children}</div>

              {footer ?? (
                <div className="flex shrink-0 items-center justify-end gap-2 border-t border-border-secondary px-6 py-4">
                  <Drawer.Close className="cursor-pointer rounded-lg border border-border-secondary px-4 py-2 text-sm font-medium text-fg-secondary shadow-card hover:bg-bg-tertiary">
                    {t("admin.common.cancel")}
                  </Drawer.Close>
                  <button
                    type="submit"
                    disabled={isPending}
                    className="cursor-pointer rounded-lg bg-accent px-5 py-2 text-sm font-medium text-accent-fg shadow-card hover:opacity-90 disabled:opacity-50"
                  >
                    {isPending ? savingText : saveText}
                  </button>
                </div>
              )}
            </form>
          </Drawer.Popup>
        </Drawer.Viewport>
      </Drawer.Portal>
    </Drawer.Root>
  );
}
