import { useState } from "react";

import { Tabs } from "@base-ui/react/tabs";
import { Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { DataTable, type Column } from "@/components/data-table";
import { EmptyState } from "@/components/empty-state";
import type { FilterDef } from "@/components/filter-bar";
import { FilterBar } from "@/components/filter-bar";
import { DownloadSimpleIcon, HandHeartIcon, PencilSimpleIcon, PlusIcon } from "@/components/icons";
import { PageHeader } from "@/components/page-header";
import { Pagination } from "@/components/pagination";
import { PersonName } from "@/components/person-name";
import { StatusBadge } from "@/components/status-badge";
import { SupportRecordDrawer } from "@/components/support-record-drawer";
import { referralKeys, sphereKeys, typeKeys } from "@/constants/support";
import { useMyProjects } from "@/hooks/use-my-projects";
import { useSupportRecords } from "@/hooks/use-support-records";
import { api } from "@/lib/api";
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

const sphereValues = [
  "housing_assistance",
  "document_recovery",
  "social_benefits",
  "property_rights",
  "employment_rights",
  "family_law",
  "healthcare_access",
  "education_access",
  "financial_aid",
  "psychological_support",
  "other",
] as const;

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
  const [exporting, setExporting] = useState(false);

  const { data: projectsData } = useMyProjects();
  const project = projectsData?.projects.find((p) => p.id === projectId);
  const canExport = project?.role === "owner" || project?.role === "manager";

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

  const sphereOptions = [
    { label: t("project.supportRecords.allSpheres"), value: "" },
    ...sphereValues.map((s) => ({
      label: sphereKeys[s] ? t(sphereKeys[s]) : s,
      value: s,
    })),
  ];

  const filters: FilterDef[] = [
    {
      type: "select",
      value: sphere,
      onValueChange: setSphere,
      options: sphereOptions,
      placeholder: t("project.supportRecords.allSpheres"),
    },
    {
      type: "date-range",
      fromValue: dateFrom,
      toValue: dateTo,
      onFromChange: setDateFrom,
      onToChange: setDateTo,
      fromPlaceholder: t("common.dateFrom"),
      toPlaceholder: t("common.dateTo"),
    },
  ];

  function openCreate() {
    setEditRecordId(null);
    setDrawerOpen(true);
  }

  function openEdit(recordId: string) {
    setEditRecordId(recordId);
    setDrawerOpen(true);
  }

  async function handleExport() {
    setExporting(true);
    try {
      const searchParams: Record<string, string> = {};
      if (typeFilter) searchParams.type = typeFilter;
      if (sphere) searchParams.sphere = sphere;
      if (dateFrom) searchParams.provided_from = dateFrom;
      if (dateTo) searchParams.provided_to = dateTo;

      const blob = await api.get(`projects/${projectId}/export/support-records`, { searchParams }).blob();
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      const date = new Date().toISOString().slice(0, 10);
      a.download = `support-records-${date}.csv`;
      a.click();
      URL.revokeObjectURL(url);
    } finally {
      setExporting(false);
    }
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
          <span className="truncate text-sm text-fg">
            <PersonName projectId={projectId} personId={r.person_id} />
          </span>
        </div>
      ),
    },
    {
      key: "type",
      header: t("project.supportRecords.type"),
      render: (r) => <StatusBadge label={typeKeys[r.type] ? t(typeKeys[r.type]) : r.type} />,
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
      <PageHeader
        title={t("project.supportRecords.title")}
        action={
          <Button icon={<PlusIcon size={16} />} onClick={openCreate}>
            {t("project.supportRecords.create")}
          </Button>
        }
      />

      <Tabs.Root value={typeFilter} className="mb-4">
        <Tabs.List className="flex gap-0 rounded-lg border border-border-secondary bg-bg-secondary p-0.5">
          {supportTypes.map((tab) => (
            <Tabs.Tab key={tab} value={tab} nativeButton={false} render={<Link
              to={
                tab
                  ? "/projects/$projectId/support-records/$type"
                  : "/projects/$projectId/support-records"
              }
              params={tab ? { projectId, type: tab } : { projectId }}
            />}
              className="cursor-pointer rounded-sm px-4 py-1.5 m-0.5 text-sm font-medium text-fg-tertiary transition-colors hover:text-fg data-active:bg-bg data-active:text-fg data-active:shadow-card"
            >
              {tabLabels[tab]}
            </Tabs.Tab>
          ))}
        </Tabs.List>
      </Tabs.Root>

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
              <Button onClick={openCreate} icon={<PlusIcon size={16} />}>
                {t("project.supportRecords.create")}
              </Button>
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
