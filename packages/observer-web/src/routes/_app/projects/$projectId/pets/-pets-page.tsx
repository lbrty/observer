import { useState } from "react";

import { Tabs } from "@base-ui/react/tabs";
import { useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { DataTable, type Column } from "@/components/data-table";
import { EmptyState } from "@/components/empty-state";
import type { FilterDef } from "@/components/filter-bar";
import { FilterBar } from "@/components/filter-bar";
import { DownloadSimpleIcon, PawPrintIcon, PencilSimpleIcon, PlusIcon } from "@/components/icons";
import { PageHeader } from "@/components/page-header";
import { Pagination } from "@/components/pagination";
import { PersonName } from "@/components/person-name";
import { PetDrawer } from "@/components/pet-drawer";
import { StatusBadge } from "@/components/status-badge";
import { TagChips } from "@/components/tag-chips";
import { TagFilter } from "@/components/tag-filter";
import { useMyProjects } from "@/hooks/use-my-projects";
import { usePets } from "@/hooks/use-pets";
import { api } from "@/lib/api";
import type { Pet } from "@/types/pet";

export type PetStatus = "" | "registered" | "adopted" | "owner_found" | "needs_shelter" | "unknown";

const statusVariants: Record<string, "foam" | "gold" | "rose" | "neutral"> = {
  registered: "gold",
  adopted: "foam",
  owner_found: "foam",
  needs_shelter: "rose",
  unknown: "neutral",
};

const statusTabs: PetStatus[] = [
  "",
  "registered",
  "adopted",
  "owner_found",
  "needs_shelter",
  "unknown",
];

export function PetsContent({
  projectId,
  statusFilter,
  page,
  onPageChange,
}: {
  projectId: string;
  statusFilter: PetStatus;
  page: number;
  onPageChange: (page: number) => void;
}) {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editPetId, setEditPetId] = useState<string | null>(null);
  const [tagIds, setTagIds] = useState<string[]>([]);
  const [dateFrom, setDateFrom] = useState("");
  const [dateTo, setDateTo] = useState("");
  const [exporting, setExporting] = useState(false);

  const { data: projectsData } = useMyProjects();
  const project = projectsData?.projects.find((p) => p.id === projectId);
  const canExport = project?.role === "owner" || project?.role === "manager";

  const params = {
    page,
    per_page: 20,
    ...(statusFilter && { status: statusFilter }),
    ...(tagIds.length > 0 && { tag_ids: tagIds }),
    ...(dateFrom && { created_from: dateFrom }),
    ...(dateTo && { created_to: dateTo }),
  };

  const { data, isLoading } = usePets(projectId, params);

  const statusLabels: Record<string, string> = {
    registered: t("project.pets.statusRegistered"),
    adopted: t("project.pets.statusAdopted"),
    owner_found: t("project.pets.statusOwnerFound"),
    needs_shelter: t("project.pets.statusNeedsShelter"),
    unknown: t("project.pets.statusUnknown"),
  };

  const filters: FilterDef[] = [
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
    setEditPetId(null);
    setDrawerOpen(true);
  }

  function openEdit(petId: string) {
    setEditPetId(petId);
    setDrawerOpen(true);
  }

  async function handleExport() {
    setExporting(true);
    try {
      const searchParams: Record<string, string> = {};
      if (statusFilter) searchParams.status = statusFilter;
      if (dateFrom) searchParams.created_from = dateFrom;
      if (dateTo) searchParams.created_to = dateTo;

      const blob = await api.get(`projects/${projectId}/export/pets`, { searchParams }).blob();
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      const date = new Date().toISOString().slice(0, 10);
      a.download = `pets-${date}.csv`;
      a.click();
      URL.revokeObjectURL(url);
    } finally {
      setExporting(false);
    }
  }

  const columns: Column<Pet>[] = [
    {
      key: "name",
      header: t("project.pets.name"),
      render: (p) => (
        <div className="flex items-center gap-3">
          <span className="inline-flex size-8 shrink-0 items-center justify-center rounded-lg bg-bg-tertiary text-fg-tertiary">
            <PawPrintIcon size={16} />
          </span>
          <div className="min-w-0">
            <p className="truncate font-medium text-fg">{p.name}</p>
          </div>
        </div>
      ),
    },
    {
      key: "status",
      header: t("project.pets.status"),
      render: (p) => <StatusBadge label={statusLabels[p.status] ?? p.status} variant={statusVariants[p.status]} />,
    },
    {
      key: "owner_id",
      header: t("project.pets.ownerId"),
      render: (p) =>
        p.owner_id ? (
          <span className="text-sm text-fg-secondary">
            <PersonName projectId={projectId} personId={p.owner_id} />
          </span>
        ) : (
          <span className="text-fg-tertiary">—</span>
        ),
    },
    {
      key: "tags",
      header: t("project.tags.title"),
      render: (p) => <TagChips projectId={projectId} tagIds={p.tag_ids} />,
    },
    {
      key: "registration_id",
      header: t("project.pets.registrationId"),
      render: (p) => (
        <span className="font-mono text-xs tabular-nums text-fg-tertiary">
          {p.registration_id ?? ""}
        </span>
      ),
    },
    {
      key: "actions",
      header: "",
      render: (p) => (
        <Button
          variant="ghost"
          className="p-1.5"
          onClick={(e) => {
            e.stopPropagation();
            openEdit(p.id);
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
        title={t("project.pets.title")}
        action={
          <Button icon={<PlusIcon size={16} />} onClick={openCreate}>
            {t("project.pets.register")}
          </Button>
        }
      />

      <Tabs.Root
        defaultValue=""
        value={statusFilter}
        onValueChange={(value) => {
          const s = value as PetStatus;
          if (s) {
            navigate({ to: "/projects/$projectId/pets/$status", params: { projectId, status: s } });
          } else {
            navigate({ to: "/projects/$projectId/pets", params: { projectId } });
          }
        }}
        className="mb-4"
      >
        <Tabs.List className="flex gap-0 rounded-lg border border-border-secondary bg-bg-secondary p-0.5">
          {statusTabs.map((tab) => (
            <Tabs.Tab
              key={tab}
              value={tab}
              className="cursor-pointer rounded-sm px-4 py-1.5 m-0.5 text-sm font-medium text-fg-tertiary transition-colors hover:text-fg data-active:bg-bg data-active:text-fg data-active:shadow-card"
            >
              {tab === "" ? t("project.pets.all") : statusLabels[tab]}
            </Tabs.Tab>
          ))}
        </Tabs.List>
      </Tabs.Root>

      <FilterBar
        filters={filters}
        trailing={
          <div className="flex items-center gap-2">
            <TagFilter projectId={projectId} selectedIds={tagIds} onChange={setTagIds} />
            {canExport && (
              <Button
                variant="secondary"
                icon={<DownloadSimpleIcon size={16} />}
                onClick={handleExport}
                disabled={exporting}
              >
                {t("common.export")}
              </Button>
            )}
          </div>
        }
      />

      <DataTable
        columns={columns}
        data={data?.pets ?? []}
        keyExtractor={(p) => p.id}
        onRowClick={(p) => openEdit(p.id)}
        isLoading={isLoading}
        emptyState={
          <EmptyState
            icon={PawPrintIcon}
            title={t("project.pets.emptyTitle")}
            description={t("project.pets.emptyDescription")}
            action={
              <Button onClick={openCreate} icon={<PlusIcon size={16} />}>
                {t("project.pets.register")}
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

      <PetDrawer
        open={drawerOpen}
        onOpenChange={setDrawerOpen}
        projectId={projectId}
        petId={editPetId}
      />
    </div>
  );
}
