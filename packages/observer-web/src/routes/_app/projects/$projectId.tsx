import {
  ChartBarIcon,
  ClockCounterClockwiseIcon,
  FilesIcon,
  HandHeartIcon,
  HouseSimpleIcon,
  PawPrintIcon,
  SlidersHorizontalIcon,
  TagIcon,
  UserCircleIcon,
  UserFocusIcon,
} from "@/components/icons";
import { SidebarLink } from "@/components/sidebar-link";
import { useMyProjects } from "@/hooks/use-my-projects";
import { useAuth } from "@/stores/auth";
import { createFileRoute, Navigate, Outlet, useLocation } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

export const Route = createFileRoute("/_app/projects/$projectId")({
  component: ProjectLayout,
});

function ReportsGroup({ projectId }: { projectId: string }) {
  const { t } = useTranslation();
  const location = useLocation();
  const isActive = location.pathname.includes(`/projects/${projectId}/reports`);

  return (
    <div>
      <div
        className={`flex items-center gap-2.5 rounded-lg px-3 py-2 text-sm transition-colors ${
          isActive
            ? "font-medium text-accent"
            : "text-fg-secondary"
        }`}
      >
        <ChartBarIcon size={18} weight={isActive ? "fill" : "regular"} />
        {t("project.nav.reports")}
      </div>
      <div className="ml-7 space-y-0.5 border-l border-border-secondary pl-2">
        <SidebarLink
          to={`/projects/${projectId}/reports/people`}
          label={t("project.nav.reportsPeople")}
          icon={UserCircleIcon}
        />
        <SidebarLink
          to={`/projects/${projectId}/reports/pets`}
          label={t("project.nav.reportsPets")}
          icon={PawPrintIcon}
        />
        <SidebarLink
          to={`/projects/${projectId}/reports/custom`}
          label={t("project.nav.reportsCustom")}
          icon={SlidersHorizontalIcon}
        />
      </div>
    </div>
  );
}

function ProjectLayout() {
  const { t } = useTranslation();
  const { projectId } = Route.useParams();
  const { user } = useAuth();
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
          {project.can_view_documents && (
            <SidebarLink
              to={`/projects/${projectId}/documents`}
              label={t("project.nav.documents")}
              icon={FilesIcon}
            />
          )}
          <SidebarLink
            to={`/projects/${projectId}/pets`}
            label={t("project.nav.pets")}
            icon={PawPrintIcon}
          />
          <ReportsGroup projectId={projectId} />
          {(project.role === "owner" || project.role === "manager") && (
            <SidebarLink
              to={`/projects/${projectId}/audit-logs`}
              label={t("project.nav.auditLogs")}
              icon={ClockCounterClockwiseIcon}
            />
          )}
          {user?.role === "consultant" && (
            <SidebarLink
              to={`/projects/${projectId}/my-stats`}
              label={t("project.nav.myStats")}
              icon={UserFocusIcon}
            />
          )}
        </nav>
      </aside>
      <main className="min-w-0 flex-1 animate-fade-in px-8 py-6">
        <Outlet />
      </main>
    </div>
  );
}
