import { useState } from "react";

import { Tabs } from "@base-ui/react/tabs";
import { Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { DataTable } from "@/components/data-table";
import { EmptyState } from "@/components/empty-state";
import { DownloadSimpleIcon, HandHeartIcon, PlusIcon } from "@/components/icons";
import { PageHeader } from "@/components/page-header";
import { Pagination } from "@/components/pagination";
import { buildSupportRecordColumns } from "@/components/support-record-columns";
import { SupportRecordFilterBar } from "@/components/support-record-filter-bar";
import { SupportRecordDrawer } from "@/components/support-record-drawer";
import { useExportCSV } from "@/hooks/use-export-csv";
import { useProjectRole } from "@/hooks/use-project-role";
import { useSupportRecords } from "@/hooks/use-support-records";
import type { SupportRecord } from "@/types/support-record";

const supportTypes = [
  "",
  "humanitarian",
  "legal",
  "social",
  "psychological",
  "medical",
  "general",
] as const;

export type SupportType = (typeof supportTypes)[number];

interface SupportRecordsContentProps {
  projectId: string;
  typeFilter: SupportType;
  page: number;
  onPageChange: (page: number) => void;
}

export function SupportRecordsContent({
  projectId,
  typeFilter,
  page,
  onPageChange,
}: SupportRecordsContentProps) {
  const { t } = useTranslation();

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editRecordId, setEditRecordId] = useState<string | null>(null);
  const [sphere, setSphere] = useState("");
  const [dateFrom, setDateFrom] = useState("");
  const [dateTo, setDateTo] = useState("");
  const { canWrite, canExport } = useProjectRole(projectId);
  const { exporting, exportCSV } = useExportCSV();

  const params = {
    page,
    per_page: 20,
    ...(typeFilter && { type: typeFilter as SupportRecord["type"] }),
    ...(sphere && { sphere: sphere as SupportRecord["sphere"] }),
    ...(dateFrom && { provided_from: dateFrom }),
    ...(dateTo && { provided_to: dateTo }),
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

  function handleExport() {
    const searchParams: Record<string, string> = {};
    if (typeFilter) searchParams.type = typeFilter;
    if (sphere) searchParams.sphere = sphere;
    if (dateFrom) searchParams.provided_from = dateFrom;
    if (dateTo) searchParams.provided_to = dateTo;
    return exportCSV(`projects/${projectId}/export/support-records`, "support-records", searchParams);
  }

  const columns = buildSupportRecordColumns({ t, canWrite, onEdit: openEdit });

  return (
    <div>
      <PageHeader
        title={t("project.supportRecords.title")}
        action={
          canWrite ? (
            <Button icon={<PlusIcon size={16} />} onClick={openCreate}>
              {t("project.supportRecords.create")}
            </Button>
          ) : undefined
        }
      />

      <Tabs.Root value={typeFilter} className="mb-4">
        <Tabs.List className="flex gap-0 rounded-lg border border-border-secondary bg-bg-secondary p-0.5">
          {supportTypes.map((tab) => (
            <Tabs.Tab
              key={tab}
              value={tab}
              nativeButton={false}
              render={
                <Link
                  to={
                    tab
                      ? "/projects/$projectId/support-records/$type"
                      : "/projects/$projectId/support-records"
                  }
                  params={tab ? { projectId, type: tab } : { projectId }}
                />
              }
              className="cursor-pointer rounded-sm px-4 py-1.5 m-0.5 text-sm font-medium text-fg-tertiary transition-colors hover:text-fg data-active:bg-bg data-active:text-fg data-active:shadow-card"
            >
              {tabLabels[tab]}
            </Tabs.Tab>
          ))}
        </Tabs.List>
      </Tabs.Root>

      <SupportRecordFilterBar
        sphere={sphere}
        onSphereChange={setSphere}
        dateFrom={dateFrom}
        dateTo={dateTo}
        onDateFromChange={setDateFrom}
        onDateToChange={setDateTo}
        trailing={
          canExport ? (
            <Button
              variant="secondary"
              icon={<DownloadSimpleIcon size={16} />}
              onClick={handleExport}
              disabled={exporting}
            >
              {t("common.export")}
            </Button>
          ) : undefined
        }
      />

      <DataTable
        columns={columns}
        data={data?.records ?? []}
        keyExtractor={(r) => r.id}
        onRowClick={(r) => openEdit(r.id)}
        isLoading={isLoading}
        emptyState={
          <EmptyState
            icon={HandHeartIcon}
            title={t("project.supportRecords.emptyTitle")}
            description={t("project.supportRecords.emptyDescription")}
            action={
              canWrite ? (
                <Button onClick={openCreate} icon={<PlusIcon size={16} />}>
                  {t("project.supportRecords.create")}
                </Button>
              ) : undefined
            }
          />
        }
      />

      {data && (
        <Pagination
          page={data.page}
          perPage={data.per_page}
          total={data.total}
          onChange={onPageChange}
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
