import type { ReactNode, SyntheticEvent } from "react";

import { Drawer } from "@base-ui/react/drawer";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import { XIcon } from "@/components/ui/icons";
import { Tooltip } from "@/components/ui/tooltip";
import { useAuth } from "@/stores/auth";

interface DrawerShellProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  onSubmit: (e: SyntheticEvent) => void;
  isPending?: boolean;
  submitLabel?: string;
  savingLabel?: string;
  children: ReactNode;
  footer?: ReactNode;
  size?: "md" | "lg";
}

const sizeClasses = {
  md: "sm:max-w-[560px]",
  lg: "sm:max-w-[840px]",
};

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
  size = "lg",
}: DrawerShellProps) {
  const { t } = useTranslation();
  const { user } = useAuth();
  const isGuest = user?.role === "guest";
  const saveText = submitLabel ?? t("admin.common.save");
  const savingText = savingLabel ?? t("admin.common.saving");

  return (
    <Drawer.Root open={open} onOpenChange={onOpenChange} swipeDirection="right">
      <Drawer.Portal>
        <Drawer.Backdrop className="fixed inset-0 z-50 bg-black/25 backdrop-blur-xs transition-opacity duration-200 data-ending-style:opacity-0 data-starting-style:opacity-0" />
        <Drawer.Viewport className="fixed inset-0 z-50">
          <Drawer.Popup
            className={`fixed top-0 right-0 flex h-dvh w-full flex-col border-l border-border-secondary bg-bg-secondary shadow-elevated transition-transform duration-200 ease-out data-ending-style:translate-x-full data-starting-style:translate-x-full ${sizeClasses[size]}`}
          >
            <div className="flex shrink-0 items-center justify-between border-b border-border-secondary px-6 py-4">
              <Drawer.Title className="font-serif text-lg font-semibold text-fg">
                {title}
              </Drawer.Title>
              <Tooltip label={t("admin.common.close")}>
                <Drawer.Close className="cursor-pointer rounded-lg p-1.5 text-fg-tertiary hover:bg-bg-tertiary hover:text-fg">
                  <XIcon size={18} />
                </Drawer.Close>
              </Tooltip>
            </div>

            <form
              onSubmit={isGuest ? (e) => e.preventDefault() : onSubmit}
              className="flex min-h-0 flex-1 flex-col"
            >
              <div
                inert={isGuest || undefined}
                className="flex-1 space-y-5 overflow-y-auto px-6 py-5"
              >
                {children}
              </div>

              {footer ?? (
                <div className="flex shrink-0 items-center justify-end gap-2 border-t border-border-secondary px-6 py-4 pb-[max(1rem,env(safe-area-inset-bottom))]">
                  {isGuest ? (
                    <Button variant="secondary" asChild>
                      <Drawer.Close>{t("admin.common.close")}</Drawer.Close>
                    </Button>
                  ) : (
                    <>
                      <Button variant="secondary" asChild>
                        <Drawer.Close>{t("admin.common.cancel")}</Drawer.Close>
                      </Button>
                      <Button type="submit" disabled={isPending}>
                        {isPending ? savingText : saveText}
                      </Button>
                    </>
                  )}
                </div>
              )}
            </form>
          </Drawer.Popup>
        </Drawer.Viewport>
      </Drawer.Portal>
    </Drawer.Root>
  );
}
