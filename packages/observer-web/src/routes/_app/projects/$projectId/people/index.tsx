import { useState } from "react";

import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import { PlusIcon } from "@/components/ui/icons";
import { PageHeader } from "@/components/layout/page-header";
import { buildPeopleColumns } from "@/components/people/people-columns";
import { PeopleFilterBar } from "@/components/people/people-filter-bar";
import { PeopleStatusTabs } from "@/components/people/people-status-tabs";
import { PeopleTable } from "@/components/people/people-table";
import { PersonDrawer } from "@/components/people/person-drawer";
import { useExportCSV } from "@/hooks/use-export-csv";
import { useProjectRole } from "@/hooks/users/use-project-role";
import { usePeople } from "@/hooks/people/use-people";

export const Route = createFileRoute("/_app/projects/$projectId/people/")({
  component: PeopleListPage,
  validateSearch: (search: Record<string, unknown>): { status?: string; page?: number } => ({
    status: (search.status as string) || undefined,
    page: Number(search.page) || undefined,
  }),
});

function PeopleListPage() {
  const { t } = useTranslation();
  const { projectId } = Route.useParams();
  const navigate = useNavigate();
  const { status = "", page = 1 } = Route.useSearch();

  const [search, setSearch] = useState("");
  const [sex, setSex] = useState("");
  const [ageGroup, setAgeGroup] = useState("");
  const [dateFrom, setDateFrom] = useState("");
  const [dateTo, setDateTo] = useState("");
  const [tagIds, setTagIds] = useState<string[]>([]);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editPersonId, setEditPersonId] = useState<string | null>(null);
  const { canWrite, canExport } = useProjectRole(projectId);
  const { exporting, exportCSV } = useExportCSV();

  function setStatus(value: string) {
    navigate({ from: Route.fullPath, search: { status: value || undefined }, replace: true });
  }

  function setPage(value: number) {
    navigate({
      from: Route.fullPath,
      search: (prev) => ({ ...prev, page: value > 1 ? value : undefined }),
      replace: true,
    });
  }

  const params = {
    page,
    per_page: 20,
    ...(status && { case_status: status }),
    ...(search && { search }),
    ...(sex && { sex }),
    ...(ageGroup && { age_group: ageGroup }),
    ...(dateFrom && { registered_from: dateFrom }),
    ...(dateTo && { registered_to: dateTo }),
    ...(tagIds.length > 0 && { tag_ids: tagIds }),
  };

  const { data, isLoading } = usePeople(projectId, params);

  const statusLabels: Record<string, string> = {
    new: t("project.people.new"),
    active: t("project.people.active"),
    closed: t("project.people.closed"),
    archived: t("project.people.archived"),
  };

  function openCreate() {
    setEditPersonId(null);
    setDrawerOpen(true);
  }

  function openEdit(personId: string) {
    setEditPersonId(personId);
    setDrawerOpen(true);
  }

  function handleExport() {
    const searchParams: Record<string, string> = {};
    if (status) searchParams.case_status = status;
    if (search) searchParams.search = search;
    if (sex) searchParams.sex = sex;
    if (ageGroup) searchParams.age_group = ageGroup;
    if (dateFrom) searchParams.registered_from = dateFrom;
    if (dateTo) searchParams.registered_to = dateTo;
    return exportCSV(`projects/${projectId}/export/people`, "people", searchParams);
  }

  const columns = buildPeopleColumns({
    projectId,
    t,
    onEdit: openEdit,
    canWrite,
    statusLabels,
  });

  return (
    <div>
      <PageHeader
        title={t("project.people.title")}
        action={
          canWrite ? (
            <Button icon={<PlusIcon size={16} />} onClick={openCreate}>
              {t("project.people.register")}
            </Button>
          ) : undefined
        }
      />

      <PeopleFilterBar
        projectId={projectId}
        search={search}
        onSearchChange={setSearch}
        sex={sex}
        onSexChange={setSex}
        ageGroup={ageGroup}
        onAgeGroupChange={setAgeGroup}
        dateFrom={dateFrom}
        onDateFromChange={setDateFrom}
        dateTo={dateTo}
        onDateToChange={setDateTo}
        tagIds={tagIds}
        onTagIdsChange={setTagIds}
        canExport={canExport}
        exporting={exporting}
        onExport={handleExport}
      />

      <PeopleStatusTabs
        tabs={["", "new", "active", "closed", "archived"]}
        value={status}
        onValueChange={(v) => setStatus(v as string)}
        getLabel={(tab) => (tab === "" ? t("project.people.all") : statusLabels[tab])}
      />

      <PeopleTable
        columns={columns}
        data={data?.people ?? []}
        isLoading={isLoading}
        onRowClick={(p) =>
          navigate({
            to: "/projects/$projectId/people/$personId",
            params: { projectId, personId: p.id },
          })
        }
        canWrite={canWrite}
        onCreateClick={openCreate}
        t={t}
        pagination={
          data
            ? { page: data.page, perPage: data.per_page, total: data.total, onChange: setPage }
            : undefined
        }
      />

      <PersonDrawer
        open={drawerOpen}
        onOpenChange={setDrawerOpen}
        projectId={projectId}
        personId={editPersonId}
      />
    </div>
  );
}
