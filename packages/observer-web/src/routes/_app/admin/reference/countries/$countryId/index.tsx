import { type SyntheticEvent, useState } from "react";

import { ArrowLeftIcon, GlobeIcon } from "@/components/icons";
import { Field } from "@base-ui/react/field";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { ConfirmDialog } from "@/components/confirm-dialog";
import { DataTable, type Column } from "@/components/data-table";
import { EmptyState } from "@/components/empty-state";
import { FormDialog } from "@/components/form-dialog";
import { PageHeader } from "@/components/page-header";
import { RowActions } from "@/components/row-actions";
import { StatusBadge } from "@/components/status-badge";
import { useCreateState, useDeleteState, useStates, useUpdateState } from "@/hooks/use-states";
import type { State } from "@/types/reference";

export const Route = createFileRoute("/_app/admin/reference/countries/$countryId/")({
  component: StatesPage,
});

function StatesPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { countryId } = Route.useParams();
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
          <span className="text-fg-tertiary">{"—"}</span>
        ),
    },
    {
      key: "actions",
      header: "",
      render: (s) => (
        <RowActions onEdit={() => setEditTarget(s)} onDelete={() => setDeleteTarget(s)} />
      ),
    },
  ];

  return (
    <div>
      <Link
        to="/admin/reference/countries"
        className="mb-3 inline-flex items-center gap-1 text-sm text-fg-tertiary hover:text-fg"
      >
        <ArrowLeftIcon size={14} />
        {t("admin.reference.countries.title")}
      </Link>

      <PageHeader
        title={t("admin.reference.states.title")}
        action={
          <Button onClick={() => setCreateOpen(true)}>
            {t("admin.reference.states.add")}
          </Button>
        }
      />

      <DataTable
        columns={columns}
        data={data?.states ?? []}
        keyExtractor={(s) => s.id}
        onRowClick={(s) =>
          navigate({
            to: "/admin/reference/countries/$countryId/states/$stateId",
            params: { countryId, stateId: s.id },
          })
        }
        isLoading={isLoading}
        emptyState={
          <EmptyState
            icon={GlobeIcon}
            title={t("admin.reference.states.emptyTitle")}
            description={t("admin.reference.states.emptyDescription")}
            action={
              <Button onClick={() => setCreateOpen(true)}>
                {t("admin.reference.states.add")}
              </Button>
            }
          />
        }
      />

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
  const [conflictZone, setConflictZone] = useState(initial?.conflict_zone ?? "");

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    await onSubmit(name, code, conflictZone);
    if (!initial) {
      setName("");
      setCode("");
      setConflictZone("");
    }
  }

  return (
    <FormDialog
      open={open}
      onOpenChange={onOpenChange}
      title={title}
      loading={loading}
      onSubmit={handleSubmit}
    >
      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("admin.reference.states.name")}
        </Field.Label>
        <Field.Control
          required
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="block w-full rounded-lg border border-border-secondary bg-bg h-9 px-3 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg"
        />
      </Field.Root>
      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("admin.reference.states.code")}
        </Field.Label>
        <Field.Control
          value={code}
          onChange={(e) => setCode(e.target.value)}
          className="block w-full rounded-lg border border-border-secondary bg-bg h-9 px-3 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg"
        />
      </Field.Root>
      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("admin.reference.states.conflictZone")}
        </Field.Label>
        <Field.Control
          value={conflictZone}
          onChange={(e) => setConflictZone(e.target.value)}
          className="block w-full rounded-lg border border-border-secondary bg-bg h-9 px-3 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg"
        />
      </Field.Root>
    </FormDialog>
  );
}
