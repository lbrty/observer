import {
  MagnifyingGlassIcon,
  PencilSimpleIcon,
  UserCircleIcon,
} from "@/components/icons";
import { Tabs } from "@base-ui/react/tabs";
import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { DataTable, type Column } from "@/components/data-table";
import { PageHeader } from "@/components/page-header";
import { Pagination } from "@/components/pagination";
import { PersonDrawer } from "@/components/person-drawer";
import { StatusBadge } from "@/components/status-badge";
import { usePeople } from "@/hooks/use-people";
import type { Person } from "@/types/person";

export const Route = createFileRoute("/_app/projects/$projectId/people/")({
  component: PeopleListPage,
});

const statusTabs = ["", "new", "active", "closed", "archived"] as const;

function PeopleListPage() {
  const { t } = useTranslation();
  const { projectId } = Route.useParams();

  const [page, setPage] = useState(1);
  const [status, setStatus] = useState<string>("");
  const [search, setSearch] = useState("");
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editPersonId, setEditPersonId] = useState<string | null>(null);

  const params = {
    page,
    per_page: 20,
    ...(status && { case_status: status }),
    ...(search && { search }),
  };

  const { data, isLoading } = usePeople(projectId, params);

  const tabLabels: Record<string, string> = {
    "": t("project.people.all"),
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
      render: (p) => <StatusBadge label={p.case_status} />,
    },
    {
      key: "registered",
      header: t("project.people.registered"),
      render: (p) => (
        <span className="font-mono text-xs tabular-nums text-fg-tertiary">
          {new Date(p.registered_at ?? p.created_at).toLocaleDateString(
            "en-CA",
          )}
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
        title={t("project.people.title")}
        action={
          <button
            type="button"
            onClick={openCreate}
            className="cursor-pointer rounded-lg bg-accent px-4 py-2 text-sm font-medium text-accent-fg shadow-card hover:opacity-90"
          >
            + {t("project.people.register")}
          </button>
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
              setPage(1);
            }}
            placeholder={t("project.people.search")}
            className="w-full rounded-lg border border-border-secondary bg-bg-secondary py-2 pr-3 pl-9 text-sm text-fg outline-none focus:border-accent"
          />
        </div>
      </div>

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
        data={data?.people ?? []}
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

      <PersonDrawer
        open={drawerOpen}
        onOpenChange={setDrawerOpen}
        projectId={projectId}
        personId={editPersonId}
      />
    </div>
  );
}
