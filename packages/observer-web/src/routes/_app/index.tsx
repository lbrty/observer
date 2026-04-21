import { BuildingsIcon, FolderSimpleIcon, GlobeIcon, UsersIcon } from "@/components/ui/icons";
import type { Icon } from "@/components/ui/icons";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { StatusBadge } from "@/components/ui/status-badge";
import { useCountries } from "@/hooks/reference/use-countries";
import { useOffices } from "@/hooks/reference/use-offices";
import { useMyProjects } from "@/hooks/projects/use-my-projects";
import { useProjects } from "@/hooks/projects/use-projects";
import { useUsers } from "@/hooks/users/use-users";
import { useAuth } from "@/stores/auth";

export const Route = createFileRoute("/_app/")({
  component: DashboardPage,
});

const colorClasses = {
  accent: "bg-accent/10 text-accent",
  foam: "bg-foam/10 text-foam",
  gold: "bg-gold/10 text-gold",
  rose: "bg-rose/10 text-rose",
};

interface StatAction {
  key: string;
  to: string;
  icon: Icon;
  color: keyof typeof colorClasses;
  labelKey: string;
  value: number | string | null;
}

function DashboardPage() {
  const { t } = useTranslation();
  const { user } = useAuth();

  const isAdminOrStaff = user?.role === "admin" || user?.role === "staff";

  const { data: projectsData } = useProjects({ per_page: 1 }, { enabled: isAdminOrStaff });
  const { data: usersData } = useUsers({ per_page: 1 }, { enabled: isAdminOrStaff });
  const { data: countriesData } = useCountries();
  const { data: officesData } = useOffices();
  const { data: myProjectsData } = useMyProjects();

  const myProjects = myProjectsData?.projects ?? [];

  const statActions: StatAction[] = [
    {
      key: "projects",
      to: "/admin/projects",
      icon: FolderSimpleIcon,
      color: "foam",
      labelKey: "dashboard.projects",
      value: projectsData?.total ?? "—",
    },
    {
      key: "users",
      to: "/admin/users",
      icon: UsersIcon,
      color: "accent",
      labelKey: "dashboard.users",
      value: usersData?.total ?? "—",
    },
    {
      key: "countries",
      to: "/admin/reference/countries",
      icon: GlobeIcon,
      color: "gold",
      labelKey: "dashboard.countries",
      value: countriesData?.length ?? "—",
    },
    {
      key: "offices",
      to: "/admin/reference/offices",
      icon: BuildingsIcon,
      color: "rose",
      labelKey: "dashboard.offices",
      value: officesData?.length ?? "—",
    },
  ];

  return (
    <div className="mx-auto w-full max-w-270 px-10 py-8">
      <div className="pb-8">
        <h1 className="font-serif text-3xl font-bold tracking-tight text-fg">
          {t("dashboard.greeting")}
        </h1>
        <div className="mt-2">
          <StatusBadge label={user?.role ?? ""} />
        </div>
      </div>

      {isAdminOrStaff && (
        <div className="grid grid-cols-2 gap-3 sm:grid-cols-4">
          {statActions.map(({ key, to, icon: ActionIcon, color, labelKey, value }) => (
            <Link
              key={key}
              to={to}
              className="card-bg-topo group rounded-xl border border-border-secondary bg-bg-secondary p-5 transition-shadow hover:shadow-elevated"
            >
              <span
                className={`relative mb-3 inline-flex size-9 items-center justify-center rounded-xl ${colorClasses[color]}`}
              >
                <ActionIcon size={18} weight="duotone" />
              </span>
              <p className="relative text-[11px] font-semibold uppercase tracking-wide text-fg-tertiary">
                {t(labelKey)}
              </p>
              {value !== null && (
                <p className="relative mt-1 font-mono text-2xl font-bold tabular-nums text-fg">
                  {value}
                </p>
              )}
            </Link>
          ))}
        </div>
      )}

      <h2 className="mt-8 mb-4 font-serif text-lg font-semibold text-fg">
        {t("dashboard.myProjects")}
      </h2>

      {myProjects.length === 0 ? (
        <p className="text-sm text-fg-tertiary">{t("dashboard.noProjects")}</p>
      ) : (
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {myProjects.map((project) => (
            <Link
              key={project.id}
              to="/projects/$projectId/people"
              params={{ projectId: project.id }}
              className="card-bg-waves group rounded-xl border border-border-secondary bg-bg-secondary p-5 transition-shadow hover:shadow-elevated"
            >
              <span className="relative mb-4 inline-flex size-10 items-center justify-center rounded-xl bg-foam/10 text-foam">
                <FolderSimpleIcon size={20} weight="duotone" />
              </span>
              <p className="relative text-sm font-medium text-fg">{project.name}</p>
              {project.description && (
                <p className="relative mt-0.5 truncate text-xs text-fg-tertiary">
                  {project.description}
                </p>
              )}
              <div className="relative mt-3">
                <StatusBadge label={project.role} />
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
