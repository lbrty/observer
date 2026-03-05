import { HandHeartIcon, PencilSimpleIcon } from "@/components/icons";
import { Tabs } from "@base-ui/react/tabs";
import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { DataTable, type Column } from "@/components/data-table";
import { PageHeader } from "@/components/page-header";
import { Pagination } from "@/components/pagination";
import { StatusBadge } from "@/components/status-badge";
import { SupportRecordDrawer } from "@/components/support-record-drawer";
import { useSupportRecords } from "@/hooks/use-support-records";
import type { SupportRecord } from "@/types/support-record";

export const Route = createFileRoute("/_app/projects/$projectId/support-records/")({
  component: SupportRecordsPage,
});

const typeTabs = [
  "",
  "humanitarian",
  "legal",
  "social",
  "psychological",
  "medical",
  "general",
] as const;

function SupportRecordsPage() {
  const { t } = useTranslation();
  const { projectId } = Route.useParams();

  const [page, setPage] = useState(1);
  const [typeFilter, setTypeFilter] = useState<string>("");
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editRecordId, setEditRecordId] = useState<string | null>(null);

  const params = {
    page,
    per_page: 20,
    ...(typeFilter && { type: typeFilter as SupportRecord["type"] }),
  };

  const { data, isLoading } = useSupportRecords(projectId, params);

  const tabLabels: Record<string, string> = {
    "": t("project.supportRecords.all"),
    humanitarian: t("project.supportRecords.typeHumanitarian"),
    legal: t("project.supportRecords.typeLegal"),
    social: t("project.supportRecords.typeSocial"),
    psychological: t("project.supportRecords.typePsychological"),
    medical: t("project.supportRecords.typeMedical"),
    general: t("project.supportRecords.typeGeneral"),
  };

  function openCreate() {
    setEditRecordId(null);
    setDrawerOpen(true);
  }

  function openEdit(recordId: string) {
    setEditRecordId(recordId);
    setDrawerOpen(true);
  }

  const columns: Column<SupportRecord>[] = [
    {
      key: "person_id",
      header: t("project.supportRecords.person"),
      render: (r) => (
        <div className="flex items-center gap-3">
          <span className="inline-flex size-8 shrink-0 items-center justify-center rounded-lg bg-bg-tertiary text-fg-tertiary">
            <HandHeartIcon size={16} />
          </span>
          <span className="truncate font-mono text-xs text-fg-secondary">{r.person_id}</span>
        </div>
      ),
    },
    {
      key: "type",
      header: t("project.supportRecords.type"),
      render: (r) => <StatusBadge label={r.type} />,
    },
    {
      key: "sphere",
      header: t("project.supportRecords.sphere"),
      render: (r) => <span className="text-fg-secondary">{r.sphere ?? "\u2014"}</span>,
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
          <StatusBadge label={r.referral_status} />
        ) : (
          <span className="text-fg-tertiary">{"\u2014"}</span>
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
          className="cursor-pointer rounded-lg p-1.5 text-fg-tertiary hover:bg-bg-tertiary hover:text-accent"
        >
          <PencilSimpleIcon size={16} />
        </button>
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title={t("project.supportRecords.title")}
        action={
          <button
            type="button"
            onClick={openCreate}
            className="cursor-pointer rounded-lg bg-accent px-4 py-2 text-sm font-medium text-accent-fg shadow-card hover:opacity-90"
          >
            + {t("project.supportRecords.create")}
          </button>
        }
      />

      <Tabs.Root
        defaultValue=""
        value={typeFilter}
        onValueChange={(value) => {
          setTypeFilter(value as string);
          setPage(1);
        }}
        className="mb-4"
      >
        <Tabs.List className="flex gap-0 rounded-lg border border-border-secondary bg-bg-secondary p-0.5">
          {typeTabs.map((tab) => (
            <Tabs.Tab
              key={tab}
              value={tab}
              className="cursor-pointer rounded-sm px-4 py-1.5 m-0.5 text-sm font-medium text-fg-tertiary transition-colors hover:text-fg data-active:bg-bg data-active:text-fg data-active:shadow-card"
            >
              {tabLabels[tab]}
            </Tabs.Tab>
          ))}
        </Tabs.List>
      </Tabs.Root>

      <DataTable
        columns={columns}
        data={data?.support_records ?? []}
        keyExtractor={(r) => r.id}
        onRowClick={(r) => openEdit(r.id)}
        isLoading={isLoading}
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
      />
    </div>
  );
}
