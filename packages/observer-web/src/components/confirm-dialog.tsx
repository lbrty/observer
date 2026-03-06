import { Dialog } from "@base-ui/react/dialog";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";

interface ConfirmDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  description: string;
  onConfirm: () => void;
  loading?: boolean;
}

export function ConfirmDialog({
  open,
  onOpenChange,
  title,
  description,
  onConfirm,
  loading,
}: ConfirmDialogProps) {
  const { t } = useTranslation();

  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      <Dialog.Portal>
        <Dialog.Backdrop className="fixed inset-0 bg-black/25 backdrop-blur-xs" />
        <Dialog.Popup className="fixed top-1/2 left-1/2 w-full max-w-sm -translate-x-1/2 -translate-y-1/2 rounded-2xl border border-border-secondary bg-bg-secondary p-6 shadow-elevated">
          <Dialog.Title className="font-serif text-lg font-semibold text-fg">{title}</Dialog.Title>
          <Dialog.Description className="mt-2 text-sm text-fg-secondary">
            {description}
          </Dialog.Description>
          <div className="mt-6 flex justify-end gap-2">
            <Button variant="secondary" asChild>
              <Dialog.Close>{t("admin.common.cancel")}</Dialog.Close>
            </Button>
            <Button variant="danger" disabled={loading} onClick={onConfirm}>
              {t("admin.common.delete")}
            </Button>
          </div>
        </Dialog.Popup>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
