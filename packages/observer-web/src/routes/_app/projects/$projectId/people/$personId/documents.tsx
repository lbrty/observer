import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { ConfirmDialog } from "@/components/confirm-dialog";
import { DataTable, type Column } from "@/components/data-table";
import { TrashIcon } from "@/components/icons";
import { useDeleteDocument, useDocuments } from "@/hooks/use-documents";
import type { Document } from "@/types/document";

export const Route = createFileRoute(
  "/_app/projects/$projectId/people/$personId/documents",
)({
  component: PersonDocuments,
});

function formatBytes(bytes: number): string {
  if (bytes === 0) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`;
}

function PersonDocuments() {
  const { t } = useTranslation();
  const { projectId, personId } = Route.useParams();

  const { data, isLoading } = useDocuments(projectId, personId);
  const deleteDocument = useDeleteDocument(projectId);

  const [deleteId, setDeleteId] = useState<string | null>(null);

  function handleDelete() {
    if (!deleteId) return;
    deleteDocument.mutate(deleteId, { onSuccess: () => setDeleteId(null) });
  }

  const columns: Column<Document>[] = [
    {
      key: "name",
      header: t("project.documents.name"),
      render: (doc) => (
        <span className="font-medium text-fg">{doc.name}</span>
      ),
    },
    {
      key: "mime_type",
      header: t("project.documents.mimeType"),
      render: (doc) => (
        <span className="text-fg-secondary">{doc.mime_type}</span>
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
      ),
    },
  ];

  return (
    <div>
      <h2 className="mb-4 font-serif text-lg font-semibold text-fg">
        {t("project.documents.title")}
      </h2>

      <DataTable
        columns={columns}
        data={data?.documents ?? []}
        keyExtractor={(doc) => doc.id}
        isLoading={isLoading}
      />

      {!isLoading && !data?.documents.length && (
        <p className="py-12 text-center text-sm text-fg-tertiary">
          {t("project.documents.empty")}
        </p>
      )}

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
    </div>
  );
}
