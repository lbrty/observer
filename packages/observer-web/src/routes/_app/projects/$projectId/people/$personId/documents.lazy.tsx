import { useState } from "react";

import { createLazyFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { ConfirmDialog } from "@/components/confirm-dialog";
import { DataTable, type Column } from "@/components/data-table";
import { formatBytes, mimeIcon } from "@/components/document-mime-icon";
import { DocumentPreviewDialog } from "@/components/document-preview-dialog";
import { DocumentUploadZone } from "@/components/document-upload-zone";
import { EmptyState } from "@/components/empty-state";
import {
  CheckIcon,
  DownloadSimpleIcon,
  FilesIcon,
  PencilSimpleIcon,
  TrashIcon,
  XIcon,
} from "@/components/icons";
import {
  documentDownloadUrl,
  documentThumbnailUrl,
  isImageMime,
  isPdfMime,
  useDeleteDocument,
  useDocuments,
  useUpdateDocument,
} from "@/hooks/use-documents";
import { useProjectRole } from "@/hooks/use-project-role";
import type { Document } from "@/types/document";

export const Route = createLazyFileRoute("/_app/projects/$projectId/people/$personId/documents")({
  component: PersonDocuments,
});

function PersonDocuments() {
  const { t } = useTranslation();
  const { projectId, personId } = Route.useParams();

  const { canWrite, canDelete } = useProjectRole(projectId);
  const { data, isLoading } = useDocuments(projectId, personId);
  const updateDocument = useUpdateDocument(projectId);
  const deleteDocument = useDeleteDocument(projectId);

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
          {canWrite && (
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
          )}
          {canDelete && (
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
          )}
        </div>
      ),
    },
  ];

  return (
    <div>
      <div className="mb-4 flex items-center justify-between">
        <h2 className="font-serif text-lg font-semibold text-fg">{t("project.documents.title")}</h2>
        <div>
          {canWrite && (
            <DocumentUploadZone
              projectId={projectId}
              personId={personId}
              onUploadError={setUploadError}
            />
          )}
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

      <DocumentPreviewDialog
        document={previewDoc}
        projectId={projectId}
        onClose={() => setPreviewDoc(null)}
      />
    </div>
  );
}
