import { useState } from "react";

import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import { DataTable } from "@/components/table/data-table";
import { EmptyState } from "@/components/ui/empty-state";
import { PawPrintIcon, PlusIcon } from "@/components/ui/icons";
import { PageHeader } from "@/components/layout/page-header";
import { Pagination } from "@/components/table/pagination";
import { buildPetColumns } from "@/components/pets/pet-columns";
import { PetFilterBar } from "@/components/pets/pet-filter-bar";
import { PetDrawer } from "@/components/pets/pet-drawer";
import { PetStatusTabs } from "@/components/pets/pet-status-tabs";
import { useExportCSV } from "@/hooks/use-export-csv";
import { useProjectRole } from "@/hooks/users/use-project-role";
import { usePets } from "@/hooks/pets/use-pets";
import type { Pet } from "@/types/pet";

export type PetStatus = "" | "registered" | "adopted" | "owner_found" | "needs_shelter" | "unknown";

const statusVariants: Record<string, "foam" | "gold" | "rose" | "neutral"> = {
  registered: "gold",
  adopted: "foam",
  owner_found: "foam",
  needs_shelter: "rose",
  unknown: "neutral",
};

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

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editPetId, setEditPetId] = useState<string | null>(null);
  const [tagIds, setTagIds] = useState<string[]>([]);
  const [dateFrom, setDateFrom] = useState("");
  const [dateTo, setDateTo] = useState("");
  const { canWrite, canExport } = useProjectRole(projectId);
  const { exporting, exportCSV } = useExportCSV();

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

  function openCreate() {
    setEditPetId(null);
    setDrawerOpen(true);
  }

  function openEdit(petId: string) {
    setEditPetId(petId);
    setDrawerOpen(true);
  }

  function handleExport() {
    const searchParams: Record<string, string> = {};
    if (statusFilter) searchParams.status = statusFilter;
    if (dateFrom) searchParams.created_from = dateFrom;
    if (dateTo) searchParams.created_to = dateTo;
    return exportCSV(`projects/${projectId}/export/pets`, "pets", searchParams);
  }

  const columns = buildPetColumns({
    t,
    projectId,
    statusLabels,
    statusVariants,
    canWrite,
    onEdit: openEdit,
  });

  return (
    <div>
      <PageHeader
        title={t("project.pets.title")}
        action={
          canWrite ? (
            <Button icon={<PlusIcon size={16} />} onClick={openCreate}>
              {t("project.pets.register")}
            </Button>
          ) : undefined
        }
      />

      <PetStatusTabs
        projectId={projectId}
        statusFilter={statusFilter}
        statusLabels={statusLabels}
      />

      <PetFilterBar
        projectId={projectId}
        dateFrom={dateFrom}
        dateTo={dateTo}
        tagIds={tagIds}
        canExport={canExport}
        exporting={exporting}
        onDateFromChange={setDateFrom}
        onDateToChange={setDateTo}
        onTagsChange={setTagIds}
        onExport={handleExport}
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
              canWrite ? (
                <Button onClick={openCreate} icon={<PlusIcon size={16} />}>
                  {t("project.pets.register")}
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

      <PetDrawer
        open={drawerOpen}
        onOpenChange={setDrawerOpen}
        projectId={projectId}
        petId={editPetId}
      />
    </div>
  );
}
