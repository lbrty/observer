import { BuildingsIcon, FolderSimpleIcon, GlobeIcon, TagIcon, UsersIcon } from "@/components/icons";
import { createFileRoute, Navigate, Outlet } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { SidebarLink } from "@/components/sidebar-link";
import { useAuth } from "@/stores/auth";

export const Route = createFileRoute("/_app/admin")({
  component: AdminLayout,
});

function AdminLayout() {
  const { t } = useTranslation();
  const { user } = useAuth();

  const isAdmin = user?.role === "admin";
  const isStaff = user?.role === "staff";

  if (!isAdmin && !isStaff) {
    return <Navigate to="/" />;
  }

  return (
    <div className="rays-admin flex flex-1">
      <aside className="w-52 shrink-0 border-r border-border-secondary">
        <nav className="sticky top-13 space-y-0.5 px-3 py-5">
          <div className="pb-1.5 pl-3 text-[11px] font-semibold uppercase tracking-wide text-fg-tertiary">
            {t("admin.title")}
          </div>
          <SidebarLink to="/admin/users" label={t("admin.nav.users")} icon={UsersIcon} />

          {isAdmin && (
            <SidebarLink
              to="/admin/projects"
              label={t("admin.nav.projects")}
              icon={FolderSimpleIcon}
            />
          )}

          <div className="pt-5 pb-1.5 pl-3 text-[11px] font-semibold uppercase tracking-wide text-fg-tertiary">
            {t("admin.nav.reference")}
          </div>
          <SidebarLink
            to="/admin/reference/countries"
            label={t("admin.nav.countries")}
            icon={GlobeIcon}
          />
          <SidebarLink
            to="/admin/reference/offices"
            label={t("admin.nav.offices")}
            icon={BuildingsIcon}
          />
          <SidebarLink
            to="/admin/reference/categories"
            label={t("admin.nav.categories")}
            icon={TagIcon}
          />
        </nav>
      </aside>
      <main className="min-w-0 flex-1 px-8 py-6">
        <Outlet />
      </main>
    </div>
  );
}
