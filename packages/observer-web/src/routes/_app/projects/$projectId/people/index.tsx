import { MagnifyingGlassIcon, PencilSimpleIcon, PlusIcon, UserCircleIcon } from "@/components/icons";
import { Tabs } from "@base-ui/react/tabs";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { DataTable, type Column } from "@/components/data-table";
import { EmptyState } from "@/components/empty-state";
import { PageHeader } from "@/components/page-header";
import { Pagination } from "@/components/pagination";
import { PersonDrawer } from "@/components/person-drawer";
import { StatusBadge } from "@/components/status-badge";
import { TagChips } from "@/components/tag-chips";
import { TagFilter } from "@/components/tag-filter";
import { usePeople } from "@/hooks/use-people";
import type { Person } from "@/types/person";

export const Route = createFileRoute("/_app/projects/$projectId/people/")({
  component: PeopleListPage,
  validateSearch: (search: Record<string, unknown>): { status?: string; page?: number } => ({
    status: (search.status as string) || undefined,
    page: Number(search.page) || undefined,
  }),
});

const statusTabs = ["", "new", "active", "closed", "archived"] as const;

function PeopleListPage() {
  const { t } = useTranslation();
  const { projectId } = Route.useParams();
  const navigate = useNavigate();
  const { status = "", page = 1 } = Route.useSearch();

  const [search, setSearch] = useState("");
  const [tagIds, setTagIds] = useState<string[]>([]);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editPersonId, setEditPersonId] = useState<string | null>(null);

  function setStatus(value: string) {
    navigate({ from: Route.fullPath, search: { status: value || undefined }, replace: true });
  }

  function setPage(value: number) {
    navigate({ from: Route.fullPath, search: (prev) => ({ ...prev, page: value > 1 ? value : undefined }), replace: true });
  }

  const params = {
    page,
    per_page: 20,
    ...(status && { case_status: status }),
    ...(search && { search }),
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

  const columns: Column<Person>[] = [
    {
      key: "name",
      header: t("project.people.name"),
      render: (p) => (
        <div className="flex items-center gap-3">
          <span className="inline-flex size-8 shrink-0 items-center justify-center rounded-lg bg-bg-tertiary text-fg-tertiary">
            <UserCircleIcon size={16} />
          </span>
          <div className="min-w-0">
            <p className="truncate font-medium text-fg">
              {p.first_name}
              {p.last_name ? ` ${p.last_name}` : ""}
            </p>
          </div>
        </div>
      ),
    },
    {
      key: "sex",
      header: t("project.people.sex"),
      render: (p) => <span className="text-fg-secondary">{p.sex}</span>,
    },
    {
      key: "case_status",
      header: t("project.people.caseStatus"),
      render: (p) => <StatusBadge label={statusLabels[p.case_status] ?? p.case_status} />,
    },
    {
      key: "tags",
      header: t("project.tags.title"),
      render: (p) => <TagChips projectId={projectId} tagIds={p.tag_ids} />,
    },
    {
      key: "registered",
      header: t("project.people.registered"),
      render: (p) => (
        <span className="font-mono text-xs tabular-nums text-fg-tertiary">
          {new Date(p.registered_at ?? p.created_at).toLocaleDateString("en-CA")}
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
        title={t("project.people.title")}
        action={
          <Button icon={<PlusIcon size={16} />} onClick={openCreate}>
            {t("project.people.register")}
          </Button>
        }
      />

      <div className="mb-4 flex items-center gap-3">
        <div className="relative max-w-xs flex-1">
          <MagnifyingGlassIcon
            size={16}
            className="absolute top-1/2 left-3 -translate-y-1/2 text-fg-tertiary"
          />
          <input
            value={search}
            onChange={(e) => {
              setSearch(e.target.value);
            }}
            placeholder={t("project.people.search")}
            className="w-full rounded-lg border border-border-secondary bg-bg-secondary py-2 pr-3 pl-9 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg"
          />
        </div>
        <TagFilter projectId={projectId} selectedIds={tagIds} onChange={setTagIds} />
      </div>

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
              {tab === "" ? t("project.people.all") : statusLabels[tab]}
            </Tabs.Tab>
          ))}
        </Tabs.List>
      </Tabs.Root>

      <DataTable
        columns={columns}
        data={data?.people ?? []}
        keyExtractor={(p) => p.id}
        onRowClick={(p) =>
          navigate({
            to: "/projects/$projectId/people/$personId",
            params: { projectId, personId: p.id },
          })
        }
        isLoading={isLoading}
        emptyState={
          <EmptyState
            icon={UserCircleIcon}
            title={t("project.people.emptyTitle")}
            description={t("project.people.emptyDescription")}
            action={
              <Button onClick={openCreate} icon={<PlusIcon size={16} />}>
                {t("project.people.register")}
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

      <PersonDrawer
        open={drawerOpen}
        onOpenChange={setDrawerOpen}
        projectId={projectId}
        personId={editPersonId}
      />
    </div>
  );
}
