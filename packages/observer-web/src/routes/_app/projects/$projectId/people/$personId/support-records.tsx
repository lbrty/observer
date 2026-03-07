import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { DataTable, type Column } from "@/components/data-table";
import { EmptyState } from "@/components/empty-state";
import { HandHeartIcon, PencilSimpleIcon, PlusIcon } from "@/components/icons";
import { Pagination } from "@/components/pagination";
import { StatusBadge } from "@/components/status-badge";
import { SupportRecordDrawer } from "@/components/support-record-drawer";
import { referralKeys, sphereKeys, typeKeys } from "@/constants/support";
import { useSupportRecords } from "@/hooks/use-support-records";
import type { SupportRecord } from "@/types/support-record";

export const Route = createFileRoute("/_app/projects/$projectId/people/$personId/support-records")({
  component: PersonSupportRecords,
});

function PersonSupportRecords() {
  const { t } = useTranslation();
  const { projectId, personId } = Route.useParams();

  const [page, setPage] = useState(1);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editRecordId, setEditRecordId] = useState<string | null>(null);

  const { data, isLoading } = useSupportRecords(projectId, {
    person_id: personId,
    page,
    per_page: 20,
  });

  function openCreate() {
    setEditRecordId(null);
    setDrawerOpen(true);
  }

  function openEdit(id: string) {
    setEditRecordId(id);
    setDrawerOpen(true);
  }

  const columns: Column<SupportRecord>[] = [
    {
      key: "type",
      header: t("project.supportRecords.type"),
      render: (r) => (
        <div className="flex items-center gap-3">
          <span className="inline-flex size-8 shrink-0 items-center justify-center rounded-lg bg-bg-tertiary text-fg-tertiary">
            <HandHeartIcon size={16} />
          </span>
          <StatusBadge label={typeKeys[r.type] ? t(typeKeys[r.type]) : r.type} />
        </div>
      ),
    },
    {
      key: "sphere",
      header: t("project.supportRecords.sphere"),
      render: (r) => <span className="text-fg-secondary">{r.sphere ? t(sphereKeys[r.sphere] ?? r.sphere) : "\u2014"}</span>,
    },
    {
      key: "provided_at",
      header: t("project.supportRecords.providedAt"),
      render: (r) => (
        <span className="font-mono text-xs tabular-nums text-fg-tertiary">
          {r.provided_at ? new Date(r.provided_at).toLocaleDateString("en-CA") : "\u2014"}
        </span>
      ),
    },
    {
      key: "referral_status",
      header: t("project.supportRecords.referralStatus"),
      render: (r) =>
        r.referral_status ? (
          <StatusBadge label={referralKeys[r.referral_status] ? t(referralKeys[r.referral_status]) : r.referral_status} />
        ) : (
          <span className="text-fg-tertiary">{"\u2014"}</span>
        ),
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
          <PencilSimpleIcon size={16} />
        </Button>
      ),
    },
  ];

  return (
    <div>
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-sm font-semibold text-fg-secondary">
          {t("project.supportRecords.title")}
        </h2>
        <Button icon={<PlusIcon size={16} />} onClick={openCreate}>
          {t("project.supportRecords.create")}
        </Button>
      </div>

      <DataTable
        columns={columns}
        data={data?.records ?? []}
        keyExtractor={(r) => r.id}
        onRowClick={(r) => openEdit(r.id)}
        isLoading={isLoading}
        emptyState={
          <EmptyState
            icon={HandHeartIcon}
            title={t("project.people.supportRecordsEmptyTitle")}
            description={t("project.people.supportRecordsEmptyDescription")}
          />
        }
      />

      {data && (
        <Pagination
          page={data.page}
          perPage={data.per_page}
          total={data.total}
          onChange={setPage}
        />
      )}

      <SupportRecordDrawer
        open={drawerOpen}
        onOpenChange={setDrawerOpen}
        projectId={projectId}
        recordId={editRecordId}
        personId={personId}
      />
    </div>
  );
}
