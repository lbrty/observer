import { Dialog } from "@base-ui/react/dialog";
import { useTranslation } from "react-i18next";

import { ErrorBanner } from "@/components/alert-banner";
import { Button } from "@/components/button";

interface AddReferenceDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  children: React.ReactNode;
  onSubmit: () => void;
  isPending: boolean;
  error: string;
}

export function AddReferenceDialog({
  open,
  onOpenChange,
  title,
  children,
  onSubmit,
  isPending,
  error,
}: AddReferenceDialogProps) {
  const { t } = useTranslation();

  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      <Dialog.Portal>
        <Dialog.Backdrop className="fixed inset-0 z-50 bg-black/25 backdrop-blur-xs" />
        <Dialog.Popup className="fixed top-1/2 left-1/2 z-50 w-full max-w-md -translate-x-1/2 -translate-y-1/2 rounded-xl border border-border-secondary bg-bg-secondary p-6 shadow-elevated">
          <Dialog.Title className="font-serif text-lg font-semibold text-fg">{title}</Dialog.Title>
          <form
            onSubmit={(e) => {
              e.preventDefault();
              onSubmit();
            }}
          >
            {error && (
              <div className="mt-4">
                <ErrorBanner message={error} />
              </div>
            )}
            <div className="mt-4">{children}</div>
            <div className="mt-6 flex justify-end gap-2">
              <Button variant="secondary" asChild>
                <Dialog.Close>{t("admin.common.cancel")}</Dialog.Close>
              </Button>
              <Button type="submit" disabled={isPending}>
                {isPending ? t("project.people.saving") : t("project.people.save")}
              </Button>
            </div>
          </form>
        </Dialog.Popup>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
