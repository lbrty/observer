import {
  ChartBarIcon,
  FilesIcon,
  HandHeartIcon,
  HouseSimpleIcon,
  PathIcon,
  PawPrintIcon,
  TagIcon,
  UserCircleIcon,
} from "@/components/icons";
import { SidebarLink } from "@/components/sidebar-link";
import { useMyProjects } from "@/hooks/use-my-projects";
import { createFileRoute, Navigate, Outlet } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

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
          <SidebarLink
            to={`/projects/${projectId}/support-records`}
            label={t("project.nav.supportRecords")}
            icon={HandHeartIcon}
          />
          <SidebarLink
            to={`/projects/${projectId}/households`}
            label={t("project.nav.households")}
            icon={HouseSimpleIcon}
          />
          <SidebarLink
            to={`/projects/${projectId}/tags`}
            label={t("project.nav.tags")}
            icon={TagIcon}
          />
          <SidebarLink
            to={`/projects/${projectId}/documents`}
            label={t("project.nav.documents")}
            icon={FilesIcon}
          />
          <SidebarLink
            to={`/projects/${projectId}/pets`}
            label={t("project.nav.pets")}
            icon={PawPrintIcon}
          />
          <SidebarLink
            to={`/projects/${projectId}/reports`}
            label={t("project.nav.reports")}
            icon={ChartBarIcon}
          />
        </nav>
      </aside>
      <main className="min-w-0 flex-1 px-8 py-6">
        <Outlet />
      </main>
    </div>
  );
}
