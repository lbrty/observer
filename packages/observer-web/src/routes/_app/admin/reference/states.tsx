import { Dialog } from "@base-ui/react/dialog";
import { Field } from "@base-ui/react/field";
import { createFileRoute } from "@tanstack/react-router";
import { type FormEvent, useState } from "react";
import { useTranslation } from "react-i18next";

import { ConfirmDialog } from "@/components/confirm-dialog";
import { DataTable, type Column } from "@/components/data-table";
import { PageHeader } from "@/components/page-header";
import { StatusBadge } from "@/components/status-badge";
import { useCountries } from "@/hooks/use-countries";
import {
  useCreateState,
  useDeleteState,
  useStates,
  useUpdateState,
} from "@/hooks/use-states";
import type { State } from "@/types/reference";

export const Route = createFileRoute("/_app/admin/reference/states")({
  component: StatesPage,
});

function StatesPage() {
  const { t } = useTranslation();
  const { data: countriesData } = useCountries();
  const countries = countriesData?.countries ?? [];

  const [countryId, setCountryId] = useState("");
  const { data, isLoading } = useStates(countryId);
  const createState = useCreateState();
  const updateState = useUpdateState();
  const deleteState = useDeleteState();

  const [createOpen, setCreateOpen] = useState(false);
  const [editTarget, setEditTarget] = useState<State | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<State | null>(null);

  const columns: Column<State>[] = [
    {
      key: "name",
      header: t("admin.reference.states.name"),
      render: (s) => <span className="font-medium text-fg">{s.name}</span>,
    },
    {
      key: "code",
      header: t("admin.reference.states.code"),
      render: (s) => <span className="text-fg-secondary">{s.code ?? "—"}</span>,
    },
    {
      key: "conflict_zone",
      header: t("admin.reference.states.conflictZone"),
      render: (s) =>
        s.conflict_zone ? (
          <StatusBadge label={s.conflict_zone} variant="rose" />
        ) : (
          <span className="text-fg-tertiary">—</span>
        ),
    },
    {
      key: "actions",
      header: t("admin.common.actions"),
      render: (s) => (
        <div className="flex gap-2">
          <button
            type="button"
            onClick={() => setEditTarget(s)}
            className="cursor-pointer text-xs text-accent hover:underline"
          >
            {t("admin.common.edit")}
          </button>
          <button
            type="button"
            onClick={() => setDeleteTarget(s)}
            className="cursor-pointer text-xs text-rose hover:underline"
          >
            {t("admin.common.delete")}
          </button>
        </div>
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title={t("admin.reference.states.title")}
        action={
          countryId ? (
            <button
              type="button"
              onClick={() => setCreateOpen(true)}
              className="cursor-pointer rounded-md bg-accent px-3 py-1.5 text-sm font-medium text-accent-fg hover:opacity-90"
            >
              {t("admin.reference.states.add")}
            </button>
          ) : undefined
        }
      />

      <div className="mb-4">
        <select
          value={countryId}
          onChange={(e) => setCountryId(e.target.value)}
          className="rounded-md border border-border-secondary bg-bg-secondary px-3 py-1.5 text-sm text-fg outline-none"
        >
          <option value="">{t("admin.reference.states.selectCountry")}</option>
          {countries.map((c) => (
            <option key={c.id} value={c.id}>
              {c.name}
            </option>
          ))}
        </select>
      </div>

      {countryId && (
        <DataTable
          columns={columns}
          data={data?.states ?? []}
          keyExtractor={(s) => s.id}
          isLoading={isLoading}
        />
      )}

      <StateFormDialog
        open={createOpen}
        onOpenChange={setCreateOpen}
        title={t("admin.reference.states.addTitle")}
        onSubmit={async (name, code, conflictZone) => {
          await createState.mutateAsync({
            countryId,
            data: {
              name,
              code: code || undefined,
              conflict_zone: conflictZone || undefined,
            },
          });
          setCreateOpen(false);
        }}
        loading={createState.isPending}
      />

      {editTarget && (
        <StateFormDialog
          open={!!editTarget}
          onOpenChange={(open) => !open && setEditTarget(null)}
          title={t("admin.reference.states.editTitle")}
          initial={editTarget}
          onSubmit={async (name, code, conflictZone) => {
            await updateState.mutateAsync({
              id: editTarget.id,
              data: {
                name,
                code: code || undefined,
                conflict_zone: conflictZone || undefined,
              },
            });
            setEditTarget(null);
          }}
          loading={updateState.isPending}
        />
      )}

      <ConfirmDialog
        open={!!deleteTarget}
        onOpenChange={(open) => !open && setDeleteTarget(null)}
        title={t("admin.common.delete")}
        description={t("admin.reference.states.deleteConfirm")}
        onConfirm={async () => {
          if (deleteTarget) {
            await deleteState.mutateAsync(deleteTarget.id);
            setDeleteTarget(null);
          }
        }}
        loading={deleteState.isPending}
      />
    </div>
  );
}

function StateFormDialog({
  open,
  onOpenChange,
  title,
  initial,
  onSubmit,
  loading,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  initial?: State;
  onSubmit: (name: string, code: string, conflictZone: string) => Promise<void>;
  loading: boolean;
}) {
  const { t } = useTranslation();
  const [name, setName] = useState(initial?.name ?? "");
  const [code, setCode] = useState(initial?.code ?? "");
  const [conflictZone, setConflictZone] = useState(
    initial?.conflict_zone ?? "",
  );

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    await onSubmit(name, code, conflictZone);
    if (!initial) {
      setName("");
      setCode("");
      setConflictZone("");
    }
  }

  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      <Dialog.Portal>
        <Dialog.Backdrop className="fixed inset-0 bg-black/40" />
        <Dialog.Popup className="fixed top-1/2 left-1/2 w-full max-w-sm -translate-x-1/2 -translate-y-1/2 rounded-lg border border-border-secondary bg-bg-secondary p-6 shadow-elevated">
          <Dialog.Title className="text-lg font-semibold text-fg">
            {title}
          </Dialog.Title>
          <form onSubmit={handleSubmit} className="mt-4 space-y-3">
            <Field.Root>
              <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                {t("admin.reference.states.name")}
              </Field.Label>
              <Field.Control
                required
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="block w-full rounded-md border border-border-secondary bg-bg px-3 py-2 text-sm text-fg outline-none focus:border-accent"
              />
            </Field.Root>
            <Field.Root>
              <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                {t("admin.reference.states.code")}
              </Field.Label>
              <Field.Control
                value={code}
                onChange={(e) => setCode(e.target.value)}
                className="block w-full rounded-md border border-border-secondary bg-bg px-3 py-2 text-sm text-fg outline-none focus:border-accent"
              />
            </Field.Root>
            <Field.Root>
              <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                {t("admin.reference.states.conflictZone")}
              </Field.Label>
              <Field.Control
                value={conflictZone}
                onChange={(e) => setConflictZone(e.target.value)}
                className="block w-full rounded-md border border-border-secondary bg-bg px-3 py-2 text-sm text-fg outline-none focus:border-accent"
              />
            </Field.Root>
            <div className="flex justify-end gap-3 pt-2">
              <Dialog.Close className="cursor-pointer rounded-md border border-border-secondary px-3 py-1.5 text-sm text-fg-secondary hover:bg-bg-tertiary">
                {t("admin.common.cancel")}
              </Dialog.Close>
              <button
                type="submit"
                disabled={loading}
                className="cursor-pointer rounded-md bg-accent px-3 py-1.5 text-sm font-medium text-accent-fg hover:opacity-90 disabled:opacity-50"
              >
                {loading ? t("admin.common.saving") : t("admin.common.save")}
              </button>
            </div>
          </form>
        </Dialog.Popup>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
