import { Link, Outlet } from "@tanstack/react-router";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { ArrowLeftIcon } from "@/components/icons";
import { StatusBadge } from "@/components/status-badge";
import { usePerson } from "@/hooks/use-people";

export const Route = createFileRoute("/_app/projects/$projectId/people/$personId")({
  component: PersonDetailLayout,
});

function PersonDetailLayout() {
  const { t } = useTranslation();
  const { projectId, personId } = Route.useParams();
  const { data: person, isLoading } = usePerson(projectId, personId);

  const tabs = [
    {
      to: "/projects/$projectId/people/$personId" as const,
      label: t("project.people.overview"),
      exact: true,
    },
    {
      to: "/projects/$projectId/people/$personId/support-records" as const,
      label: t("project.people.supportTab"),
    },
    {
      to: "/projects/$projectId/people/$personId/notes" as const,
      label: t("project.people.notesTab"),
    },
    {
      to: "/projects/$projectId/people/$personId/migration-records" as const,
      label: t("project.people.migrationRecordsTab"),
    },
    {
      to: "/projects/$projectId/people/$personId/documents" as const,
      label: t("project.people.documentsTab"),
    },
  ];

  if (isLoading) {
    return (
      <div className="space-y-4">
        <div className="h-6 w-48 animate-pulse rounded bg-bg-tertiary" />
        <div className="h-10 w-72 animate-pulse rounded bg-bg-tertiary" />
      </div>
    );
  }

  if (!person) return null;

  const fullName = [person.first_name, person.last_name].filter(Boolean).join(" ");

  return (
    <div className="page-bg-people">
      <Link
        to="/projects/$projectId/people"
        params={{ projectId }}
        className="mb-4 inline-flex items-center gap-1.5 text-sm text-fg-tertiary transition-colors hover:text-fg"
      >
        <ArrowLeftIcon size={14} />
        {t("project.people.backToPeople")}
      </Link>

      <div className="mb-6 flex items-center gap-3">
        <h1 className="font-serif text-xl font-bold tracking-tight text-fg">{fullName}</h1>
        <StatusBadge label={person.case_status} />
      </div>

      <div className="mb-6 flex gap-0 rounded-lg border border-border-secondary bg-bg-secondary p-0.5">
        {tabs.map((tab) => (
          <Link
            key={tab.to}
            to={tab.to}
            params={{ projectId, personId }}
            activeOptions={{ exact: tab.exact ?? false }}
            className="rounded-sm px-4 py-1.5 m-0.5 text-sm font-medium transition-colors"
            activeProps={{ className: "bg-bg text-fg shadow-card" }}
            inactiveProps={{ className: "text-fg-tertiary hover:text-fg" }}
          >
            {tab.label}
          </Link>
        ))}
      </div>

      <Outlet />
    </div>
  );
}
