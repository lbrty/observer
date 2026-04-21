import type { TFunction } from "i18next";

import { type Column } from "@/components/table/data-table";
import { formatBytes, mimeIcon } from "@/components/documents/document-mime-icon";
import {
  CheckIcon,
  DownloadSimpleIcon,
  PencilSimpleIcon,
  TrashIcon,
  XIcon,
} from "@/components/ui/icons";
import {
  documentDownloadUrl,
  documentThumbnailUrl,
  isImageMime,
  isPdfMime,
} from "@/hooks/documents/use-documents";
import type { Document } from "@/types/document";

interface DocumentColumnsConfig {
  t: TFunction;
  projectId: string;
  canWrite: boolean;
  canDelete: boolean;
  editId: string | null;
  editName: string;
  setEditName: (name: string) => void;
  startEdit: (doc: Document) => void;
  saveEdit: () => void;
  cancelEdit: () => void;
  setDeleteId: (id: string | null) => void;
  setPreviewDoc: (doc: Document | null) => void;
  updateDocument: { isPending: boolean };
}

export function buildDocumentColumns(cfg: DocumentColumnsConfig): Column<Document>[] {
  return [
    {
      key: "preview",
      header: "",
      render: (doc) => {
        if (isImageMime(doc.mime_type)) {
          return (
            <button
              type="button"
              onClick={() => cfg.setPreviewDoc(doc)}
              className="w-12 cursor-pointer overflow-hidden rounded"
              style={{ aspectRatio: "4 / 3" }}
            >
              <img
                src={documentThumbnailUrl(cfg.projectId, doc.id)}
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
            onClick={() => cfg.setPreviewDoc(doc)}
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
      header: cfg.t("project.documents.name"),
      render: (doc) =>
        cfg.editId === doc.id ? (
          <div className="flex items-center gap-1">
            <input
              type="text"
              value={cfg.editName}
              onChange={(e) => cfg.setEditName(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter") cfg.saveEdit();
                if (e.key === "Escape") cfg.cancelEdit();
              }}
              className="h-7 rounded border border-accent bg-bg-secondary px-2 text-sm text-fg outline-none"
              autoFocus
            />
            <button
              type="button"
              onClick={cfg.saveEdit}
              disabled={!cfg.editName.trim() || cfg.updateDocument.isPending}
              className="cursor-pointer rounded p-1 text-accent hover:bg-bg-tertiary disabled:opacity-50"
            >
              <CheckIcon size={14} />
            </button>
            <button
              type="button"
              onClick={cfg.cancelEdit}
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
      header: cfg.t("project.documents.size"),
      render: (doc) => (
        <span className="font-mono text-xs tabular-nums text-fg-tertiary">
          {formatBytes(doc.size)}
        </span>
      ),
    },
    {
      key: "created_at",
      header: cfg.t("project.people.registered"),
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
            href={documentDownloadUrl(cfg.projectId, doc.id)}
            title={cfg.t("project.documents.download")}
            className="cursor-pointer rounded-lg p-1.5 text-fg-tertiary hover:bg-bg-tertiary hover:text-fg"
          >
            <DownloadSimpleIcon size={16} />
          </a>
          {cfg.canWrite && (
            <button
              type="button"
              onClick={(e) => {
                e.stopPropagation();
                cfg.startEdit(doc);
              }}
              className="cursor-pointer rounded-lg p-1.5 text-fg-tertiary hover:bg-bg-tertiary hover:text-fg"
            >
              <PencilSimpleIcon size={16} />
            </button>
          )}
          {cfg.canDelete && (
            <button
              type="button"
              onClick={(e) => {
                e.stopPropagation();
                cfg.setDeleteId(doc.id);
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
}
