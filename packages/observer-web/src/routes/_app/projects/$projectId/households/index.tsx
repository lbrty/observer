import { HouseSimpleIcon, PencilSimpleIcon, PlusIcon } from "@/components/icons";
import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { DataTable, type Column } from "@/components/data-table";
import { EmptyState } from "@/components/empty-state";
import { HouseholdDrawer } from "@/components/household-drawer";
import { PageHeader } from "@/components/page-header";
import { Pagination } from "@/components/pagination";
import { useHouseholds } from "@/hooks/use-households";
import type { Household } from "@/types/household";

export const Route = createFileRoute("/_app/projects/$projectId/households/")({
  component: HouseholdsListPage,
});

function HouseholdsListPage() {
  const { t } = useTranslation();
  const { projectId } = Route.useParams();

  const [page, setPage] = useState(1);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editHouseholdId, setEditHouseholdId] = useState<string | null>(null);

  const params = {
    page,
    per_page: 20,
  };

  const { data, isLoading } = useHouseholds(projectId, params);

  function openCreate() {
    setEditHouseholdId(null);
    setDrawerOpen(true);
  }

  function openEdit(householdId: string) {
    setEditHouseholdId(householdId);
    setDrawerOpen(true);
  }

  const columns: Column<Household>[] = [
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
      key: "head_person_id",
      header: t("project.households.headPerson"),
      render: (h) => (
        <span className="font-mono text-xs text-fg-secondary">{h.head_person_id || "-"}</span>
      ),
    },
    {
      key: "members",
      header: t("project.households.members"),
      render: (h) => <span className="text-fg-secondary">{h.member_count}</span>,
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
    {
      key: "actions",
      header: "",
      render: (h) => (
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
    },
  ];

  return (
    <div>
      <PageHeader
        title={t("project.households.title")}
        action={
          <Button icon={<PlusIcon size={16} />} onClick={openCreate}>
            {t("project.households.create")}
          </Button>
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
              <Button onClick={openCreate} icon={<PlusIcon size={16} />}>
                {t("project.households.create")}
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

      <HouseholdDrawer
        open={drawerOpen}
        onOpenChange={setDrawerOpen}
        projectId={projectId}
        householdId={editHouseholdId}
      />
    </div>
  );
}
