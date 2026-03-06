import { useState } from "react";

import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { DataTable, type Column } from "@/components/data-table";
import { EmptyState } from "@/components/empty-state";
import { PathIcon, PencilSimpleIcon, PlusIcon } from "@/components/icons";
import { JourneyTimeline } from "@/components/journey-timeline";
import { MigrationRecordDrawer } from "@/components/migration-record-drawer";
import { useMigrationRecords } from "@/hooks/use-migration-records";
import { usePlaces } from "@/hooks/use-places";
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
  const { data: placesData } = usePlaces();

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editId, setEditId] = useState<string | null>(null);
  const [view, setView] = useState<"timeline" | "table">("timeline");

  const records = data?.records ?? [];
  const places = placesData?.places ?? [];

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
      key: "from",
      header: t("project.migrationRecords.from"),
      render: (r) => {
        const name = r.from_place_id ? places.find((p) => p.id === r.from_place_id)?.name : null;
        return <span className="text-fg-secondary">{name ?? "—"}</span>;
      },
    },
    {
      key: "to",
      header: t("project.migrationRecords.to"),
      render: (r) => {
        const name = r.destination_place_id
          ? places.find((p) => p.id === r.destination_place_id)?.name
          : null;
        return <span className="text-fg-secondary">{name ?? "—"}</span>;
      },
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
      key: "actions",
      header: "",
      render: (r) => (
        <Button
          variant="ghost"
          className="p-1.5"
          onClick={(e) => {
            e.stopPropagation();
            openEdit(r.id);
          }}
        >
          <PencilSimpleIcon size={14} />
        </Button>
      ),
    },
  ];

  return (
    <div>
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-sm font-semibold text-fg">{t("project.migrationRecords.title")}</h2>
        <div className="flex items-center gap-2">
          {records.length > 0 && (
            <div className="flex rounded-lg border border-border-secondary bg-bg-secondary p-0.5">
              <button
                type="button"
                onClick={() => setView("timeline")}
                className={`cursor-pointer rounded-sm px-3 py-1 text-xs font-medium transition-colors ${
                  view === "timeline"
                    ? "bg-bg text-fg shadow-card"
                    : "text-fg-tertiary hover:text-fg"
                }`}
              >
                {t("project.migrationRecords.timelineView")}
              </button>
              <button
                type="button"
                onClick={() => setView("table")}
                className={`cursor-pointer rounded-sm px-3 py-1 text-xs font-medium transition-colors ${
                  view === "table"
                    ? "bg-bg text-fg shadow-card"
                    : "text-fg-tertiary hover:text-fg"
                }`}
              >
                {t("project.migrationRecords.tableView")}
              </button>
            </div>
          )}
          <Button size="sm" icon={<PlusIcon size={14} weight="bold" />} onClick={openCreate}>
            {t("admin.common.add")}
          </Button>
        </div>
      </div>

      {records.length === 0 && !isLoading && (
        <EmptyState
          icon={PathIcon}
          title={t("project.people.migrationRecordsEmptyTitle")}
          description={t("project.people.migrationRecordsEmptyDescription")}
        />
      )}

      {records.length > 0 && view === "timeline" && (
        <JourneyTimeline records={records} places={places} onEdit={openEdit} />
      )}

      {records.length > 0 && view === "table" && (
        <DataTable
          columns={columns}
          data={records}
          keyExtractor={(r) => r.id}
          isLoading={isLoading}
        />
      )}

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
