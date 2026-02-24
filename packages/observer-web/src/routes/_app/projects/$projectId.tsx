import { UserCircleIcon } from "@/components/icons";
import { createFileRoute, Navigate, Outlet } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { SidebarLink } from "@/components/sidebar-link";
import { useMyProjects } from "@/hooks/use-my-projects";

export const Route = createFileRoute("/_app/projects/$projectId")({
  component: ProjectLayout,
});

function ProjectLayout() {
  const { t } = useTranslation();
  const { projectId } = Route.useParams();
  const { data, isLoading } = useMyProjects();

  if (isLoading) return null;

  const project = data?.projects.find((p) => p.id === projectId);

  if (!project) {
    return <Navigate to="/" />;
  }

  return (
    <div className="flex flex-1">
      <aside className="w-52 shrink-0 border-r border-border-secondary">
        <nav className="sticky top-13 space-y-0.5 px-3 py-5">
          <div className="pb-1.5 pl-3 text-[11px] font-semibold uppercase tracking-wide text-fg-tertiary">
            {project.name}
          </div>
          <SidebarLink
            to={`/projects/${projectId}/people`}
            label={t("project.nav.people")}
            icon={UserCircleIcon}
          />
        </nav>
      </aside>
      <main className="min-w-0 flex-1 px-8 py-6">
        <Outlet />
      </main>
    </div>
  );
}
