import { useRef } from "react";

import { useTranslation } from "react-i18next";

import { UploadSimpleIcon } from "@/components/icons";
import { useUploadDocument } from "@/hooks/use-documents";
import { handleApiError } from "@/lib/form-error";

interface DocumentUploadZoneProps {
  projectId: string;
  personId: string;
  onUploadError: (msg: string) => void;
}

export function DocumentUploadZone({ projectId, personId, onUploadError }: DocumentUploadZoneProps) {
  const { t } = useTranslation();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const uploadDocument = useUploadDocument(projectId, personId);

  async function handleFileSelect(e: React.ChangeEvent<HTMLInputElement>) {
    const files = e.target.files;
    if (!files?.length) return;

    onUploadError("");

    for (const file of Array.from(files)) {
      try {
        await uploadDocument.mutateAsync(file);
      } catch (err) {
        onUploadError(await handleApiError(err, t));
        break;
      }
    }

    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  }

  return (
    <>
      <input
        ref={fileInputRef}
        type="file"
        multiple
        accept=".pdf,image/*,video/*,audio/*,text/*,.doc,.docx,.xls,.xlsx,.ppt,.pptx,.zip,.rar,.gz"
        onChange={handleFileSelect}
        className="hidden"
      />
      <button
        type="button"
        onClick={() => fileInputRef.current?.click()}
        disabled={uploadDocument.isPending}
        className="inline-flex cursor-pointer items-center gap-1.5 rounded-lg border border-border-secondary bg-bg-secondary px-3 py-1.5 text-sm font-medium text-fg hover:bg-bg-tertiary disabled:opacity-50"
      >
        <UploadSimpleIcon size={16} />
        {uploadDocument.isPending
          ? t("project.documents.uploading")
          : t("project.documents.upload")}
      </button>
    </>
  );
}
