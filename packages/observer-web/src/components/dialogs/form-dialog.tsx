import type { SyntheticEvent, ReactNode } from "react";

import { Dialog } from "@base-ui/react/dialog";
import { useTranslation } from "react-i18next";

import { ErrorBanner } from "@/components/alert-banner";
import { Button } from "@/components/button";

interface FormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  description?: string;
  loading?: boolean;
  onSubmit: (e: SyntheticEvent) => void;
  children: ReactNode;
  maxWidth?: "sm" | "md";
  error?: string;
  isPending?: boolean;
}

export function FormDialog({
  open,
  onOpenChange,
  title,
  description,
  loading,
  onSubmit,
  children,
  maxWidth = "sm",
  error,
  isPending,
}: FormDialogProps) {
  const { t } = useTranslation();

  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      <Dialog.Portal>
        <Dialog.Backdrop className="fixed inset-0 bg-black/25 backdrop-blur-xs" />
        <Dialog.Popup
          className={`fixed top-1/2 left-1/2 w-full -translate-x-1/2 -translate-y-1/2 rounded-2xl border border-border-secondary bg-bg-secondary p-6 shadow-elevated ${maxWidth === "md" ? "max-w-md" : "max-w-sm"}`}
        >
          <Dialog.Title className="font-serif text-lg font-semibold text-fg">{title}</Dialog.Title>
          {description && (
            <Dialog.Description className="mt-1 text-sm text-fg-tertiary">
              {description}
            </Dialog.Description>
          )}
          <form onSubmit={onSubmit} className="mt-5 space-y-4">
            {error && <ErrorBanner message={error} />}
            {children}
            <div className="flex justify-end gap-2 pt-2">
              <Button variant="secondary" asChild>
                <Dialog.Close>{t("admin.common.cancel")}</Dialog.Close>
              </Button>
              <Button type="submit" disabled={loading || isPending}>
                {loading || isPending ? t("admin.common.saving") : t("admin.common.save")}
              </Button>
            </div>
          </form>
        </Dialog.Popup>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
