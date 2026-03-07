import { useRef, useState } from "react";

import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { Dialog } from "@base-ui/react/dialog";

import { ConfirmDialog } from "@/components/confirm-dialog";
import { DataTable, type Column } from "@/components/data-table";
import { EmptyState } from "@/components/empty-state";
import {
  CheckIcon,
  DownloadSimpleIcon,
  FileArchiveIcon,
  FileAudioIcon,
  FileCsvIcon,
  FileDocIcon,
  FileIcon,
  FileImageIcon,
  FilePdfIcon,
  FilePngIcon,
  FilePptIcon,
  FileSvgIcon,
  FileTextIcon,
  FileVideoIcon,
  FileXlsIcon,
  FilesIcon,
  PencilSimpleIcon,
  TrashIcon,
  UploadSimpleIcon,
  XIcon,
} from "@/components/icons";
import type { Icon } from "@/components/icons";
import {
  documentDownloadUrl,
  documentStreamUrl,
  documentThumbnailUrl,
  isImageMime,
  isPdfMime,
  useDeleteDocument,
  useDocuments,
  useUpdateDocument,
  useUploadDocument,
} from "@/hooks/use-documents";
import { handleApiError } from "@/lib/form-error";
import type { Document } from "@/types/document";

export const Route = createFileRoute("/_app/projects/$projectId/people/$personId/documents")({
  component: PersonDocuments,
});

function formatBytes(bytes: number): string {
  if (bytes === 0) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`;
}

function mimeIcon(mime: string): Icon {
  if (mime === "application/pdf") return FilePdfIcon;
  if (mime.startsWith("image/png")) return FilePngIcon;
  if (mime.startsWith("image/svg")) return FileSvgIcon;
  if (mime.startsWith("image/")) return FileImageIcon;
  if (mime.startsWith("video/")) return FileVideoIcon;
  if (mime.startsWith("audio/")) return FileAudioIcon;
  if (mime === "text/csv") return FileCsvIcon;
  if (mime.startsWith("text/")) return FileTextIcon;
  if (
    mime === "application/msword" ||
    mime === "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
  )
    return FileDocIcon;
  if (
    mime === "application/vnd.ms-excel" ||
    mime === "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
  )
    return FileXlsIcon;
  if (
    mime === "application/vnd.ms-powerpoint" ||
    mime === "application/vnd.openxmlformats-officedocument.presentationml.presentation"
  )
    return FilePptIcon;
  if (mime === "application/zip" || mime === "application/x-rar-compressed" || mime === "application/gzip")
    return FileArchiveIcon;
  return FileIcon;
}

function PersonDocuments() {
  const { t } = useTranslation();
  const { projectId, personId } = Route.useParams();

  const { data, isLoading } = useDocuments(projectId, personId);
  const updateDocument = useUpdateDocument(projectId);
  const deleteDocument = useDeleteDocument(projectId);
  const uploadDocument = useUploadDocument(projectId, personId);

  const fileInputRef = useRef<HTMLInputElement>(null);

  const [deleteId, setDeleteId] = useState<string | null>(null);
  const [editId, setEditId] = useState<string | null>(null);
  const [editName, setEditName] = useState("");
  const [uploadError, setUploadError] = useState("");
  const [previewDoc, setPreviewDoc] = useState<Document | null>(null);

  function startEdit(doc: Document) {
    setEditId(doc.id);
    setEditName(doc.name);
  }

  function cancelEdit() {
    setEditId(null);
    setEditName("");
  }

  function saveEdit() {
    if (!editId || !editName.trim()) return;
    updateDocument.mutate(
      { id: editId, data: { name: editName.trim() } },
      { onSuccess: () => cancelEdit() },
    );
  }

  function handleDelete() {
    if (!deleteId) return;
    deleteDocument.mutate(deleteId, { onSuccess: () => setDeleteId(null) });
  }

  async function handleFileSelect(e: React.ChangeEvent<HTMLInputElement>) {
    const files = e.target.files;
    if (!files?.length) return;

    setUploadError("");

    for (const file of Array.from(files)) {
      try {
        await uploadDocument.mutateAsync(file);
      } catch (err) {
        setUploadError(await handleApiError(err, t));
        break;
      }
    }

    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  }

  const columns: Column<Document>[] = [
    {
      key: "preview",
      header: "",
      render: (doc) => {
        if (isImageMime(doc.mime_type)) {
          return (
            <button
              type="button"
              onClick={() => setPreviewDoc(doc)}
              className="w-12 cursor-pointer overflow-hidden rounded"
              style={{ aspectRatio: "4 / 3" }}
            >
              <img
                src={documentThumbnailUrl(projectId, doc.id)}
                alt={doc.name}
                className="h-full w-full object-cover"
                loading="lazy"
              />
            </button>
          );
        }
        const IconComponent = mimeIcon(doc.mime_type);
        const clickable = isPdfMime(doc.mime_type);
        return clickable ? (
          <button
            type="button"
            onClick={() => setPreviewDoc(doc)}
            className="cursor-pointer text-fg-tertiary hover:text-fg"
          >
            <IconComponent size={28} />
          </button>
        ) : (
          <span className="text-fg-tertiary">
            <IconComponent size={28} />
          </span>
        );
      },
    },
    {
      key: "name",
      header: t("project.documents.name"),
      render: (doc) =>
        editId === doc.id ? (
          <div className="flex items-center gap-1">
            <input
              type="text"
              value={editName}
              onChange={(e) => setEditName(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter") saveEdit();
                if (e.key === "Escape") cancelEdit();
              }}
              className="h-7 rounded border border-accent bg-bg-secondary px-2 text-sm text-fg outline-none"
              autoFocus
            />
            <button
              type="button"
              onClick={saveEdit}
              disabled={!editName.trim() || updateDocument.isPending}
              className="cursor-pointer rounded p-1 text-accent hover:bg-bg-tertiary disabled:opacity-50"
            >
              <CheckIcon size={14} />
            </button>
            <button
              type="button"
              onClick={cancelEdit}
              className="cursor-pointer rounded p-1 text-fg-tertiary hover:bg-bg-tertiary"
            >
              <XIcon size={14} />
            </button>
          </div>
        ) : (
          <span className="font-medium text-fg">{doc.name}</span>
        ),
    },
    {
      key: "size",
      header: t("project.documents.size"),
      render: (doc) => (
        <span className="font-mono text-xs tabular-nums text-fg-tertiary">
          {formatBytes(doc.size)}
        </span>
      ),
    },
    {
      key: "created_at",
      header: t("project.people.registered"),
      render: (doc) => (
        <span className="font-mono text-xs tabular-nums text-fg-tertiary">
          {new Date(doc.created_at).toLocaleDateString("en-CA")}
        </span>
      ),
    },
    {
      key: "actions",
      header: "",
      render: (doc) => (
        <div className="flex gap-1">
          <a
            href={documentDownloadUrl(projectId, doc.id)}
            title={t("project.documents.download")}
            className="cursor-pointer rounded-lg p-1.5 text-fg-tertiary hover:bg-bg-tertiary hover:text-fg"
          >
            <DownloadSimpleIcon size={16} />
          </a>
          <button
            type="button"
            onClick={(e) => {
              e.stopPropagation();
              startEdit(doc);
            }}
            className="cursor-pointer rounded-lg p-1.5 text-fg-tertiary hover:bg-bg-tertiary hover:text-fg"
          >
            <PencilSimpleIcon size={16} />
          </button>
          <button
            type="button"
            onClick={(e) => {
              e.stopPropagation();
              setDeleteId(doc.id);
            }}
            className="cursor-pointer rounded-lg p-1.5 text-fg-tertiary hover:bg-bg-tertiary hover:text-rose"
          >
            <TrashIcon size={16} />
          </button>
        </div>
      ),
    },
  ];

  return (
    <div>
      <div className="mb-4 flex items-center justify-between">
        <h2 className="font-serif text-lg font-semibold text-fg">
          {t("project.documents.title")}
        </h2>
        <div>
          <input
            ref={fileInputRef}
            type="file"
            multiple
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
        </div>
      </div>

      {uploadError && (
        <div className="mb-4 rounded-lg bg-rose/10 px-3 py-2 text-sm text-rose">{uploadError}</div>
      )}

      <DataTable
        columns={columns}
        data={data?.documents ?? []}
        keyExtractor={(doc) => doc.id}
        isLoading={isLoading}
        emptyState={
          <EmptyState
            icon={FilesIcon}
            title={t("project.people.documentsEmptyTitle")}
            description={t("project.people.documentsEmptyDescription")}
          />
        }
      />

      <ConfirmDialog
        open={!!deleteId}
        onOpenChange={(open) => {
          if (!open) setDeleteId(null);
        }}
        title={t("admin.common.delete")}
        description={t("project.documents.deleteConfirm")}
        onConfirm={handleDelete}
        loading={deleteDocument.isPending}
      />

      <Dialog.Root
        open={!!previewDoc}
        onOpenChange={(open) => {
          if (!open) setPreviewDoc(null);
        }}
      >
        <Dialog.Portal>
          <Dialog.Backdrop className="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm" />
          <Dialog.Popup
            className="fixed inset-0 z-50 flex cursor-pointer items-center justify-center p-8"
            onClick={() => setPreviewDoc(null)}
          >
            <div
              className={`relative cursor-default ${
                previewDoc && isPdfMime(previewDoc.mime_type)
                  ? "flex h-[90vh] w-[60vw] flex-col"
                  : "max-h-full max-w-full"
              }`}
              onClick={(e) => e.stopPropagation()}
            >
              <Dialog.Close className="absolute -top-3 -right-3 z-10 cursor-pointer rounded-full bg-bg-secondary p-1.5 text-fg-secondary shadow-elevated hover:text-fg">
                <XIcon size={18} />
              </Dialog.Close>
              {previewDoc && isPdfMime(previewDoc.mime_type) ? (
                <iframe
                  src={documentStreamUrl(projectId, previewDoc.id)}
                  title={previewDoc.name}
                  className="h-full w-full flex-1 rounded-lg shadow-elevated"
                />
              ) : previewDoc ? (
                <img
                  src={documentStreamUrl(projectId, previewDoc.id)}
                  alt={previewDoc.name}
                  className="max-h-[80vh] max-w-[80vw] rounded-lg object-contain shadow-elevated"
                />
              ) : null}
              {previewDoc && (
                <p className="mt-2 text-center text-sm text-white/70">{previewDoc.name}</p>
              )}
            </div>
          </Dialog.Popup>
        </Dialog.Portal>
      </Dialog.Root>
    </div>
  );
}
