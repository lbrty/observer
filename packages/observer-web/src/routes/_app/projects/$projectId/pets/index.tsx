import { PawPrintIcon, PencilSimpleIcon, PlusIcon } from "@/components/icons";
import { Tabs } from "@base-ui/react/tabs";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { DataTable, type Column } from "@/components/data-table";
import { EmptyState } from "@/components/empty-state";
import { PageHeader } from "@/components/page-header";
import { Pagination } from "@/components/pagination";
import { PetDrawer } from "@/components/pet-drawer";
import { StatusBadge } from "@/components/status-badge";
import { usePets } from "@/hooks/use-pets";
import type { Pet } from "@/types/pet";

export const Route = createFileRoute("/_app/projects/$projectId/pets/")({
  component: PetsListPage,
  validateSearch: (search: Record<string, unknown>): { status?: string; page?: number } => ({
    status: (search.status as string) || undefined,
    page: Number(search.page) || undefined,
  }),
});

const statusTabs = [
  "",
  "registered",
  "adopted",
  "owner_found",
  "needs_shelter",
  "unknown",
] as const;

function PetsListPage() {
  const { t } = useTranslation();
  const { projectId } = Route.useParams();
  const navigate = useNavigate();
  const { status = "", page = 1 } = Route.useSearch();

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editPetId, setEditPetId] = useState<string | null>(null);

  function setStatus(value: string) {
    navigate({ from: Route.fullPath, search: { status: value || undefined }, replace: true });
  }

  function setPage(value: number) {
    navigate({ from: Route.fullPath, search: (prev) => ({ ...prev, page: value > 1 ? value : undefined }), replace: true });
  }

  const params = {
    page,
    per_page: 20,
    ...(status && { status }),
  };

  const { data, isLoading } = usePets(projectId, params);

  const tabLabels: Record<string, string> = {
    "": t("project.pets.all"),
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
      render: (p) => <StatusBadge label={p.status} />,
    },
    {
      key: "owner_id",
      header: t("project.pets.ownerId"),
      render: (p) => <span className="text-fg-secondary">{p.owner_id ?? ""}</span>,
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
        value={status}
        onValueChange={(value) => {
          setStatus(value as string);
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
              {tabLabels[tab]}
            </Tabs.Tab>
          ))}
        </Tabs.List>
      </Tabs.Root>

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
          onChange={setPage}
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
