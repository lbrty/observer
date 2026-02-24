import {
  FolderSimpleIcon,
  GlobeIcon,
  TagIcon,
  UsersIcon,
} from "@/components/icons";
import type { Icon } from "@/components/icons";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { StatusBadge } from "@/components/status-badge";
import { useCountries } from "@/hooks/use-countries";
import { useMyProjects } from "@/hooks/use-my-projects";
import { useProjects } from "@/hooks/use-projects";
import { useUsers } from "@/hooks/use-users";
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

interface QuickAction {
  to: string;
  icon: Icon;
  color: keyof typeof colorClasses;
  titleKey: string;
  descKey: string;
}

const quickActions: QuickAction[] = [
  {
    to: "/admin/users",
    icon: UsersIcon,
    color: "accent",
    titleKey: "dashboard.manageUsers",
    descKey: "dashboard.manageUsersDesc",
  },
  {
    to: "/admin/projects",
    icon: FolderSimpleIcon,
    color: "foam",
    titleKey: "dashboard.projectsAction",
    descKey: "dashboard.projectsActionDesc",
  },
  {
    to: "/admin/reference/countries",
    icon: GlobeIcon,
    color: "gold",
    titleKey: "dashboard.referenceData",
    descKey: "dashboard.referenceDataDesc",
  },
  {
    to: "/admin/reference/categories",
    icon: TagIcon,
    color: "rose",
    titleKey: "dashboard.categories",
    descKey: "dashboard.categoriesDesc",
  },
];

function DashboardPage() {
  const { t } = useTranslation();
  const { user } = useAuth();

  const isAdminOrStaff = user?.role === "admin" || user?.role === "staff";

  const { data: projectsData } = useProjects({ per_page: 1 });
  const { data: usersData } = useUsers({ per_page: 1 });
  const { data: activeUsersData } = useUsers({ per_page: 1, is_active: true });
  const { data: countriesData } = useCountries();
  const { data: myProjectsData } = useMyProjects();

  const stats = [
    { label: t("dashboard.projects"), value: projectsData?.total ?? "—" },
    { label: t("dashboard.users"), value: usersData?.total ?? "—" },
    { label: t("dashboard.active"), value: activeUsersData?.total ?? "—" },
    {
      label: t("dashboard.countries"),
      value: countriesData?.length ?? "—",
    },
  ];

  const myProjects = myProjectsData?.projects ?? [];

  return (
    <div className="mx-auto w-full max-w-[1080px] px-10 py-8">
      <div className="pb-8">
        <h1 className="font-serif text-3xl font-bold tracking-tight text-fg">
          {t("dashboard.greeting")}
        </h1>
        <p className="mt-2 text-sm text-fg-secondary">
          {t("dashboard.role")} <StatusBadge label={user?.role ?? ""} />
        </p>
      </div>

      {isAdminOrStaff && (
        <>
          <div className="grid grid-cols-2 gap-3 sm:grid-cols-4">
            {stats.map((stat) => (
              <div
                key={stat.label}
                className="card-bg-topo rounded-xl border border-border-secondary bg-bg-secondary p-5"
              >
                <p className="relative text-[11px] font-semibold uppercase tracking-wide text-fg-tertiary">
                  {stat.label}
                </p>
                <p className="relative mt-2 font-mono text-2xl font-bold tabular-nums text-fg">
                  {stat.value}
                </p>
              </div>
            ))}
          </div>

          <h2 className="mt-8 mb-4 font-serif text-lg font-semibold text-fg">
            {t("dashboard.quickActions")}
          </h2>

          <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
            {quickActions.map(
              ({ to, icon: ActionIcon, color, titleKey, descKey }) => (
                <Link
                  key={to}
                  to={to}
                  className="card-bg-dots group rounded-xl border border-border-secondary bg-bg-secondary p-5 transition-shadow hover:shadow-elevated"
                >
                  <span
                    className={`relative mb-4 inline-flex size-10 items-center justify-center rounded-xl ${colorClasses[color]}`}
                  >
                    <ActionIcon size={20} weight="duotone" />
                  </span>
                  <p className="relative text-sm font-medium text-fg">
                    {t(titleKey)}
                  </p>
                  <p className="relative mt-0.5 text-xs text-fg-tertiary">
                    {t(descKey)}
                  </p>
                </Link>
              ),
            )}
          </div>
        </>
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
              <p className="relative text-sm font-medium text-fg">
                {project.name}
              </p>
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
