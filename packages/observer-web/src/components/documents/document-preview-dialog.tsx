import { Dialog } from "@base-ui/react/dialog";

import { XIcon } from "@/components/ui/icons";
import { documentStreamUrl, isPdfMime } from "@/hooks/documents/use-documents";
import type { Document } from "@/types/document";

interface DocumentPreviewDialogProps {
  document: Document | null;
  projectId: string;
  onClose: () => void;
}

export function DocumentPreviewDialog({
  document,
  projectId,
  onClose,
}: DocumentPreviewDialogProps) {
  return (
    <Dialog.Root
      open={!!document}
      onOpenChange={(open) => {
        if (!open) onClose();
      }}
    >
      <Dialog.Portal>
        <Dialog.Backdrop className="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm" />
        <Dialog.Popup
          className="fixed inset-0 z-50 flex cursor-pointer items-center justify-center p-8"
          onClick={onClose}
        >
          <div
            className={`relative cursor-default ${
              document && isPdfMime(document.mime_type)
                ? "flex h-[90vh] w-[60vw] flex-col"
                : "max-h-full max-w-full"
            }`}
            onClick={(e) => e.stopPropagation()}
          >
            <Dialog.Close className="absolute -top-3 -right-3 z-10 cursor-pointer rounded-full bg-bg-secondary p-1.5 text-fg-secondary shadow-elevated hover:text-fg">
              <XIcon size={18} />
            </Dialog.Close>
            {document && isPdfMime(document.mime_type) ? (
              <iframe
                src={documentStreamUrl(projectId, document.id)}
                title={document.name}
                className="h-full w-full flex-1 rounded-lg shadow-elevated"
              />
            ) : document ? (
              <img
                src={documentStreamUrl(projectId, document.id)}
                alt={document.name}
                className="max-h-[80vh] max-w-[80vw] rounded-lg object-contain shadow-elevated"
              />
            ) : null}
            {document && <p className="mt-2 text-center text-sm text-white/70">{document.name}</p>}
          </div>
        </Dialog.Popup>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
