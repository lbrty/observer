import { useState } from "react";

import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { DataTable, type Column } from "@/components/data-table";
import { EmptyState } from "@/components/empty-state";
import { PathIcon, PencilSimpleIcon, PlusIcon } from "@/components/icons";
import { MigrationRecordDrawer } from "@/components/migration-record-drawer";
import { useMigrationRecords } from "@/hooks/use-migration-records";
import type { MigrationRecord } from "@/types/migration-record";

export const Route = createFileRoute(
  "/_app/projects/$projectId/people/$personId/migration-records",
)({
  component: PersonMigrationRecords,
});

function PersonMigrationRecords() {
  const { t } = useTranslation();
  const { projectId, personId } = Route.useParams();
  const { data, isLoading } = useMigrationRecords(projectId, personId);

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editId, setEditId] = useState<string | null>(null);

  function openCreate() {
    setEditId(null);
    setDrawerOpen(true);
  }

  function openEdit(id: string) {
    setEditId(id);
    setDrawerOpen(true);
  }

  const columns: Column<MigrationRecord>[] = [
    {
      key: "migration_date",
      header: t("project.migrationRecords.date"),
      render: (r) => (
        <span className="font-mono text-xs tabular-nums text-fg-tertiary">
          {r.migration_date ? new Date(r.migration_date).toLocaleDateString("en-CA") : "—"}
        </span>
      ),
    },
    {
      key: "movement_reason",
      header: t("project.migrationRecords.reason"),
      render: (r) => <span className="text-fg-secondary">{r.movement_reason ?? "—"}</span>,
    },
    {
      key: "housing",
      header: t("project.migrationRecords.housing"),
      render: (r) => <span className="text-fg-secondary">{r.housing_at_destination ?? "—"}</span>,
    },
    {
      key: "created_at",
      header: t("project.migrationRecords.created"),
      render: (r) => (
        <span className="font-mono text-xs tabular-nums text-fg-tertiary">
          {new Date(r.created_at).toLocaleDateString("en-CA")}
        </span>
      ),
    },
    {
      key: "actions",
      header: "",
      render: (r) => (
        <button
          type="button"
          onClick={(e) => {
            e.stopPropagation();
            openEdit(r.id);
          }}
          className="cursor-pointer rounded-lg p-1.5 text-fg-tertiary hover:bg-bg-tertiary hover:text-fg"
        >
          <PencilSimpleIcon size={14} />
        </button>
      ),
    },
  ];

  return (
    <div>
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-sm font-semibold text-fg">{t("project.migrationRecords.title")}</h2>
        <button
          type="button"
          onClick={openCreate}
          className="inline-flex cursor-pointer items-center gap-1.5 rounded-lg bg-accent px-3 py-1.5 text-sm font-medium text-accent-fg shadow-card hover:opacity-90"
        >
          <PlusIcon size={14} weight="bold" />
          {t("admin.common.add")}
        </button>
      </div>

      <DataTable
        columns={columns}
        data={data?.records ?? []}
        keyExtractor={(r) => r.id}
        isLoading={isLoading}
        emptyState={
          <EmptyState
            icon={PathIcon}
            title={t("project.people.migrationRecordsEmptyTitle")}
            description={t("project.people.migrationRecordsEmptyDescription")}
          />
        }
      />

      <MigrationRecordDrawer
        open={drawerOpen}
        onOpenChange={setDrawerOpen}
        projectId={projectId}
        personId={personId}
        recordId={editId}
      />
    </div>
  );
}
