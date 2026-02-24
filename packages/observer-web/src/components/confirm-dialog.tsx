import { Dialog } from "@base-ui/react/dialog";
import { useTranslation } from "react-i18next";

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
        <Dialog.Backdrop className="fixed inset-0 bg-black/40" />
        <Dialog.Popup className="fixed top-1/2 left-1/2 w-full max-w-sm -translate-x-1/2 -translate-y-1/2 rounded-lg border border-border-secondary bg-bg-secondary p-6 shadow-elevated">
          <Dialog.Title className="text-lg font-semibold text-fg">
            {title}
          </Dialog.Title>
          <Dialog.Description className="mt-2 text-sm text-fg-secondary">
            {description}
          </Dialog.Description>
          <div className="mt-6 flex justify-end gap-3">
            <Dialog.Close className="cursor-pointer rounded-md border border-border-secondary px-3 py-1.5 text-sm text-fg-secondary hover:bg-bg-tertiary">
              {t("admin.common.cancel")}
            </Dialog.Close>
            <button
              type="button"
              disabled={loading}
              onClick={onConfirm}
              className="cursor-pointer rounded-md bg-rose px-3 py-1.5 text-sm font-medium text-white hover:opacity-90 disabled:opacity-50"
            >
              {t("admin.common.delete")}
            </button>
          </div>
        </Dialog.Popup>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
