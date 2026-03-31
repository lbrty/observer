import { useState } from "react";

import { createLazyFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { DataTable, type Column } from "@/components/data-table";
import { EmptyState } from "@/components/empty-state";
import type { FilterDef } from "@/components/filter-bar";
import { FilterBar } from "@/components/filter-bar";
import { HouseholdDrawer } from "@/components/household-drawer";
import {
  DownloadSimpleIcon,
  HouseSimpleIcon,
  PencilSimpleIcon,
  PlusIcon,
} from "@/components/icons";
import { PageHeader } from "@/components/page-header";
import { Pagination } from "@/components/pagination";
import { useExportCSV } from "@/hooks/use-export-csv";
import { useHouseholds } from "@/hooks/use-households";
import { useProjectRole } from "@/hooks/use-project-role";
import type { Household } from "@/types/household";

export const Route = createLazyFileRoute("/_app/projects/$projectId/households/")({
  component: HouseholdsListPage,
});

function HouseholdsListPage() {
  const { t } = useTranslation();
  const { projectId } = Route.useParams();

  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [dateFrom, setDateFrom] = useState("");
  const [dateTo, setDateTo] = useState("");
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editHouseholdId, setEditHouseholdId] = useState<string | null>(null);
  const { canWrite, canExport } = useProjectRole(projectId);
  const { exporting, exportCSV } = useExportCSV();

  const params = {
    page,
    per_page: 20,
    ...(search && { search }),
    ...(dateFrom && { created_from: dateFrom }),
    ...(dateTo && { created_to: dateTo }),
  };

  const { data, isLoading } = useHouseholds(projectId, params);

  const filters: FilterDef[] = [
    {
      type: "search",
      placeholder: t("project.households.search"),
      value: search,
      onChange: (v) => {
        setSearch(v);
        setPage(1);
      },
    },
    {
      type: "date-range",
      fromValue: dateFrom,
      toValue: dateTo,
      onFromChange: (v) => {
        setDateFrom(v);
        setPage(1);
      },
      onToChange: (v) => {
        setDateTo(v);
        setPage(1);
      },
      fromPlaceholder: t("common.dateFrom"),
      toPlaceholder: t("common.dateTo"),
    },
  ];

  function openCreate() {
    setEditHouseholdId(null);
    setDrawerOpen(true);
  }

  function openEdit(householdId: string) {
    setEditHouseholdId(householdId);
    setDrawerOpen(true);
  }

  function handleExport() {
    const searchParams: Record<string, string> = {};
    if (search) searchParams.search = search;
    if (dateFrom) searchParams.created_from = dateFrom;
    if (dateTo) searchParams.created_to = dateTo;
    return exportCSV(`projects/${projectId}/export/households`, "households", searchParams);
  }

  const columns: Column<Household>[] = [
    {
      key: "head_person_id",
      header: t("project.households.headPerson"),
      render: (h) =>
        h.head_person_name ? (
          <span className="text-sm text-fg-secondary">{h.head_person_name}</span>
        ) : (
          <span className="text-fg-tertiary">—</span>
        ),
    },
    {
      key: "members",
      header: t("project.households.members"),
      render: (h) => <span className="text-fg-secondary">{h.member_count}</span>,
    },
    {
      key: "reference_number",
      header: t("project.households.referenceNumber"),
      render: (h) => (
        <div className="flex items-center gap-3">
          <span className="inline-flex size-8 shrink-0 items-center justify-center rounded-lg bg-bg-tertiary text-fg-tertiary">
            <HouseSimpleIcon size={16} />
          </span>
          <span className="font-medium text-fg">{h.reference_number || "-"}</span>
        </div>
      ),
    },
    {
      key: "created_at",
      header: t("project.households.createdAt"),
      render: (h) => (
        <span className="font-mono text-xs tabular-nums text-fg-tertiary">
          {new Date(h.created_at).toLocaleDateString("en-CA")}
        </span>
      ),
    },
    ...(canWrite
      ? [
          {
            key: "actions",
            header: "",
            render: (h: Household) => (
              <Button
                variant="ghost"
                className="p-1.5"
                onClick={(e) => {
                  e.stopPropagation();
                  openEdit(h.id);
                }}
              >
                <PencilSimpleIcon size={16} />
              </Button>
            ),
          } satisfies Column<Household>,
        ]
      : []),
  ];

  return (
    <div>
      <PageHeader
        title={t("project.households.title")}
        action={
          canWrite ? (
            <Button icon={<PlusIcon size={16} />} onClick={openCreate}>
              {t("project.households.create")}
            </Button>
          ) : undefined
        }
      />

      <FilterBar
        filters={filters}
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
        data={data?.households ?? []}
        keyExtractor={(h) => h.id}
        onRowClick={(h) => openEdit(h.id)}
        isLoading={isLoading}
        emptyState={
          <EmptyState
            icon={HouseSimpleIcon}
            title={t("project.households.emptyTitle")}
            description={t("project.households.emptyDescription")}
            action={
              canWrite ? (
                <Button onClick={openCreate} icon={<PlusIcon size={16} />}>
                  {t("project.households.create")}
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
          onChange={setPage}
        />
      )}

      <HouseholdDrawer
        open={drawerOpen}
        onOpenChange={setDrawerOpen}
        projectId={projectId}
        householdId={editHouseholdId}
      />
    </div>
  );
}
