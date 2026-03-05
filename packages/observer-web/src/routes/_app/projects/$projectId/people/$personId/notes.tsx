import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { ConfirmDialog } from "@/components/confirm-dialog";
import { TrashIcon } from "@/components/icons";
import { useCreateNote, useDeleteNote, useNotes } from "@/hooks/use-notes";

export const Route = createFileRoute(
  "/_app/projects/$projectId/people/$personId/notes",
)({
  component: PersonNotes,
});

function PersonNotes() {
  const { t } = useTranslation();
  const { projectId, personId } = Route.useParams();

  const { data, isLoading } = useNotes(projectId, personId);
  const createNote = useCreateNote(projectId, personId);
  const deleteNote = useDeleteNote(projectId, personId);

  const [body, setBody] = useState("");
  const [deleteId, setDeleteId] = useState<string | null>(null);

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    const trimmed = body.trim();
    if (!trimmed) return;
    createNote.mutate({ body: trimmed }, { onSuccess: () => setBody("") });
  }

  function handleDelete() {
    if (!deleteId) return;
    deleteNote.mutate(deleteId, { onSuccess: () => setDeleteId(null) });
  }

  return (
    <div>
      <h2 className="mb-4 font-serif text-lg font-semibold text-fg">
        {t("project.notes.title")}
      </h2>

      <form onSubmit={handleSubmit} className="mb-6">
        <textarea
          value={body}
          onChange={(e) => setBody(e.target.value)}
          placeholder={t("project.notes.body")}
          rows={3}
          className="w-full rounded-xl border border-border-secondary bg-bg-secondary px-4 py-3 text-sm text-fg outline-none transition-colors focus:border-accent"
        />
        <div className="mt-2 flex justify-end">
          <button
            type="submit"
            disabled={!body.trim() || createNote.isPending}
            className="cursor-pointer rounded-lg bg-accent px-4 py-2 text-sm font-medium text-accent-fg shadow-card hover:opacity-90 disabled:cursor-default disabled:opacity-50"
          >
            {t("project.notes.add")}
          </button>
        </div>
      </form>

      {isLoading ? (
        <div className="space-y-3">
          {Array.from({ length: 3 }, (_, i) => (
            <div
              key={i}
              className="h-24 animate-pulse rounded-xl bg-bg-tertiary"
            />
          ))}
        </div>
      ) : !data?.notes.length ? (
        <p className="py-12 text-center text-sm text-fg-tertiary">
          {t("project.notes.empty")}
        </p>
      ) : (
        <div className="space-y-3">
          {data.notes.map((note) => (
            <div
              key={note.id}
              className="rounded-xl border border-border-secondary bg-bg-secondary p-4"
            >
              <p className="whitespace-pre-wrap text-sm text-fg">{note.body}</p>
              <div className="mt-2 flex items-center justify-between">
                <span className="font-mono text-xs tabular-nums text-fg-tertiary">
                  {new Date(note.created_at).toLocaleString("en-CA")}
                </span>
                <button
                  type="button"
                  onClick={() => setDeleteId(note.id)}
                  className="cursor-pointer rounded-lg p-1.5 text-fg-tertiary hover:bg-bg-tertiary hover:text-rose"
                >
                  <TrashIcon size={14} />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      <ConfirmDialog
        open={!!deleteId}
        onOpenChange={(open) => {
          if (!open) setDeleteId(null);
        }}
        title={t("admin.common.delete")}
        description={t("project.notes.deleteConfirm")}
        onConfirm={handleDelete}
        loading={deleteNote.isPending}
      />
    </div>
  );
}
