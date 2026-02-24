import type { FormEvent, ReactNode } from "react";

import { Dialog } from "@base-ui/react/dialog";
import { useTranslation } from "react-i18next";

interface FormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  description?: string;
  loading?: boolean;
  onSubmit: (e: FormEvent) => void;
  children: ReactNode;
  maxWidth?: "sm" | "md";
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
}: FormDialogProps) {
  const { t } = useTranslation();

  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      <Dialog.Portal>
        <Dialog.Backdrop className="fixed inset-0 bg-black/25 backdrop-blur-xs" />
        <Dialog.Popup
          className={`fixed top-1/2 left-1/2 w-full -translate-x-1/2 -translate-y-1/2 rounded-2xl border border-border-secondary bg-bg-secondary p-6 shadow-elevated ${maxWidth === "md" ? "max-w-md" : "max-w-sm"}`}
        >
          <Dialog.Title className="font-serif text-lg font-semibold text-fg">
            {title}
          </Dialog.Title>
          {description && (
            <Dialog.Description className="mt-1 text-sm text-fg-tertiary">
              {description}
            </Dialog.Description>
          )}
          <form onSubmit={onSubmit} className="mt-5 space-y-4">
            {children}
            <div className="flex justify-end gap-2 pt-2">
              <Dialog.Close className="cursor-pointer rounded-lg border border-border-secondary px-4 py-2 text-sm font-medium text-fg-secondary shadow-card hover:bg-bg-tertiary">
                {t("admin.common.cancel")}
              </Dialog.Close>
              <button
                type="submit"
                disabled={loading}
                className="cursor-pointer rounded-lg bg-accent px-4 py-2 text-sm font-medium text-accent-fg shadow-card hover:opacity-90 disabled:opacity-50"
              >
                {loading ? t("admin.common.saving") : t("admin.common.save")}
              </button>
            </div>
          </form>
        </Dialog.Popup>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
