import { useState } from "react";

import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import { DataTable } from "@/components/table/data-table";
import { EmptyState } from "@/components/ui/empty-state";
import { DownloadSimpleIcon, HandHeartIcon, PlusIcon } from "@/components/ui/icons";
import { PageHeader } from "@/components/layout/page-header";
import { Pagination } from "@/components/table/pagination";
import { buildSupportRecordColumns } from "@/components/support/support-record-columns";
import { SupportRecordDrawer } from "@/components/support/support-record-drawer";
import { SupportRecordFilterBar } from "@/components/support/support-record-filter-bar";
import { SupportRecordStatusTabs } from "@/components/support/support-record-status-tabs";
import { useExportCSV } from "@/hooks/use-export-csv";
import { useProjectRole } from "@/hooks/users/use-project-role";
import { useSupportRecords } from "@/hooks/support/use-support-records";
import type { SupportRecord } from "@/types/support-record";

export type SupportType =
  | ""
  | "humanitarian"
  | "legal"
  | "social"
  | "psychological"
  | "medical"
  | "general";

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
    return exportCSV(
      `projects/${projectId}/export/support-records`,
      "support-records",
      searchParams,
    );
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

      <SupportRecordStatusTabs projectId={projectId} typeFilter={typeFilter} />

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
