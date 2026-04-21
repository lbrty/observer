import { useState } from "react";

import { createLazyFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { ConfirmDialog } from "@/components/dialogs/confirm-dialog";
import { DataTable } from "@/components/table/data-table";
import { buildDocumentColumns } from "@/components/documents/document-columns";
import { DocumentPreviewDialog } from "@/components/documents/document-preview-dialog";
import { DocumentUploadZone } from "@/components/documents/document-upload-zone";
import { EmptyState } from "@/components/ui/empty-state";
import { FilesIcon } from "@/components/ui/icons";
import { useDeleteDocument, useDocuments, useUpdateDocument } from "@/hooks/documents/use-documents";
import { useProjectRole } from "@/hooks/users/use-project-role";
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

  const columns = buildDocumentColumns({
    t,
    projectId,
    canWrite,
    canDelete,
    editId,
    editName,
    setEditName,
    startEdit,
    saveEdit,
    cancelEdit,
    setDeleteId,
    setPreviewDoc,
    updateDocument,
  });

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
