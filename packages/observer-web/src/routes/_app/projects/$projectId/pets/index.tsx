import { PawPrintIcon, PencilSimpleIcon } from "@/components/icons";
import { Tabs } from "@base-ui/react/tabs";
import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { DataTable, type Column } from "@/components/data-table";
import { PageHeader } from "@/components/page-header";
import { Pagination } from "@/components/pagination";
import { PetDrawer } from "@/components/pet-drawer";
import { StatusBadge } from "@/components/status-badge";
import { usePets } from "@/hooks/use-pets";
import type { Pet } from "@/types/pet";

export const Route = createFileRoute("/_app/projects/$projectId/pets/")({
  component: PetsListPage,
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

  const [page, setPage] = useState(1);
  const [status, setStatus] = useState<string>("");
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editPetId, setEditPetId] = useState<string | null>(null);

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
      render: (p) => (
        <span className="text-fg-secondary">{p.owner_id ?? ""}</span>
      ),
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
        <button
          type="button"
          onClick={(e) => {
            e.stopPropagation();
            openEdit(p.id);
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
        title={t("project.pets.title")}
        action={
          <button
            type="button"
            onClick={openCreate}
            className="cursor-pointer rounded-lg bg-accent px-4 py-2 text-sm font-medium text-accent-fg shadow-card hover:opacity-90"
          >
            + {t("project.pets.register")}
          </button>
        }
      />

      <Tabs.Root
        defaultValue=""
        value={status}
        onValueChange={(value) => {
          setStatus(value as string);
          setPage(1);
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
